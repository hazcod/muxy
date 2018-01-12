package muxy

import (
	log "github.com/golang/glog"
	"errors"
	"net/http"
	"path/filepath"
	"strings"
	"strconv"
	"net/url"
	"time"
)

func waitForNextSegment() {
	time.Sleep(9 * time.Second)
}

func startChannelStream(writer http.ResponseWriter, channelPlaylist string) {

	streamID := filepath.Base(channelPlaylist)
	streamID = strings.TrimSuffix(streamID, filepath.Ext(streamID))

	segmentHostUrl, err := url.Parse(channelPlaylist)
	if err != nil {
		log.Error("Could not parse host from " + channelPlaylist)
		sendError(writer)
		return
	}

	segmentHost := segmentHostUrl.Scheme + "://" + segmentHostUrl.Host

	log.Info("Streaming stream " + streamID)

	for true {

		segments, err := FetchStreamSegments(channelPlaylist, streamID)
		if err != nil {
			log.Error("Could not fetch channel playlist: " + err.Error())
			sendError(writer)
			return
		}

		for _, segment := range segments {

			if ! strings.HasPrefix(segment.url, ".ts") {
				log.Error("Not a TS file: " + segment.url)
				sendError(writer)
				return
			}

			fullSegmentUrl := segmentHost + segment.url

			segmentBytes, err := downloadFile(fullSegmentUrl)
			if err != nil {
				log.Warning("Error when fetching segment " + fullSegmentUrl + " : " + err.Error())
				sendError(writer)
				return
			}

			log.Info("Sending to client " + strconv.Itoa(len(segmentBytes)) + " bytes")
			writer.Write(segmentBytes)

			waitForNextSegment()
		}

	}

}

func FetchStreamSegments(url string, streamID string) ([]Channel, error) {

	log.Info("Fetching segments for stream " + streamID)

	mediaPlayList, err := parseM3UFile(url)
	if err != nil {
		return nil, errors.New("Could not get channel playlist: " + err.Error())
	}

	var channels []Channel
	for index, segment := range mediaPlayList.Segments {

		if true == strings.Contains(segment.Title, "▬") {
			continue
		}

		log.Info("Adding Segment{" + "0." + strconv.Itoa(index) + "," + segment.Title + "," + segment.URI + "}")
		channels = append(channels, Channel{"0." + strconv.Itoa(index), segment.Title, segment.URI})
	}

	return channels, nil
}