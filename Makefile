.PHONY: upload
build: clean
	go run ./cmd/builder/main.go --keepoutput

.PHONY: build-setup
build-setup:
	go run ./cmd/builder/main.go -v -b setup

.PHONY: build-video-clips
build-video-clips:
	go run ./cmd/builder/main.go -v -b video-clips --playlistids UUaXkIU1QidjPwiAYu6GcHjg --count 1

.PHONY: build-model
build-model:
	go run ./cmd/builder/main.go -v -b model

.PHONY: build-promote
build-promote:
	go run ./cmd/builder/main.go -v -b promote

.PHONY: build-cleanup
build-cleanup:
	go run ./cmd/builder/main.go -v -b cleanup

.PHONY: clean
clean:
	-rm -rf dist/*

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
	docker run --rm --env-file .env -p 8080:8080 andyinabox/yesterdaysnews-server

.PHONY: docker-push-server
docker-push-server:
	docker push andyinabox/yesterdaysnews-server



.PHONY: docker-build-builder
docker-build-builder:
	docker buildx build --platform linux/arm64 -f docker/builder.Dockerfile -t andyinabox/yesterdaysnews-builder .

.PHONY: docker-run-builder
docker-run-builder:
	mkdir -p dist
	docker run --rm --env-file .env  -v ./dist:/dist andyinabox/yesterdaysnews-builder --output /dist -v

.PHONY: docker-push-builder
docker-push-builder:
	docker push andyinabox/yesterdaysnews-builder
