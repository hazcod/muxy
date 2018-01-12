default: build

execName=muxyProxy

clean:
	rm "${execName}" || true

build:
	GOARCH=amd64 CGO_ENALED=0 go build -ldflags="-w -s" -o "${execName}"
	chmod +x "${execName}"

run:
	"./${execName}"
