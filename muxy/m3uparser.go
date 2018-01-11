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
)

type Channel struct {
	number int
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

func downloadM3UFile(url string) (string, error) {
	out, err := os.Create(tempM3Upath)

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

	log.Info("Downloaded to " + tempM3Upath)
	return tempM3Upath, nil
}

func parseM3UFile(path string) ([]Channel, error) {
	if isValidUrl(path) == true {
		log.Info("Downloading " + path)

		downloadedPath, err := downloadM3UFile(path)

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

		case m3u8.MASTER:
			return nil, errors.New("We can't do nothing with a MASTER playlist.")

		default:
			return nil, errors.New("Unknown m3u type")
	}

	var channels []Channel;
	for _, segment := range mediaWrappper.Segments {
		log.Info("Adding channel{" + strconv.Itoa(int(segment.SeqId)) + "," + segment.Title + "," + segment.URI + "}")
		channels = append(channels, Channel{int(segment.SeqId), segment.Title, segment.URI})
	}

	return channels, nil
}

