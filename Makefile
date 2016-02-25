export GOPATH=${CURDIR}

prepare:
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...

install: assets
	go install cmds/ssh-execute
	go install cmds/ssh-http

dev: devassets
	go install cmds/ssh-execute
	go install cmds/ssh-http

assets: src/assets/
	PATH=${PATH}:${CURDIR}/../../bin go-bindata-assetfs -pkg communication $<...
	cp bindata_assetfs.go src/communication/bindata_assetfs.go

devassets: src/assets/
	PATH=${PATH}:${CURDIR}/../../bin go-bindata-assetfs -debug -pkg communication $<...
	cp bindata_assetfs.go src/communication/bindata_assetfs.go
