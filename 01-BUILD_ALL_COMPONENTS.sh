#!/bin/bash

dirs=("erc" "innot" "lcm-workflow" "obs" "nmt")

for dir in "${dirs[@]}"; do
  cd ./pkg/"$dir" || exit
  ./build.sh
  cd ../..
done
