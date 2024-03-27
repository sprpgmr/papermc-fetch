build: test
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o papermc-fetch-linux-amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o papermc-fetch-linux-arm
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o papermc-fetch-linux-arm64

	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o papermc-fetch-windows-64bit.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o papermc-fetch-windows-32bit.exe

	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o papermc-fetch-macos-intel
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o papermc-fetch-macos-apple-silicon

test:
	go test ./...