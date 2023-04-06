#!/bin/bash

make all

sudo docker tag rl-agent pmatysiaq/rl-agent:latest

sudo docker push pmatysiaq/rl-agent:latest
