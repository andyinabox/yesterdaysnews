.PHONY: build
build: clean output/downloads/cnn output/downloads/msnbc output/downloads/foxnews

.PHONY: clean
clean:
	rm -rf output

output:
	-mkdir -p output


output/downloads/cnn: output/cnn.txt
	go run ./cmd/download-multiple/main.go -v -f output/cnn.txt -o output/downloads/cnn

output/downloads/msnbc: output/msnbc.txt
	go run ./cmd/download-multiple/main.go -v -f output/msnbc.txt -o output/downloads/msnbc

output/downloads/foxnews: output/foxnews.txt
	go run ./cmd/download-multiple/main.go -v -f output/foxnews.txt -o output/downloads/foxnews

output/cnn.txt: output
	go run ./cmd/get-vid-ids/main.go -u @CNN -c 10 > output/cnn.txt

output/msnbc.txt: output
	go run ./cmd/get-vid-ids/main.go -u @msnbc -c 10 > output/msnbc.txt

output/foxnews.txt: output
	go run ./cmd/get-vid-ids/main.go -u @FoxNews -c 10 > output/foxnews.txt