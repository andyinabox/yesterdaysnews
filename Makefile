#
# app runners
# 


# builder

.PHONY: builder
builder:
	go run ./app/builder/main.go

.PHONY: builder-local
builder-local:
	go run ./app/builder/main.go --keepoutput --skipupload


# server

# run server, using object store credential in .env
.PHONY: server
server:
	go run ./app/server/main.go -a -v -m 20s


# run server using objectstoremock
# (need to already be running objectstoremock)
.PHONY: server-local
server-local:
	YN_OBJECTSTORE_URL=http://localhost:9000 YN_CDN_URL=http://localhost:9000 go run ./app/server/main.go -a -v -m 20s


#
# utils
# 

.PHONY: objectstoremock
objectstoremock:
	go run ./cmd/objectstoremock/main.go

.PHONY: clean-dist clean-bin clean-assets
clean:

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
# docker
# 

# server

.PHONY: docker-run-server
docker-run-server: clean-bin clean-assets bin/server-linux-amd64
	docker buildx build --platform linux/amd64 -f app/server/Dockerfile -t andyinabox/yesterdaysnews-server:dev .
	docker run --rm --env-file .env -p 8080:8080 andyinabox/yesterdaysnews-server:dev

# builder

.PHONY: docker-run-builder
docker-run-builder: clean-bin bin/builder-linux-amd64
	docker buildx build --platform linux/amd64 -f app/builder/Dockerfile -t andyinabox/yesterdaysnews-builder:dev .
	mkdir -p dist
	docker run --rm --env-file .env  -v ./dist:/dist andyinabox/yesterdaysnews-builder:dev --output /dist -v

#
# file-based targets
# 

# binaries

bin/server-linux-amd64: app/server/.assets
	GOOS=linux GOARCH=amd64 go build -o $@ ./app/server/main.go

bin/builder-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o $@ ./app/builder/main.go

# assets

app/server/.assets: app/server/.assets/styles.css app/server/.assets/script.js

app/server/.assets/styles.css:
		go run ./cmd/esbuild/main.go app/server/assets/styles.css --bundle --minify --outfile=app/server/.assets/styles.css

app/server/.assets/script.js:
		go run ./cmd/esbuild/main.go app/server/assets/script.js --bundle --minify --outfile=app/server/.assets/script.js



