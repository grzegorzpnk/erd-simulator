#!/bin/bash

make all

sudo docker tag simu pmatysiaq/simu:latest

sudo docker push pmatysiaq/simu:latest
