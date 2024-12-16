.PHONY: build
build: output/corpus/cnn.txt output/corpus/msnbc.txt output/corpus/foxnews.txt

.PHONY: clean
clean:
	rm -rf output

#
# non-phony targets
#

# combined subtitle text

output/corpus/cnn.txt: corpus output/downloads/cnn
	go run ./cmd/vtt-to-corpus/main.go -f 'output/downloads/cnn/*.vtt' > output/corpus/cnn.txt

output/corpus/msnbc.txt: corpus output/downloads/msnbc
	go run ./cmd/vtt-to-corpus/main.go -f 'output/downloads/msnbc/*.vtt' > output/corpus/msnbc.txt

output/corpus/foxnews.txt: corpus output/downloads/foxnews
	go run ./cmd/vtt-to-corpus/main.go -f 'output/downloads/foxnews/*.vtt' > output/corpus/foxnews.txt

# video and subtitle downloads

output/downloads/cnn: output/cnn.txt
	go run ./cmd/download-multiple/main.go -v -f output/cnn.txt -o output/downloads/cnn

output/downloads/msnbc: output/msnbc.txt
	go run ./cmd/download-multiple/main.go -v -f output/msnbc.txt -o output/downloads/msnbc

output/downloads/foxnews: output/foxnews.txt
	go run ./cmd/download-multiple/main.go -v -f output/foxnews.txt -o output/downloads/foxnews

# video download lists

output/cnn.txt: output
	go run ./cmd/get-vid-ids/main.go -u @CNN -c 10 > output/cnn.txt

output/msnbc.txt: output
	go run ./cmd/get-vid-ids/main.go -u @msnbc -c 10 > output/msnbc.txt

output/foxnews.txt: output
	go run ./cmd/get-vid-ids/main.go -u @FoxNews -c 10 > output/foxnews.txt

# build dirs

corpus:
	-mkdir -p output/corpus

output:
	-mkdir -p output