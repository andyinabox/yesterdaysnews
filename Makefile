.PHONY: build
build: output/corpus/cnn.txt output/corpus/msnbc.txt output/corpus/foxnews.txt

.PHONY: clean
clean:
	rm -rf output

#
# non-phony targets
#

# video clips

# output/clips/cnn: output/downloads/cnn
# 	go run ./cmd/cut-multiple-videos/main.go -i 'output/downloads/cnn/*.mp4' -o $@

# output/clips/msnbc: output/downloads/msnbc
# 	go run ./cmd/cut-multiple-videos/main.go -i 'output/downloads/msnbc/*.mp4' -o $@

# output/clips/foxnews: output/downloads/foxnews
# 	go run ./cmd/cut-multiple-videos/main.go -i 'output/downloads/foxnews/*.mp4' -o $@


# markov models

# output/models/cnn.json: models output/corpus/cnn.txt
# 	go run ./cmd/build-model/main.go -f output/corpus/cnn.txt > $@

# output/models/msnbc.json: models output/corpus/msnbc.txt
# 	go run ./cmd/build-model/main.go -f output/corpus/msnbc.txt > $@

# output/models/foxnews.json: models output/corpus/foxnews.txt
# 	go run ./cmd/build-model/main.go -f output/corpus/foxnews.txt > $@


# combined subtitle text

output/corpus/cnn.txt: corpus output/downloads/cnn
	go run ./cmd/vtt-to-corpus/main.go -f 'output/downloads/cnn/*.vtt' > $@

output/corpus/msnbc.txt: corpus output/downloads/msnbc
	go run ./cmd/vtt-to-corpus/main.go -f 'output/downloads/msnbc/*.vtt' > $@

output/corpus/foxnews.txt: corpus output/downloads/foxnews
	go run ./cmd/vtt-to-corpus/main.go -f 'output/downloads/foxnews/*.vtt' > $@

# video and subtitle downloads

output/downloads/cnn: output/cnn.txt
	go run ./cmd/download-multiple/main.go -v -f output/cnn.txt -o $@

output/downloads/msnbc: output/msnbc.txt
	go run ./cmd/download-multiple/main.go -v -f output/msnbc.txt -o $@

output/downloads/foxnews: output/foxnews.txt
	go run ./cmd/download-multiple/main.go -v -f output/foxnews.txt -o $@

# video download lists

output/cnn.txt: output
	go run ./cmd/get-vid-ids/main.go -u @CNN -c 10 > $@

output/msnbc.txt: output
	go run ./cmd/get-vid-ids/main.go -u @msnbc -c 10 > $@

output/foxnews.txt: output
	go run ./cmd/get-vid-ids/main.go -u @FoxNews -c 10 > $@

# build dirs

output:
	-mkdir -p output

corpus:
	-mkdir -p output/corpus

models:
	-mkdir -p output/models