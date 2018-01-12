package muxy

import (
	"net/http"
	"encoding/json"
	"strconv"
	"github.com/gorilla/mux"
	log "github.com/golang/glog"
)

func sendError(w http.ResponseWriter) {
	w.WriteHeader(500)
}

func sendJson(w http.ResponseWriter, data interface{}) {
	var dataBytes []byte

	if data != "" {
		var jsonData, err= json.Marshal(data)

		if err != nil {
			log.Error("Could not encode to JSON: " + err.Error())
			sendError(w)
			return
		}

		dataBytes = jsonData
	} else {
		dataBytes = []byte(data.(string))
	}

	log.Info("Sending: " + string(dataBytes))

	w.Write(dataBytes)
}

func getLineupStatus(w http.ResponseWriter, r *http.Request) {
	sendJson(w, map[string]interface{}{
		"ScanInProgress": "0",
		"ScanPossible": "0",
		"Source": "Cable",
		"SourceList": []string{"Cable", "Antenna"},
	})
}

func getLineup(w http.ResponseWriter, r *http.Request) {
	channels, err := parseM3UFile(m3ufile)

	if err != nil {
		log.Error("Could not get channels: " + err.Error())
		sendError(w)
		return
	}

	var lineup []map[string]string;

	for _, channel := range channels {
		lineup = append(lineup, map[string]string{
			"GuideNumber": channel.number,
			"GuideName": channel.name,
			"URL": channel.url,
		})
	}

	sendJson(w, lineup)
}

func getDeviceInfo(w http.ResponseWriter, r *http.Request) {
	sendJson(w, deviceInfo)
}

func doNothing(w http.ResponseWriter, r *http.Request) {
	sendJson(w, nil)
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Info(r.RemoteAddr + " " + r.Method + " " + r.URL.Path)

		if r.Method == "POST" {
			r.ParseForm()
			log.Info("Body: " + r.Form.Encode())
		}

		handler.ServeHTTP(w, r)
	})
}

func SetM3UFile(path string) {
	m3ufile = path
}

func SetMaxStreams(num int) {
	tunerCount = num
}

func SetListenHost(host string) {
	listenHost = host
}

func SetListenPort(port int) {
	listenPort = port
}

func SetTempM3UPath(path string) {
	tempM3Upath = path
}

func RunListener() {
	router := mux.NewRouter()

	router.HandleFunc("/device.json", getDeviceInfo).Methods("GET", "POST")
	router.HandleFunc("/discover.json", getDeviceInfo).Methods("GET", "POST")

	router.HandleFunc("/lineup_status.json", getLineupStatus).Methods("GET", "POST")
	router.HandleFunc("/lineup.json", getLineup).Methods("GET", "POST")

	router.HandleFunc("/lineup.post", doNothing).Methods("GET", "POST")

	removeTempFile()

	err := http.ListenAndServe(
		listenHost + ":" + strconv.Itoa(listenPort),
		logRequest(router),
	)

	if err != nil {
		log.Error("Could not start listener: " + err.Error())
	}
}
