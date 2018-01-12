package muxy

import (
	"github.com/grafov/m3u8"
	"os"
	"net/http"
	"io"
	"errors"
	log "github.com/golang/glog"
	"strings"
	"strconv"
	"time"
)

var channelCache []Channel

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

func downloadFile(url string, result string) (string, error) {
	out, err := os.Create(result)

	if err != nil {
		return "", errors.New("Could not create temp M3U file: " + err.Error())
	}

	defer out.Close()

	resp, err := http.Get(url)

	if err != nil {
		return "", errors.New("HTTP request failed: " + err.Error())
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		return "", errors.New("Could not copy to temp path: " + err.Error())
	}

	log.Info("Downloaded to " + result)
	return result, nil
}

func removeTempFile() {
	os.Remove(tempM3Upath)
}

func isFileFresh(path string) bool {
	info, err := os.Stat(path)

	if err != nil {
		return false
	}

	expiryDate := time.Now().Add(time.Duration(cacheTimeMinutes * -1) * time.Minute)
	return true == info.ModTime().After(expiryDate)
}

func parseM3UFile(path string) ([]Channel, error) {

	if channelCache != nil && isFileFresh(tempM3Upath) == true {
		log.Info("Using channel cache")
		return channelCache, nil
	}

	if isValidUrl(path) == true {
		log.Info("Downloading " + path)

		downloadedPath, err := downloadFile(path, tempM3Upath)

		if err != nil {
			return nil, err
		}

		path = downloadedPath
	}

	log.Info("Using M3U file: " + path)

	readFile, err := os.Open(path)

	if err != nil {
		return nil, errors.New("Could not open file " + path)
	}

	defer readFile.Close()

	fileStat, err := readFile.Stat()

	if err != nil {
		return nil, errors.New("Could not examine file " + path)
	}

	if fileStat.Size() == 0 {
		return nil, errors.New("M3U file is empty")
	}

	playList, listType, err := m3u8.DecodeFrom(readFile, false)

	if err != nil {
		return nil, errors.New("Could not parse M3U: " + err.Error())
	}

	var mediaWrappper MediaPlaylistWrapper
	switch listType {
		case m3u8.MEDIA:
			mediaWrappper.MediaPlaylist = playList.(*m3u8.MediaPlaylist)
			mediaWrappper.BaseUrl = path
			mediaWrappper.MediaPlaylist.Segments = mediaWrappper.MediaPlaylist.Segments[0 : mediaWrappper.MediaPlaylist.Count()-1]

		default:
			return nil, errors.New("Unknown m3u type")
	}

	var channels []Channel;
	for index, segment := range mediaWrappper.Segments {

		if true == strings.Contains(segment.Title, "▬") {
			continue
		}

		// FetchStreamUrl(segment.URI)

		log.Info("Adding channel{" + "0." + strconv.Itoa(index) + "," + segment.Title + "," + segment.URI + "}")
		channels = append(channels, Channel{"0." + strconv.Itoa(index), segment.Title, segment.URI})
	}

	channelCache = channels

	return channelCache, nil
}

