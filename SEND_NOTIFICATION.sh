for i in {1..100};
do
    subId=$(expr $((1 + $RANDOM % 100)))
    curl -X POST -d {} http://10.254.185.44:32137/v1/intermediate-notifier/subscriptions/$subId/handle
    echo "Sent notification to subscriber[$subId]"
    echo "Sleep 15s"
    sleep 15
done
