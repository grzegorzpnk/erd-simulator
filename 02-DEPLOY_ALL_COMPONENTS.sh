#!/bin/bash

dirs=("erc" "innot" "lcm-workflow" "relocate-workflow" "obs" "nmt")

for dir in "${dirs[@]}"; do
  cd ./pkg/"$dir" || exit
  ./deploy.sh
  cd ../..
done
