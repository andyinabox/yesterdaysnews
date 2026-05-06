# Detect host architecture for local builds
ARCH := $(shell uname -m)
ifeq ($(ARCH),arm64)
  GOARCH := arm64
else ifeq ($(ARCH),aarch64)
  GOARCH := arm64
else
  GOARCH := amd64
endif

# Pass YN_* env vars (loaded by direnv) through to Docker containers
DOCKER_ENV := \
  -e YN_OBJECTSTORE_URL \
  -e YN_CDN_URL \
  -e YN_GOOGLE_API_KEY \
  -e YN_S3_ENDPOINT \
  -e YN_S3_REGION \
  -e YN_S3_BUCKET_NAME \
  -e YN_S3_ACCESS_KEY \
  -e YN_S3_SECRET_ACCESS_KEY

#
# app runners
#


# builder

.PHONY: builder
builder:
	go run ./app/builder/main.go -v --keepoutput 2>&1 | tee builder.log

.PHONY: builder-local
builder-local:
	go run ./app/builder/main.go -v --keepoutput --skipupload 2>&1 | tee builder-local.log

.PHONY: builder-docker
builder-docker: clean-bin bin/builder-linux-$(GOARCH)
	docker buildx build --platform linux/$(GOARCH) -f app/builder/Dockerfile -t yesterdaysnews-builder:dev .
	mkdir -p dist
	docker run --rm $(DOCKER_ENV) -v ./dist:/dist yesterdaysnews-builder:dev --output /dist -v --keepoutput --skipupload --throttledl 1s

.PHONY: builder-docker-amd64
builder-docker-amd64: clean-bin bin/builder-linux-amd64
	docker buildx build --platform linux/amd64 -f app/builder/Dockerfile -t yesterdaysnews-builder:dev .

# server

# run server, using object store credentials from direnv-loaded environment
.PHONY: server
server:
	go run ./app/server/main.go -a -v -m 20s


# run server using objectstoremock
# (need to already be running objectstoremock)
.PHONY: server-local
server-local:
	YN_OBJECTSTORE_URL=http://localhost:9000 YN_CDN_URL=http://localhost:9000 go run ./app/server/main.go -a -v -m 20s


.PHONY: server-docker
server-docker: clean-bin clean-assets bin/server-linux-amd64
	docker buildx build --platform linux/amd64 -f app/server/Dockerfile -t yesterdaysnews-server:dev .
	docker run --rm $(DOCKER_ENV) -p 8080:8080 yesterdaysnews-server:dev

#
# test
#

.PHONY: test
test:
	go test ./...

#
# utils
#

.PHONY: objectstoremock
objectstoremock:
	go run ./cmd/objectstoremock/main.go

.PHONY: clean
clean: clean-dist clean-bin clean-assets

.PHONY: clean-dist
clean-dist:
	-rm -rf dist/*

.PHONY: clean-bin
clean-bin:
	-rm -rf bin/*

.PHONY: clean-assets
clean-assets:
	-rm -rf app/server/.assets

#
# file-based targets
# 

# binaries

bin/server-linux-amd64: clean-assets app/server/.assets
	GOOS=linux GOARCH=amd64 go build -o $@ ./app/server/main.go

bin/builder-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o $@ ./app/builder/main.go

bin/builder-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o $@ ./app/builder/main.go

# assets

app/server/.assets: app/server/.assets/styles.css app/server/.assets/script.js
	cp app/server/assets/*.png app/server/.assets/
	cp app/server/assets/*.svg app/server/.assets/
	cp -r app/server/assets/icon app/server/.assets/icon

app/server/.assets/styles.css:
		go run ./cmd/esbuild/main.go app/server/assets/styles.css --bundle --minify --outfile=app/server/.assets/styles.css

app/server/.assets/script.js:
		go run ./cmd/esbuild/main.go app/server/assets/script.js --bundle --minify --outfile=app/server/.assets/script.js



