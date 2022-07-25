.PHONY: all build

all: clean compile build

build:
	@echo "Building innot app"
	sudo docker build --rm -f ./build/Dockerfile -t innot .

compile:
	@echo "Compiling innot app"
	/bin/mkdir -p bin; \
	cd src/cmd; go build -o ../../bin/innot main.go; cd ../../;

clean:
	@echo "Cleaning..."
	/bin/rm -rf bin/

test:
	@echo "No tests defined"
