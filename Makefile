VERSION = 0.6

GO_FMT = gofmt -s -w -l .
GO_XC = goxc -os="linux darwin windows"

GOXC_FILE = .goxc.local.json

all: deps compile

compile: goxc

goxc:
	$(shell echo '{\n "ConfigVersion": "0.9",\n "PackageVersion": "$(VERSION)",' > $(GOXC_FILE))
	$(shell echo ' "TaskSettings": {' >> $(GOXC_FILE))
	$(shell echo '  "bintray": {\n   "apikey": "$(BINTRAY_APIKEY)"' >> $(GOXC_FILE))
	$(shell echo '  }\n } \n}' >> $(GOXC_FILE))
	$(GO_XC) 

deps:
	go get

format: 
	$(GO_FMT) 

bintray:
	$(GO_XC) bintray
