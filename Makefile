version = 0.0.1

build:
	@mkdir -p dist
	CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -ldflags="-w -s" -o bin/linux-amd64/pass-server .
	GOOS=darwin  GOARCH=amd64 go build -o bin/macos-amd64/pass-server .
	GOOS=darwin  GOARCH=arm64 go build -o bin/macos-arm64/pass-server .
	GOOS=windows GOARCH=amd64 go build -o bin/windows-amd64/pass-server .

clean:
	rm -fr bin/*
	rm -fr dist/*

pkg: build
	@mkdir -p dist
	zip -j dist/linux-${version}.zip bin/linux-amd64/pass-server
	zip -j dist/macos-${version}.zip bin/macos-amd64/pass-server
	zip -j dist/macos-${version}.zip bin/macos-arm64/pass-server
	zip -j dist/windows-${version}.zip bin/windows-amd64/pass-server
