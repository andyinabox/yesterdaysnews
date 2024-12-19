.PHONY: build
build: dist/manifest.json

.PHONY: clean
clean:
	-rm -rf dist
	-rm -rf downloads

.PHONY: clobber
clobber: clean
	-rm -rf downloads

#
# non-phony targets
#

dist/manifest.json: dist/yesterdays-news.mp4 dist/yesterdays-news.model.json
	go run ./cmd/create-dist-manifest/main.go

dist/yesterdays-news.mp4: dist/clips/manifest.json
	go run ./cmd/combine-clips/main.go -v -i 'dist/clips/*.webm' -o 'dist/yesterdays-news.mp4'

dist/yesterdays-news.model.json: download/cnn/manifest.json download/msnbc/manifest.json download/foxnews/manifest.json
	go run ./cmd/output-markov-model/main.go

# video clips

dist/clips/manifest.json: download/cnn/manifest.json download/msnbc/manifest.json download/foxnews/manifest.json
	go run ./cmd/cut-multiple-videos/main.go -i 'download/*/*.webm' -o dist/clips

# combined subtitle text

dist/txt/cnn.txt: download/cnn/manifest.json
	-mkdir -p dist/txt
	go run ./cmd/vtt-to-corpus/main.go -f 'download/cnn/*.vtt' > $@

dist/txt/msnbc.txt: download/msnbc/manifest.json
	-mkdir -p dist/txt
	go run ./cmd/vtt-to-corpus/main.go -f 'download/msnbc/*.vtt' > $@

dist/txt/foxnews.txt: download/foxnews/manifest.json
	-mkdir -p dist/txt
	go run ./cmd/vtt-to-corpus/main.go -f 'download/foxnews/*.vtt' > $@

# video and subtitle downloads

download/cnn/manifest.json:
	go run ./cmd/download-multiple/main.go -n '@cnn' -o download/cnn

download/msnbc/manifest.json:
	go run ./cmd/download-multiple/main.go -n '@msnbc' -o download/msnbc

download/foxnews/manifest.json:
	go run ./cmd/download-multiple/main.go -n '@msnbc' -o download/foxnews

