VERSION = 0.9.1

GO_FMT = gofmt -s -w -l .
GO_XC = goxc -os="linux darwin windows" -tasks-="rmbin"

GOXC_FILE = .goxc.local.json

all: deps compile

compile: goxc

goxc:
	$(shell echo '{\n "ConfigVersion": "0.9",\n "PackageVersion": "$(VERSION)",' > $(GOXC_FILE))
	$(shell echo ' "TaskSettings": {' >> $(GOXC_FILE))
	$(shell echo '  "bintray": {\n   "apikey": "$(BINTRAY_APIKEY)"' >> $(GOXC_FILE))
	$(shell echo '  },' >> $(GOXC_FILE))
	$(shell echo '  "publish-github": {' >> $(GOXC_FILE))
	$(shell echo '     "apikey": "$(GITHUB_APIKEY)",' >> $(GOXC_FILE))
	$(shell echo '     "body": "",' >> $(GOXC_FILE))
	$(shell echo '     "include": "*.tar.gz,*.deb,depcon-linux64,depcon-osx64,depcon-win64.exe"' >> $(GOXC_FILE))
	$(shell echo '  }\n } \n}' >> $(GOXC_FILE))
	$(GO_XC)
	cp build/$(VERSION)/linux_amd64/depcon build/$(VERSION)/depcon-linux64
	cp build/$(VERSION)/darwin_amd64/depcon build/$(VERSION)/depcon-osx64
	cp build/$(VERSION)/windows_amd64/depcon.exe build/$(VERSION)/depcon-win64.exe

deps:
	go get

format:
	$(GO_FMT)

bintray:
	$(GO_XC) bintray

github:
	$(GO_XC) publish-github

docker-build:
	cp build/$(VERSION)/linux_amd64/depcon docker-release/depcon
	docker build -t containx/depcon docker-release/
	docker tag containx/depcon containx/depcon:$(VERSION)

docker-push:
	docker push containx/depcon
	docker push containx/depcon:$(VERSION)
