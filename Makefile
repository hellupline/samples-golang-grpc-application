APPLICATION_NAME=sample-application

.PHONY: run-server
run-server:
	GO111MODULE=on go run ./cmd/server

.PHONY: run-application
run-application:
	npm --prefix ./app run start


.PHONY: clean
clean:
	rm -rf ${APPLICATION_NAME}-server ${APPLICATION_NAME}-client _swagger/apidocs.swagger.json _tls/ internal/static/ pkg/api/favorites/


# build binaries
.PHONY: go-build
go-build: ${APPLICATION_NAME}-server ${APPLICATION_NAME}-client
${APPLICATION_NAME}-%: cmd/%/ | protoc static-tlsdata static-openapidata go-mod-vendor
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o $@ ./$<
.PHONY: go-mod-vendor
go-mod-vendor: vendor/modules.txt
vendor/modules.txt: go.mod
	GO111MODULE=on GOOS=linux GOARCH=amd64 go mod vendor
.PHONY: go-vet
go-vet:
	GO111MODULE=on go vet ./...
.PHONY: go-fmt
go-fmt:
	GO111MODULE=on go fmt ./...
.PHONY: go-test
go-test:
	GO111MODULE=on go test ./...


# build api files
.PHONY: protoc
protoc: _swagger/apidocs.swagger.json
_swagger/apidocs.swagger.json: protobuf/*.proto
	protoc \
		-I /usr/local/include \
		-I ${GOPATH}/src \
		-I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		-I ./protobuf \
		--swagger_out=logtostderr=true,allow_merge=true:_swagger \
		--go_out=plugins=grpc:pkg/api \
		--grpc-gateway_out=logtostderr=true:pkg/api \
		$^


# bundle react app
.PHONY: static-appdata
static-appdata: internal/static/appdata/static.go
internal/static/appdata/static.go: app/build/index.html
	statik -f -m -p appdata -src app/build -dest internal/static/
app/build/index.html:
	npm --prefix ./app ci
	npm --prefix ./app run build


# bundle swagger docs
.PHONY: static-openapidata
static-openapidata: internal/static/openapidata/static.go
internal/static/openapidata/static.go: _swagger/apidocs.swagger.json
	statik -f -m -p openapidata -src _swagger -dest internal/static/


# bundle tls certs
.PHONY: static-tlsdata
static-tlsdata: internal/static/tlsdata/static.go
internal/static/tlsdata/static.go: _tls/service.pem
	statik -f -m -p tlsdata -src _tls -dest internal/static/
_tls:
	mkdir $@
_tls/rootca.key: | _tls
	openssl genrsa -out $@ 4096
_tls/rootca.cert: _tls/rootca.key | _tls
	openssl req -x509 -new -sha256 -days 3650 -key $< -out $@ -subj "/O=${APPLICATION_NAME}"
_tls/service.key: | _tls
	openssl genrsa -out $@ 4096
_tls/service.csr: _tls/service.key | _tls
	openssl req -new -sha256 -days 3650 -addext "subjectAltName = DNS:localhost,IP:::1,IP:127.0.0.1" -key $< -out $@ -subj "/O=${APPLICATION_NAME}/CN=localhost"
_tls/service.pem: SHELL:=/bin/bash
_tls/service.pem: _tls/service.csr _tls/rootca.cert _tls/rootca.key
	openssl x509 -req -sha256 -CA _tls/rootca.cert -CAkey _tls/rootca.key -CAcreateserial -in $< -out $@ -extfile <( \
		echo 'authorityKeyIdentifier=keyid,issuer'; \
		echo 'basicConstraints=CA:FALSE'; \
		echo 'keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment'; \
		echo 'subjectAltName = DNS:localhost,IP:127.0.0.1'; \
	)


.PHONY: tools-grpcui
tools-grpcui: | static-tlsdata
	grpcui \
		-import-path ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		-import-path ./protobuf \
		-cacert _tls/rootca.cert -key _tls/service.key -cert _tls/service.pem \
		-bind localhost -port 4444 localhost:50051
		# -proto favorites.proto \

.PHONY: tools-boltdbweb
tools-boltdbweb:
	boltdbweb --db-name=application.db --port=8089

.PHONY: tools-install
tools-install:
	# github.com/rakyll/statik -> github.com/goware/statik -> goware allow multiple assets modules
	GO111MODULE=off go get -u \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
		github.com/goware/statik \
		github.com/fullstorydev/grpcui/cmd/grpcui \
		github.com/evnix/boltdbweb
