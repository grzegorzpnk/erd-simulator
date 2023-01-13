.PHONY: all build

all: clean compile build

build:
	@echo "Building simu app"
	sudo docker build --rm -f ./build/docker/Dockerfile -t simu .

compile:
	@echo "Compiling simu app"
	/bin/mkdir -p bin; \
	cd src/cmd; go build -o ../../bin/simu main.go; cd ../../;

clean:
	@echo "Cleaning..."
	/bin/rm -rf bin/

test:
	@echo "No tests defined"
