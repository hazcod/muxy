package muxy

import (
	"os"
	"github.com/grafov/m3u8"
	log "github.com/golang/glog"
	"errors"
)

func FetchStreamUrl(url string) (string, error) {

	tmp, err := downloadFile(url, tempTSpath)

	if err != nil {
		return "", errors.New("Could not download TS file " + url)
	}
	readFile, err := os.Open(tmp)

	if err != nil {
		return "", errors.New("Could not open TS file " + tmp)
	}

	defer readFile.Close()

	fileStat, err := readFile.Stat()

	if err != nil {
		return "", errors.New("Could not examine TS file " + tmp)
	}

	if fileStat.Size() == 0 {
		return "", errors.New("Downloaded TS file is empty")
	}

	playList, listType, err := m3u8.DecodeFrom(readFile, false)

	if err != nil {
		return "", errors.New("Could not parse M3U: " + err.Error())
	}

	var mediaWrappper MediaPlaylistWrapper
	switch listType {
		case m3u8.MEDIA:
			mediaWrappper.MediaPlaylist = playList.(*m3u8.MediaPlaylist)
			mediaWrappper.BaseUrl = tmp
			mediaWrappper.MediaPlaylist.Segments = mediaWrappper.MediaPlaylist.Segments[0 : mediaWrappper.MediaPlaylist.Count()-1]

		default:
			return "", errors.New("Unknown m3u type")
	}

	for _, segment := range mediaWrappper.Segments {
		log.Info("File download URL: " + segment.URI)
	}

	return "", nil
}