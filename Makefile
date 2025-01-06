.PHONY: upload
build: clobber
	go run ./cmd/build/main.go

.PHONY: clean
clean:
	-rm -rf dist

.PHONY: clobber
clobber: clean
	-rm -rf download

.PHONY: objectstoremock
objectstoremock:
	go run ./cmd/objectstoremock/main.go

.PHONY: serve
serve:
	go run . -a -v -m 20s

.PHONY: docker-build-server
docker-build-server:
	docker build -f docker/server.Dockerfile -t andyinabox/yesterdaysnews-server .

.PHONY: docker-run-server
docker-run-server:
	docker run --env YN_OBJECTSTORE_URL=https://s3.pub1.infomaniak.cloud/object/v1/AUTH_b28bee57175648379ec940f55adfe842/yesterdaysnews/current -p 8080:8080 andyinabox/yesterdaysnews-server

.PHONY: docker-push-server
docker-push-server:
	docker push andyinabox/yesterdaysnews-server

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

