package muxy

import "strconv"

var m3ufile = "playlist.m3u"
var listenPort = 8080
var listenHost = "127.0.0.1"
var tunerCount = 1
var userAgent = "QuickTime/7.5";
var listenUrl = "http://" + listenHost + ":" + strconv.Itoa(listenPort)

var deviceInfo = map[string]interface{}{
	"FriendlyName": "muxy",
	"Manufacturer" : "Silicondust",
	"ModelNumber": "HDHR4-2US",
	"FirmwareName": "hdhomerun4_atsc",
	"TunerCount": tunerCount,
	"FirmwareVersion": "20170930",
	"DeviceID": "10439EFD",
	"DeviceAuth": "KOxavUdByRLBRKZRsV/ge8lS",
	"BaseURL": listenUrl,
	"LineupURL": listenUrl + "/lineup.json",
}
