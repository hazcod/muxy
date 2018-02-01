# muxy
Emulates a HDHomeRun device while streaming from M3U IPTV streams.

Add `localhost:8080` as a DVR tuner in Plex.

## WARNING
As of now, Plex does accept TS streams, meaning this code is far too complex for what it's worth.
Please check out a much simpler implementation such as `telly` instead: https://github.com/tombowditch/telly

## Usage
`./muxyProxy http://site.com/my-iptv-playlist.m3u8`

## Why?
Plex currently does not support IPTV directly, only using network tuners such as a HDHomeRun.
tvhProxy is cumbersome because it requires a tvheadend installation, resulting in extra latency.
Ideally, Plex would need practically direct access to the IPTV streams and do all the DVR stuff itself.

## How does it work?
It reads your M3U8 playlist to show Plex a channel list. The download link in the playlist is a reference to `muxy` itself,
with the download link base64 encoded in the URI. When plex requests a file, `muxy` downloads the cronological TS files
and serves the MPEG stream to Plex.

## Debugging
Run muxy in verbose logging mode: `./muxyProxy -v 10 -logtostderr <path-to-m3u>`
Try accessing `http://localhost:8080/lineup.json` and try out one of the stream URLs. (e.g. `http://localhost:8080/stream/xxxxx`)

## Building
You just need `go` and `automake`. To build, do a `make` in the source directory. Your executable will be `muxyProxy`.

## Credits
Big thanks goes out to @jkaberg for his work on `tvhProxy`, where I got the idea from.
