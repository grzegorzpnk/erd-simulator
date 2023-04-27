#!/bin/bash

dirs=("rl-agent" )

for dir in "${dirs[@]}"; do
  cd ./pkg/"$dir" || exit
  ./build.sh
  cd ../..
done
