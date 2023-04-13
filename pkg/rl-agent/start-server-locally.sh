#!/usr/bin/env sh

export GET_PREDICTION_SCHEMA_PATH=$(pwd)/get-prediction-schema.json
export CONFI_PATH=$(pwd)/config.json

python ./rlagent/manage.py runserver 8080
