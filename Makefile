.PHONY: build
build: clean output/cnn.txt output/msnbc.txt output/foxnews.txt

.PHONY: clean
clean:
	rm -rf output

output:
	-mkdir -p output

output/cnn.txt: output
	go run ./cmd/get-vid-ids/main.go -u @CNN -c 10 > output/cnn.txt

output/msnbc.txt: output
	go run ./cmd/get-vid-ids/main.go -u @msnbc -c 10 > output/msnbc.txt

output/foxnews.txt: output
	go run ./cmd/get-vid-ids/main.go -u @FoxNews -c 10 > output/foxnews.txt