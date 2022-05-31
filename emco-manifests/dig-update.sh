#!/bin/bash

curl --location --request POST 'http://10.254.185.70:30415/v2/projects/towards5gs/composite-apps/free5gc-ca/v1/deployment-intent-groups/free5gc-deployment-intent/update' \
--header 'Content-Type;' \
--data-raw '""'
