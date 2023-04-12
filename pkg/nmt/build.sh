#!/bin/bash

make all

sudo docker tag nmt pmatysiaq/nmt:latest

sudo docker push pmatysiaq/nmt:latest
