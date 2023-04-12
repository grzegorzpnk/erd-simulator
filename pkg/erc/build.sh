#!/bin/bash

make all

sudo docker tag erc pmatysiaq/erc:latest

sudo docker push pmatysiaq/erc:latest
