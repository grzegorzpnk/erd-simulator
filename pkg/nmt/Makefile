.PHONY: all build

all: clean compile build

build:
	@echo "Building nmt app"
	sudo docker build --rm -f ./build/docker/Dockerfile -t nmt .

compile:
	@echo "Compiling nmt app"
	/bin/mkdir -p bin; \
	cd src/cmd; go build -o ../../bin/nmt main.go; cd ../../;

clean:
	@echo "Cleaning..."
	/bin/rm -rf bin/

test:
	@echo "No tests defined"
