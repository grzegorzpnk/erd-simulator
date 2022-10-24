declare -a notifyList=()
declare -a temp=()

{ read -a notifyList; } < notification_order

if [[ ${#notifyList[@]} == "100" ]];
then
  printf "Notification order fetched from the file!\n"
else
  printf "Notification order will be generated\n"

  while [[ ${#notifyList[@]} != "100" ]];
  do
    subId=$(expr $((1 + $RANDOM % 100)))
    # shellcheck disable=SC2076
    if [[ ! " ${notifyList[*]} " =~ " ${subId} " ]]; then
      notifyList+=("$subId")
    fi
  done

  echo "${notifyList[*]}" > notification_order
fi

for subId in "${notifyList[@]}";
do
    curl -X POST -d {} http://10.254.185.44:32137/v1/intermediate-notifier/subscriptions/$subId/handle
    echo "Sent notification to subscriber[$subId]"
    echo "Sleep 15s"
    sleep 15
done
