VERSION=0.0.1
default: build
BUILD_TARGET = build/darwin_arm64
BUILD_TARGET_LINUX = build/linux_amd64

clean:
	rm -rf $(CURDIR)/build

install: clean 
	install -d $(CURDIR)/$(BUILD_TARGET) && install -d $(CURDIR)/$(BUILD_TARGET_LINUX)

build: install
	go build -C $(CURDIR)/cmd/load-test-rds -o $(CURDIR)/$(BUILD_TARGET)
	export GOOS=linux && export GOARCH=amd64 && go build -C $(CURDIR)/cmd/load-test-rds -o $(CURDIR)/$(BUILD_TARGET_LINUX)
	cp $(CURDIR)/scripts/startWriter.sh $(CURDIR)/$(BUILD_TARGET)
	cp $(CURDIR)/scripts/startWriter.sh $(CURDIR)/$(BUILD_TARGET_LINUX)