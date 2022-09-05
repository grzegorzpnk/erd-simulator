#!/bin/bash

make all

sudo docker tag obs pmatysiaq/obs:latest

sudo docker push pmatysiaq/obs:latest
