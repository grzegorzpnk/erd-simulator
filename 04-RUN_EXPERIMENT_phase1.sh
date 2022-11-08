#!/bin/bash

INNOT_ENDPOINT=10.254.185.44
INNOT_PORT=32137

declare -a notifyList=()
declare -a temp=()

{ read -a notifyList; } < notification_order_1

if [[ ${#notifyList[@]} == "50" ]];
then
  printf "Notification order fetched from the file!\n"
else
  printf "Notification order will be generated\n"

  while [[ ${#notifyList[@]} != "50" ]];
  do
    subId=$(expr $((1 + $RANDOM % 50)))
    notifyList+=("$subId")
#    # shellcheck disable=SC2076
#    if [[ ! " ${notifyList[*]} " =~ " ${subId} " ]]; then
#      notifyList+=("$subId")
#    fi
  done

  echo "${notifyList[*]}" > notification_order_1
fi

for subId in "${notifyList[@]}";
do
    curl -X POST -d {} http://$INNOT_ENDPOINT:$INNOT_PORT/v1/intermediate-notifier/subscriptions/$subId/handle
    echo "Sent notification to subscriber[$subId]"
    echo "Sleep 15s"
    sleep 15
done
