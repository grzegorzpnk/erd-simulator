#!/bin/bash

generate_apps=''
generate_clusters=''
deploy=''
uninstall=''

print_usage() {
  echo "Usage:"
  echo "./03-DEPLOY_EDGE_APPS.sh [options]"
  echo "options:"
  echo "  -a        generate new applications (values and helm charts)"
  echo "  -c        generate new target clusters for applications"
  echo "  -d        deploy edge applications"
  echo "  -u        uninstall edge applications"
}

while getopts 'acdhu' flag; do
  case "${flag}" in
    a) generate_apps=true ;;
    c) generate_clusters=true ;;
    d) deploy=true ;;
    u) uninstall=true ;;
    h) print_usage ;;

    *) print_usage
       exit 1 ;;
  esac
done

cd ./pkg/lcm-workflow/samples/lcm-example || exit

if [[ $generate_apps ]]; then
  ./01-GENERATE_APPS_VALUES.sh
  echo "GENERATED APPS & HELM CHARTS!"
fi

if [[ $generate_clusters ]]; then
  ./02-GENERATE_TARGET_CLUSTERS.sh
  echo "GENERATED NEW TARGET CLUSTERS!"
fi

if [[ ! $generate_apps ]] && [[ ! $generate_clusters ]]; then
  print_usage
fi

cd ./emco-manifests-v2 || exit

if [[ $uninstall ]]; then
  ./00-CLEAR_UP.sh
  echo "APPLICATIONS UNINSTALLED!"
  exit 1
fi

if [[ $deploy ]]; then
  ./01-PREPARE_ENV.sh
  sleep 2
  ./02-DEPLOY_APPS.sh
  echo "APPLICATIONS DEPLOYED!"
fi

cd ../../../../../