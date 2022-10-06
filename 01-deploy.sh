#!/bin/bash

dirs=("erc" "innot" "lcm-workflow" "relocate-workflow" "obs")

# shellcheck disable=SC2068
for dir in ${dirs[@]}; do
  cd pkg/"$dir"
  ./deploy.sh
  cd ../..
done