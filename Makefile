.PHONY: upload
build: clean
	go run ./cmd/build/main.go --keepoutput

.PHONY: build-setup
build-setup:
	go run ./cmd/build/main.go -v -b setup

.PHONY: build-video-clips
build-video-clips:
	go run ./cmd/build/main.go -v -b video-clips --playlistids UUaXkIU1QidjPwiAYu6GcHjg --count 1

.PHONY: build-model
build-model:
	go run ./cmd/build/main.go -v -b model

.PHONY: build-promote
build-promote:
	go run ./cmd/build/main.go -v -b promote

.PHONY: build-cleanup
build-cleanup:
	go run ./cmd/build/main.go -v -b cleanup

.PHONY: clean
clean:
	-rm -rf dist

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
	docker run --env YN_OBJECTSTORE_URL=https://sos-ch-dk-2.exo.io/yesterdaysnews -p 8080:8080 andyinabox/yesterdaysnews-server

.PHONY: docker-push-server
docker-push-server:
	docker push andyinabox/yesterdaysnews-server
