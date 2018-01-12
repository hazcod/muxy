package muxy

import (
	"github.com/grafov/m3u8"
	"os"
	"net/http"
	"errors"
	log "github.com/golang/glog"
	"strings"
	"strconv"
	"encoding/base64"
	"io/ioutil"
	"time"
	"regexp"
)

var cleanNameRegex, _ = regexp.Compile("[^a-zA-Z0-9]+")

type Channel struct {
	number string
	name string
	url string
}


type MediaPlaylistWrapper struct {
	*m3u8.MediaPlaylist
	BaseUrl     string
	VariantInfo string
}

func isValidUrl(toTest string) bool {
	return strings.Contains(strings.ToLower(toTest), "http")
}

func downloadFile(url string) ([]byte, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("Could not request file: " + err.Error())
	}

	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)

	if err != nil {
		return nil, errors.New("HTTP request failed: " + err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Download returned code: " + strconv.Itoa(resp.StatusCode))
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.New("Could not read file contents: " + err.Error())
	}

	if len(bodyBytes) == 0 {
		// retry because something went wrong apparently
		time.Sleep(time.Second * 5)
		return downloadFile(url)
	}

	return bodyBytes, nil
}

func sanitizeName(input string) string {
	return strings.TrimSpace(cleanNameRegex.ReplaceAllString(input, ""))
}

func getChannelPlaylist(m3uPath string) ([]Channel, error) {

	mediaPlayList, err := parseM3UFile(m3uPath)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	var channels []Channel
	for index, segment := range mediaPlayList.Segments {

		if true == strings.Contains(segment.Title, "▬") {
			continue
		}

		encodedStreamURL := base64.StdEncoding.EncodeToString([]byte(segment.URI))
		modifiedSegmentURI := listenUrl + "/stream/" + encodedStreamURL

		cleanSegmentName := sanitizeName(segment.Title)

		log.Info("Adding channel{" + "0." + strconv.Itoa(index) + "," + cleanSegmentName + "," + modifiedSegmentURI + "}")
		channels = append(channels, Channel{"0." + strconv.Itoa(index), cleanSegmentName, modifiedSegmentURI})
	}

	return channels, nil
}

func parseM3UFile(path string) (MediaPlaylistWrapper, error) {

	var mediaWrappper MediaPlaylistWrapper
	var m3uContent string

	if isValidUrl(path) == true {
		log.Info("Downloading " + path)

		str, err := downloadFile(path)
		if err != nil {
			return mediaWrappper, err
		}

		m3uContent = string(str)
	} else {
		log.Info("Using M3U file: " + path)

		readFile, err := os.Open(path)
		if err != nil {
			return mediaWrappper, errors.New("Could not open file " + path)
		}

		defer readFile.Close()

		fileStat, err := readFile.Stat()
		if err != nil {
			return mediaWrappper, errors.New("Could not examine file " + path)
		}

		if fileStat.Size() == 0 {
			return mediaWrappper, errors.New("M3U file is empty")
		}

		bytesRead, err := ioutil.ReadAll(readFile)
		if err != nil {
			return mediaWrappper, errors.New("Could not read M3U: " + err.Error())
		}

		m3uContent = string(bytesRead)
	}

	playList, listType, err := m3u8.DecodeFrom(strings.NewReader(m3uContent), false)

	if err != nil {
		return mediaWrappper, errors.New("Could not parse M3U: " + err.Error())
	}

	switch listType {
		case m3u8.MEDIA:
			mediaWrappper.MediaPlaylist = playList.(*m3u8.MediaPlaylist)
			mediaWrappper.BaseUrl = path
			mediaWrappper.MediaPlaylist.Segments = mediaWrappper.MediaPlaylist.Segments[0 : mediaWrappper.MediaPlaylist.Count()-1]

		default:
			return mediaWrappper, errors.New("unknown m3u type")
	}

	return mediaWrappper, nil
}

