#!/bin/bash

dirs=("simulator" "erc" "nmt" "rl-agent")

for dir in "${dirs[@]}"; do
  cd ./pkg/"$dir" || exit
  ./deploy.sh
  cd ../..
done
