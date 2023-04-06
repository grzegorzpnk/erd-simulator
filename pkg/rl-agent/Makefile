.PHONY: all build

all: clean compile build

build:
	@echo "Building rl-agent app"
	sudo docker build --rm -f ./Dockerfile -t rl-agent .

compile:
	@echo "No compilation required for python"

clean:
	@echo "Nothing to clean"

test:
	@echo "No tests defined"
