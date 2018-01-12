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
	"bufio"
)

var readBufSizeByes = 10000;

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

			if ! strings.HasSuffix(segment.url, ".ts") {
				log.Error("Not a TS file: " + segment.url)
				sendError(writer)
				return
			}

			fullSegmentUrl := segment.url
			if strings.HasPrefix(fullSegmentUrl, "/") {
				fullSegmentUrl = segmentHost + fullSegmentUrl
			}

			body, err := downloadStreamFile(fullSegmentUrl)
			if err != nil {
				body.Close()
				log.Error("Could not download segment: " + err.Error())
				continue
			}

			writer.Header().Set("Content-Type", "video/mp2t")

			reader := bufio.NewReader(body)

			for {

				line, _, err := reader.ReadLine()

				if err != nil {
					log.Error("Reading line error: " + err.Error())
					body.Close()
					break
				}

				writer.Write(line)
			}

			body.Close()
			waitForNextSegment()
		}

	}

}

func FetchStreamSegments(url string, streamID string) ([]Channel, error) {

	if strings.HasSuffix(url, ".ts") {
		log.Info("No channel playlist, so returning .ts url")
		return []Channel{ {"0.0", streamID,url} }, nil
	}

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

		cleanSegmentTitle := sanitizeName(segment.Title)

		log.Info("Adding Segment{" + "0." + strconv.Itoa(index) + "," + cleanSegmentTitle + "," + segment.URI + "}")
		channels = append(channels, Channel{"0." + strconv.Itoa(index), cleanSegmentTitle, segment.URI})
	}

	return channels, nil
}