.PHONY: upload
build: clean
	go run ./cmd/build/main.go --keepoutput

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
