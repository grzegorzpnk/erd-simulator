#!/bin/bash

dirs=("erc" "nmt" "simulator")

for dir in "${dirs[@]}"; do
  cd ./pkg/"$dir" || exit
  ./deploy.sh
  cd ../..
done
