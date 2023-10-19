#!/bin/bash

dirs=("erc" "nmt" "simulator" "rl-agent")

for dir in "${dirs[@]}"; do
  cd ./pkg/"$dir" || exit
  ./build.sh
  cd ../..
done
