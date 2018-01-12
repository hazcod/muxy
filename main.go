package main

import (
	"muxy/muxy"
	"flag"
	"fmt"
)

func main() {
	listenHost := flag.String("host", "localhost", "What address to listen on.")
	listenPort := flag.Int("port", 8080, "What port to listen on.")
	tempPath   := flag.String("temp", "/tmp/m3u", "What location to store the downloaded M3U file.")
	maxStreams := flag.Int("streams", 1, "How many streams can be played simultaneously.")
	flag.Parse()

	m3uPath := flag.Arg(0)

	if m3uPath == "" {
		fmt.Print("Usage: ./muxyProxy <path-to-m3u>")
		return
	}

	muxy.SetListenHost(*listenHost)
	muxy.SetListenPort(*listenPort)
	muxy.SetTempM3UPath(*tempPath)
	muxy.SetMaxStreams(*maxStreams)
	muxy.SetM3UFile(m3uPath)

	muxy.RunListener()
}
