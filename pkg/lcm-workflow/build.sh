#!/bin/bash

make all

sudo docker tag lcm-workflow-worker pmatysiaq/lcm-workflow-worker:latest
sudo docker tag workflow-client pmatysiaq/lcm-workflow-client:latest

sudo docker push pmatysiaq/lcm-workflow-client:latest
sudo docker push pmatysiaq/lcm-workflow-worker:latest
