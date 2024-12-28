.PHONY: upload
build: dist/manifest.json
	go run ./cmd/upload/main.go

.PHONY: clean
clean:
	-rm -rf dist

.PHONY: clobber
clobber: clean
	-rm -rf downloads

.PHONY: objectstoremock
objectstoremock: .cert/localhost.crt
	go run ./cmd/objectstoremock/main.go

.PHONY: serve
serve:
	go run . -a -v -m 20s

#
# non-phony targets
#

dist/manifest.json: dist/yesterdays-news.mp4 dist/yesterdays-news.model.json
	go run ./cmd/createmanifest/main.go

dist/yesterdays-news.mp4: dist/clips/manifest.json
	go run ./cmd/combineclips/main.go -v -i 'dist/clips/*.webm' -o 'dist/yesterdays-news.mp4'

dist/yesterdays-news.model.json: download/cnn/manifest.json download/msnbc/manifest.json download/foxnews/manifest.json
	go run ./cmd/buildmodel/main.go

# video clips

dist/clips/manifest.json: download/cnn/manifest.json download/msnbc/manifest.json download/foxnews/manifest.json
	go run ./cmd/cutvideos/main.go -i 'download/*/*.webm' -o dist/clips

# video and subtitle downloads

download/cnn/manifest.json:
	go run ./cmd/fetchvideos/main.go -n '@cnn' -o download/cnn

download/msnbc/manifest.json:
	go run ./cmd/fetchvideos/main.go -n '@msnbc' -o download/msnbc

download/foxnews/manifest.json:
	go run ./cmd/fetchvideos/main.go -n '@msnbc' -o download/foxnews

