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
	"encoding/base64"
)

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

	if err != nil || (resp != nil && (resp.StatusCode >= 400)){
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
	os.Remove(tempM3UFile)
}

/*
func isFileFresh(path string) bool {
	info, err := os.Stat(path)

	if err != nil {
		return false
	}

	expiryDate := time.Now().Add(time.Duration(cacheTimeMinutes * -1) * time.Minute)
	return true == info.ModTime().After(expiryDate)
}
*/

func getChannelPlaylist(m3uPath string) ([]Channel, error) {

	mediaPlayList, err := parseM3UFile(m3uPath, tempM3UFile)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	var channels []Channel;
	for index, segment := range mediaPlayList.Segments {

		if true == strings.Contains(segment.Title, "▬") {
			continue
		}

		encodedStreamURL := base64.StdEncoding.EncodeToString([]byte(segment.URI))
		modifiedSegmentURI := listenUrl + "/stream/" + encodedStreamURL

		log.Info("Adding channel{" + "0." + strconv.Itoa(index) + "," + segment.Title + "," + modifiedSegmentURI + "}")
		channels = append(channels, Channel{"0." + strconv.Itoa(index), segment.Title, modifiedSegmentURI})
	}

	return channels, nil
}

func parseM3UFile(path string, temp string) (MediaPlaylistWrapper, error) {

	var mediaWrappper MediaPlaylistWrapper

	if isValidUrl(path) == true {
		log.Info("Downloading " + path)

		downloadedPath, err := downloadFile(path, temp)

		if err != nil {
			return mediaWrappper, err
		}

		path = downloadedPath
	}

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

	playList, listType, err := m3u8.DecodeFrom(readFile, false)

	if err != nil {
		return mediaWrappper, errors.New("Could not parse M3U: " + err.Error())
	}

	switch listType {
		case m3u8.MEDIA:
			mediaWrappper.MediaPlaylist = playList.(*m3u8.MediaPlaylist)
			mediaWrappper.BaseUrl = path
			mediaWrappper.MediaPlaylist.Segments = mediaWrappper.MediaPlaylist.Segments[0 : mediaWrappper.MediaPlaylist.Count()-1]

		default:
			return mediaWrappper, errors.New("Unknown m3u type")
	}

	return mediaWrappper, nil
}

