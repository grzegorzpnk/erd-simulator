.PHONY: all build

all: clean compile build

build:
	@echo "Building pmc app"
	sudo docker build --rm -f ./build/Dockerfile -t pmc .

compile:
	@echo "Compiling pmc app"
	/bin/mkdir -p bin; \
	cd src/cmd; go build -o ../../bin/pmc main.go; cd ../../;

clean:
	@echo "Cleaning..."
	/bin/rm -rf bin/

test:
	@echo "No tests defined"
