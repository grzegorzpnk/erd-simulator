#!/bin/bash

make all

sudo docker tag simu grzegorzpnk/simu:latest

sudo docker push grzegorzpnk/simu:latest
