wget --no-check-certificate --quiet \
  --method PUT \
  --timeout=0 \
  --header 'Content-Type: application/json' \
  --body-data '{
    "metadata": {
        "name": "upf-placement-intent",
        "description": "free5gc upf placement intent"
    },
    "spec": {
        "app-name": "free5gc-upf",
        "intent": {
            "allOf": [
                {
                    "provider-name": "orange",
                    "cluster-label-name": "LabelEdge"
                },
                {
                    "provider-name": "orange",
                    "cluster-label-name": "LabelCloud"
                }

            ]
        }
    }
}' \
   'http://10.254.185.70:30415/v2/projects/towards5gs/composite-apps/free5gc-ca/v1/deployment-intent-groups/free5gc-deployment-intent/generic-placement-intents/towards5gs-generic-placement/app-intents/upf-placement-intent'
