package muxy

import "strconv"

var m3ufile = "playlist.m3u"
var listenPort = 8080
var listenHost = "localhost"
var tunerCount = 1
var userAgent = "QuickTime/7.5";
var listenUrl = "http://" + listenHost + ":" + strconv.Itoa(listenPort)

var deviceInfo = map[string]interface{}{
	"FriendlyName": "muxy",
	"Manufacturer" : "Silicondust",
	"ModelNumber": "HDTC-2US",
	"FirmwareName": "hdhomeruntc_atsc",
	"TunerCount": tunerCount,
	"FirmwareVersion": "20170930",
	"DeviceID": "12345678",
	"DeviceAuth": "lasdfkgkdskfksdkfsds",
	"BaseURL": listenUrl,
	"LineupURL": listenUrl + "/lineup.json",
}
