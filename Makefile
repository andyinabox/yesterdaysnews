#
# builder
# 

.PHONY: builder
builder:
	go run ./app/builder/main.go --keepoutput

.PHONY: builder-local
builder-local:
	go run ./app/builder/main.go --keepoutput --skipupload


#
# server
# 

# run server, using object store credential in .env
.PHONY: server
server: app/server/.assets
	go run ./app/server/main.go -a -v -m 20s


# run server using objectstoremock
# (need to already be running objectstoremock)
.PHONY: server-local
server-local: app/server/.assets
	YN_OBJECTSTORE_URL=http://localhost:9000 YN_CDN_URL=http://localhost:9000 go run ./app/server/main.go -a -v -m 20s


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

.PHONY: clean-assets
clean-assets:
	-rm -rf app/server/.assets


#
# docker
# 

.PHONY: docker-build-server
docker-build-server: clean-bin clean-assets bin/server-linux-amd64
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

bin/server-linux-amd64: app/server/.assets
	GOOS=linux GOARCH=amd64 go build -o $@ ./app/server/main.go

bin/builder-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o $@ ./app/builder/main.go


app/server/.assets: app/server/.assets/styles.css app/server/.assets/script.js

app/server/.assets/styles.css:
		go run ./cmd/esbuild/main.go app/server/assets/styles.css --bundle --minify --outfile=app/server/.assets/styles.css

app/server/.assets/script.js:
		go run ./cmd/esbuild/main.go app/server/assets/script.js --bundle --minify --outfile=app/server/.assets/script.js



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
