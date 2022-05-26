for i in {1..20}; do
    echo "----------loop $i----------"
    ./test-timing.sh $1 $2 $3
    ./clear-up.sh
    sleep 50	# wait until all pods are deleted
done
