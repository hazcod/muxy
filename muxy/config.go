package muxy

import "strconv"

var m3ufile string
var listenPort int
var listenHost string
var tunerCount int

var userAgent = "VLC";
var listenUrl = "http://" + listenHost + ":" + strconv.Itoa(listenPort)

var deviceInfo = map[string]interface{}{
	"FriendlyName": "muxy",
	"Manufacturer" : "Silicondust",
	"ModelNumber": "HDHR4-2US",
	"FirmwareName": "hdhomerun4_atsc",
	"TunerCount": tunerCount,
	"FirmwareVersion": "20150826",
	"DeviceID": "10439EFD",
	"DeviceAuth": "KOxavUdByRLBRKZRsV/ge8lS",
	"BaseURL": listenUrl,
	"LineupURL": listenUrl + "/lineup.json",
}
