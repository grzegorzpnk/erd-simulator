#!/bin/bash

dirs=("simulator" "erc" "nmt")

for dir in "${dirs[@]}"; do
  cd ./pkg/"$dir" || exit
  ./deploy.sh
  cd ../..
done
