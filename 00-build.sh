#!/bin/bash

# dirs=("erc" "innot" "lcm-workflow" "relocate-workflow" "obs")
dirs=("erc" "innot" "lcm-workflow" "obs" "nmt")

# shellcheck disable=SC2068
for dir in ${dirs[@]}; do
  cd ./pkg/"$dir"
  ./build.sh
  cd ../..
done
