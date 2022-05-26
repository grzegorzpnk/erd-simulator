# Usage: ./test-timing.sh $context1 $context2 $context3

isDeploymentReady() {
  echo $1

  while true; do
      count=0
      for status in $(kubectl --context=$1 get pod --output=jsonpath={.items..phase}); do
          if [ $status != "Running" ]
          then
            count=$((count+1))
          fi  
      done
    if [ $count -gt 0 ]; then
      #printf "%d pods not Ready; Context: %s", $count, $1
      continue
    else
      return
    fi
  done
}

isUPFReady() {
    for pod in $(kubectl --context=$1 get pod --output=jsonpath={.items..metadata.name}); do
        if [[ $pod == *"upf"* ]]; then
            status=$(kubectl --context=$1 get pod $pod --output=jsonpath={.status.phase})
            if [[ $status == "Running" ]]; then
                return 0 # return true
            else
                return 1 # return false
            fi
        else
            continue
        fi
    done
    return 1 # return false
}

isUPFGone() {
    for pod in $(kubectl --context=$1 get pod --output=jsonpath={.items..metadata.name}); do
        if [[ $pod == *"upf"* ]]; then
            return 1 # return false
        fi
    done
    return 0 # return true
}

printf "\n**Creating EMCO resources...**\n"

startTime=$(date +"%s%3N")

emcoctl --config ./emco-cfg.yaml apply -f 1-prerequisites.yaml -v values.yaml
emcoctl --config ./emco-cfg.yaml apply -f 2-instantiate-logical-cloud.yaml -v values.yaml
sleep 1 # wait until logical cloud is instantiated
emcoctl --config ./emco-cfg.yaml apply -f 3-deployment.yaml -v values.yaml

startDigInstantiateTime=$(date +"%s%3N")
printf "startDigInstantiateTime: %s\n" $startDigInstantiateTime

emcoctl --config ./emco-cfg.yaml apply -f 4-instantiate-free5gc.yaml -v values.yaml
emcoctl --config ./emco-cfg.yaml apply -f 5-instantiate-ueransim.yaml -v values.yaml

sleep 3 # wait until pods are created

printf "**Check if deployment is ready**\n"
isDeploymentReady $1
isDeploymentReady $2

endDigInstantiateTime=$(date +"%s%3N")
DIG_INSTANTIATE_TIME=$((endDigInstantiateTime-startDigInstantiateTime))

printf "**Update Placement Intent: Migrate UPF don't delete old instance**\n"
./first-pi-update.sh

startDigUpdate1Time=$(date +"%s%3N") # call update dig

./dig-update.sh

until isUPFReady $3; do
    sleep 0.1
done

endDigUpdate1Time=$(date +"%s%3N")
DIG_UPDATE1_TIME=$((endDigUpdate1Time-startDigUpdate1Time))

printf "**Update Placement Intent: Delete redundant UPF**\n"
./second-pi-update.sh

startDigUpdate2Time=$(date +"%s%3N")

./dig-update.sh

until isUPFGone $2; do
    sleep 0.1
done

endDigUpdate2Time=$(date +"%s%3N")
DIG_UPDATE2_TIME=$((endDigUpdate2Time-startDigUpdate2Time))

TEST_TIME=$((endDigUpdate2Time-startTime))

result='{"instantiate-time":"'"$DIG_INSTANTIATE_TIME"'", "update-1-time":"'"$DIG_UPDATE1_TIME"'", "update-2-time":"'"$DIG_UPDATE2_TIME"'", "test-time":"'"$TEST_TIME"'"}'

printf "%s\n" $result >> result.log
