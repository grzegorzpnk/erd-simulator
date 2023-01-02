#!/bin/bash

if [ -z "${ERD_REPO_NAME}" ] || [ -z "${ERD_REPO_TAG}" ]; then
  echo "\"ERD_REPO_NAME\" and \"ERD_REPO_TAG\" environment variables must be set!"
  exit 1
fi

make all

sudo docker tag erc "$ERD_REPO_NAME"/erc:"$ERD_REPO_TAG"

sudo docker push "$ERD_REPO_NAME"/erc:"$ERD_REPO_TAG"

