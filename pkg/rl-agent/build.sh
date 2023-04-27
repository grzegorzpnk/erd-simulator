#!/bin/bash

make all

sudo docker tag rl-agent grzegorzpnk/rl-agent:latest

sudo docker push grzegorzpnk/rl-agent:latest
