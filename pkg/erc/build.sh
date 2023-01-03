#!/bin/bash

make all

sudo docker tag erc grzegorzpnk/erc:latest

sudo docker push grzegorzpnk/erc:latest
