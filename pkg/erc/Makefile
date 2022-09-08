.PHONY: all build

all: clean compile build

build:
	@echo "Building erc app"
	sudo docker build --rm -f ./build/docker/Dockerfile -t erc .

compile:
	@echo "Compiling erc app"
	/bin/mkdir -p bin; \
	cd src/cmd; go build -o ../../bin/erc main.go; cd ../../;

clean:
	@echo "Cleaning..."
	/bin/rm -rf bin/

test:
	@echo "No tests defined"
