#!/bin/bash

make all

sudo docker tag nmt grzegorzpnk/nmt:latest

sudo docker push grzegorzpnk/nmt:latest
