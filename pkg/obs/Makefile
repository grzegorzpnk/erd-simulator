.PHONY: all build

all: clean compile build

build:
	@echo "Building obs app"
	sudo docker build --rm -f ./build/Dockerfile -t obs .

compile:
	@echo "Compiling obs app"
	/bin/mkdir -p bin; \
	cd src/cmd; go build -o ../../bin/obs main.go; cd ../../;

clean:
	@echo "Cleaning..."
	/bin/rm -rf bin/

test:
	@echo "No tests defined"
