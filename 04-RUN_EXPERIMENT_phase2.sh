#!/bin/bash

INNOT_ENDPOINT=10.254.185.44
INNOT_PORT=32137
NUMBER_OF_MOVEMENTS=50

run() {
  declare -a notifyList=()
  declare -a temp=()

  if [ -f "$(pwd)/movements-v3.0.$1/notification_order_$2" ]; then
    { read -a notifyList; } < ./movements-v3.0."$1"/notification_order_"$2"
  fi


  if [[ ${#notifyList[@]} == "$NUMBER_OF_MOVEMENTS" ]];
  then
    printf "Notification order fetched from the file!\n"
  else
    printf "Notification order will be generated\n"

    while [[ ${#notifyList[@]} != "$NUMBER_OF_MOVEMENTS" ]];
    do
      subId=$(expr $((1 + $RANDOM % 50)))
      notifyList+=("$subId")
    done

    if [ ! -d "$(pwd)/movements-v3.0.$1" ]; then
        mkdir ./movements-v3.0."$1"
    fi

    echo "${notifyList[*]}" > ./movements-v3.0."$1"/notification_order_"$2"
  fi

  for subId in "${notifyList[@]}";
  do
      curl -X POST -d {} http://$INNOT_ENDPOINT:$INNOT_PORT/v1/intermediate-notifier/subscriptions/$subId/handle
#      echo "Sent notification to subscriber[$subId]"
      sleep 15
  done
}

print_usage() {
  echo "Usage:"
  echo "./04-RUN_EXPERIMENT_phase2.sh [options]"
  echo "options:"
  echo "  -r        which run of experiment (e.g r=1 is 1st run of experiments)"
  echo "  -i        number of iterations"
}

while getopts 'r:i:h' flag; do
  case "${flag}" in
    r) runNo=$OPTARG ;;
    i) iterationsNo=$OPTARG ;;

    h) print_usage ;;

    *) print_usage
       exit 1 ;;
  esac
done

if [ -z "$runNo" ] || [ -z "$iterationsNo" ]; then
        echo 'Missing -r or -i' >&2
        print_usage
        exit 1
fi


echo "----- Run No. $runNo in progress... -----"
for i in $(seq $iterationsNo);
do
  WAIT_TIME=$(expr 15 \* "$NUMBER_OF_MOVEMENTS" / 60)
  echo "----- Iteration No. $i in progress... It will take [$WAIT_TIME minutes] in total -----"


  run "$runNo" "$i"

  echo "---- RESULTS of Iteration[$i]-----"
  sleep 15
  curl -s -X GET http://10.254.185.44:32139/v1/topology/mecHosts/metrics | json_pp -json_opt pretty
  curl -s -X GET http://10.254.185.44:32137/v1/results/csv | json_pp -json_opt pretty
  curl -X POST http://10.254.185.44:32137/v1/results/reset
done
