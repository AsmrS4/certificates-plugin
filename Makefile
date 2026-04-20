.PHONY: build clean

build:
	GOOS=wasip1 GOARCH=wasm go build -ldflags="-s -w" -o certificates-plugin.wasm .

clean:
	rm -f certificates-plugin.wasm