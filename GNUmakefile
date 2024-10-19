VERSION=0.0.1
default: build

clean:
	rm -rf $(CURDIR)/build

install: clean 
	install -d $(CURDIR)/build/darwin_arm64 && install -d $(CURDIR)/build/linux_amd64

build: install
	go build -C $(CURDIR)/cmd/load-test-rds -o $(CURDIR)/build/darwin_arm64
	export GOOS=linux && export GOARCH=amd64 && go build -C $(CURDIR)/cmd/load-test-rds -o $(CURDIR)/build/linux_amd64