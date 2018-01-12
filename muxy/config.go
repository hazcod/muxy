package muxy

import "strconv"

var m3ufile = "playlist.m3u"
var listenPort = 8080
var listenHost = "localhost"
var tunerCount = 1

var listenUrl = "http://" + listenHost + ":" + strconv.Itoa(listenPort)
var deviceInfo = map[string]string{
	"FriendlyName": "muxy",
	"Manufacturer" : "Silicondust",
	"ModelNumber": "HDTC-2US",
	"FirmwareName": "hdhomeruntc_atsc",
	"TunerCount": strconv.Itoa(tunerCount),
	"FirmwareVersion": "20150826",
	"DeviceID": "12345678",
	"DeviceAuth": "test1234",
	"BaseURL": listenUrl + "/",
	"LineupURL": listenUrl + "/lineup.json",
}
