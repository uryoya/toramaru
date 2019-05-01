GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINNAME=toramaru

all: test build

build:
	$(GOBUILD) -o $(BINNAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -fr $(BINNAME)

