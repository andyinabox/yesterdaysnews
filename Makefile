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

.PHONY: jstest
jstest:
	npx servor ./app/server js-tests.html 6006


.PHONY: clean-dist clean-bin
clean:

.PHONY: clean-dist
clean-dist:
	-rm -rf dist/*

.PHONY: clean-bin
clean-bin:
	-rm -rf bin/*

#
# docker
# 

.PHONY: docker-build-server
docker-build-server: clean-bin bin/server-linux-amd64
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
