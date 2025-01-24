
#
# builder
# 

.PHONY: build
build:
	go run ./app/builder/main.go --keepoutput

#
# server
# 

.PHONY: serve
serve:	go run ./app/server/main.go -a -v -m 20s


#
# utils
# 

.PHONY: objectstoremock
objectstoremock:
	go run ./cmd/objectstoremock/main.go

.PHONY: clean-dist clean-bin
clean:

.PHONY: clean-dist
clean-dist:
	-rm -rf dist/*

.PHONY: clean-bin
clean-bin:
	-rm -rf bin/*

.PHONY: clean-bin binaries
binaries: bin/server-linux-amd64 bin/builder-linux-amd64

#
# docker
# 

.PHONY: docker-build-server
docker-build-server: bin/server-linux-amd64
	docker build -f app/server/Dockerfile -t andyinabox/yesterdaysnews-server .

.PHONY: docker-run-server
docker-run-server:
	docker run --rm --env-file .env -p 8080:8080 andyinabox/yesterdaysnews-server

.PHONY: docker-push-server
docker-push-server:
	docker push andyinabox/yesterdaysnews-server


#
# file-based targets
# 

bin/server-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o $@ ./app/server/main.go

bin/builder-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o $@ ./app/builder/main.go



# .PHONY: docker-build-builder
# docker-build-builder:
# 	docker buildx build --platform linux/arm64 -f docker/builder.Dockerfile -t andyinabox/yesterdaysnews-builder .

# .PHONY: docker-run-builder
# docker-run-builder:
# 	mkdir -p dist
# 	docker run --rm --env-file .env  -v ./dist:/dist andyinabox/yesterdaysnews-builder --output /dist -v

# .PHONY: docker-push-builder
# docker-push-builder:
# 	docker push andyinabox/yesterdaysnews-builder
