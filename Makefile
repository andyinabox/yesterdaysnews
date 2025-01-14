.PHONY: upload
build: clobber
	go run ./cmd/build/main.go -d 10

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
