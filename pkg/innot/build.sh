#!/bin/bash

make all

sudo docker tag innot pmatysiaq/innot:latest

sudo docker push pmatysiaq/innot:latest
