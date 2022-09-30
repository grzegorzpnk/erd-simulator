## Innot (Intermediate Notifier)

This project is designed for PoC purposes. It allows to subscribe to the AMF Notifications.

---

## Design

#### `Innot` suppose to be an intermediate notifier between AMF and MEO. Due to the fact that subscribing to the Free5gc AMF notifications is not a real subscription, but rather single notification, `Innot` can emulate subscription to the AMF.

`Innot` receives AMF subscription endpoint during startup and exposes REST API for subscription. It asks AMF for the
notification once per second and compares `selected*` fields. If `specfic` field did change, it sends notification to the
MEO `eventNotifyUri` endpoint (which is provided as part of `HTTP POST Body`)
  - For now `LOCATION_REPORT` is supported, which checks if `CELL_ID` changed.
  - MEO should listen for notifications on `eventNotifyUti`
  - It will be discussed, what `Notification Body` should contain.

For now `LOCATION_REPORT` notification contains `Reason` and `New CELL_ID`
```go
# ...
type NotifyReason string
type CellId string

const (
	CELLCHANGED NotifyReason = "CELL_ID_CHANGED"
)

type CellChangedInfo struct {
	Reason NotifyReason
	Cell   CellId
}

```


`NOTE:` Please adjust `AMF IP and Port` in the `config.json` file or in the `values.yaml` while using helm chart.

`selected*` means, that it should be defined in advanced, which fields will be compared

---

## `Innot` endpoints

```yaml
# Subscribe for the notifications. Sample body is provided below.
METHOD POST: 10.254.185.45:32137/v1/intermediate-notifier/subscribe
```

```yaml
# Unsubscribe (Not supported yet)
# METHOD POST: 10.254.185.45:32137/v1/intermediate-notifier/unsubscribe
```

```yaml
# Unsubscribe based on subscriptionId (Not supported yet)
# METHOD POST: 10.254.185.45:32137/v1/intermediate-notifier/unsubscribe/{subscriptionId}
```

```yaml
# Get all subscriptions (Not supported yet)
# METHOD GET: 10.254.185.45:32137/v1/intermediate-notifier/get-all
```

## Sample Body Request

```json
{
   "subscription":{
      "eventList":[
         {
            "type":"LOCATION_REPORT",
            "immediateFlag":true,
            "areaList":[
               {
                  "presenceInfo":{
                     "trackingAreaList":[
                        {
                           "plmnId":{
                              "mcc":"208",
                              "mnc":"93"
                           },
                           "tac":"000001"
                        }
                     ]
                  },
                  "sNssai":[
                     {
                        "sd":"010203",
                        "sst":1
                     },
                     {
                        "sd":"112233",
                        "sst":1
                     }
                  ]
               }
            ],
            "udmDetectInd":true,
            "presenceInfoList":{
               
            },
            "maxResponseTime":0,
            "targetArea":{
               "taList":[
                  {
                     "plmnId":{
                        "mcc":"208",
                        "mnc":"93"
                     },
                     "tac":"1"
                  }
               ]
            }
         }
      ],
      "anyUE":true,
      "eventNotifyUri":"http://localhost/workflow-listener/cell-changed/ABCDEFGHIJ/notify",
      "options":{
         "trigger":"ONE_TIME"
      }
   }
}
```
This request body should be sent to the `Innot` and it will be used to receive notifications from AMF.

```json
"eventNotifyUri":"http://workflow-listener-ip:8181/v1/eventNotifier/notify"
```

## Sample configs provided to the `Innot`

```json
{
  "amf-endpoint": "http://10.254.185.44:31845/namf-evts/v1/subscriptions",
  "plugin-dir":"cwd",
  "service-port":"8181"
}
```

```json
{
  "amf-endpoint": "http://amf-namf:80/namf-evts/v1/subscriptions",
  "plugin-dir":"cwd",
  "service-port":"8181"
}
```

`amf-endpoint` is an endpoint to which `Innot` will send subscription requests

---
---
---

## Motivations to create `Innot` - ER WG considerations

### Note that an application may be composed of multiple microservices defined together in a Helm chart. The application may be part of a composite application, whose component applications may be running in different clusters, including public/private clouds, telco/co-located edge clusters, and on-prem data centers.

- Existing Relocation Workflow doesn't allow to place different microservices in different locations. We only can relocate applications to the single cluster.
    - Should we consider specifying target cluster for each single microservice?
    - How we can provide microservice re-configuration in a generic way? (we need to assure that each microservice will know how to reach another microservices, even
      if they are distributed across many k8s clusters)

- Existing Relocation Workflow doesn't allow to place applications based on cluster-labels.
    - Even if application is placed using label, after relocation it will be placed based on clusterName
    - EMCO provides additional challanges due to the fact that we can't specify single cluster twice (once as label and second as clusterName)

- Existing Relocation Workflow always will relocate application to the single cluster, even if application was deployed on many different clusters.
    - Should we consider specifying a list of clusters where the specific application (or microservice) should be placed?

### As an important corollary, the new instance of the application must be declared to be 'ready' only when all its associated state has been relocated.

- Existing Relocation Workflow rediness check is done using monitor, which definately is not stable.
    - Can we rely on the informations provided by the monitor?
    - Maybe we should wait additional time before moving to the TrafficSteering activity, after monitor notifies us that application is "Ready"?
    - (future-work) How to make sure that the state was relocated? (it depends how the state relocation will be implemented)


### The state may include not only databases and other configuration, but also configuration of networking policies (such as firewall rules or security policies) in the target environment. Otherwise, service continuity cannot be assured.

- How we can cover all different external requirements? We don't know if there is any firewall, etc.


### A. Listen for the events that can trigger a relocation, e.g., notifications from the 5G core regarding user mobility, a new user or PDU session, or application performance issues.

#### From the EMCO perspective there could be an API endpoint defined, which would listen for the different events.
- http://workflow-listener-ip:8181/v1/eventNotifier/notify # And metadata is fetched from HTTP POST Body. Metadata can be fully filled or it can be partially empty
- metadata should be a struct similiar to the RelocationIntent because we have to define RelocationIntent after notification was received
- temporal-workflow-intent is to complex? How to make it simplier
- Some information can be filled or empty (e.g. target provider/cluster)

     - This approach is extensible, because we can define other eventTypes and eventHandlers afterwards. 
        - Example EventTypes -> relocateStatelessApp, relocateStatefullApp /// OR CELL_CHANGED

#### Other entities should be aware about this endpoint

- Based on AMF Notification OpenAPI, we can assume that we can provide such endpoint to the 5GC.
- How to make use of notification body? It won't contain information about relocation-workflow-intent.
- (free5gc) We can't (?) subscribe to the notifications, we can only request single notification -> For PoC purposes there could exist some external entity in the middle, which will ask for notifications and
if necessary invoke eventListener endpoint.

#### If we consider such endpoint where it will be defined?

- Higher level workflow? -> can such endpoint be defined inside some workflow?
- RelocatioWorkflow, which will be running all the time (in a loop), not only when relocation is needed?
- External entity? -> If we consider that, it's not usefull solution for EMCO community.
- Inside EMCO? how?

#### Namf_EventExposure allows to subscribe to the: 

LOCATION_REPORT, PRESENCE_IN_AOI_REPORT, TIMEZONE_REPORT, ACCESS_TYPE_REPORT, REGISTRATION_STATE_REPORT,
CONNECTIVITY_STATE_REPORT, REACHABILITY_REPORT, COMMUNICATION_FAILURE_REPORT, UES_IN_AREA_REPORT, SUBSCRIPTION_ID_CHANGE, SUBSCRIPTION_ID_ADDITION,
SUBSCRIPTION_ID_ADDITION notifications

#### There are location filters such as: TAI, CELL_ID, N3IWF, UE_IP, UDP_PORT

- For example we could use LOCATION_REPORT which returns current cell id.
    - When handover is performed, the cell ID will be changed. Of course, we would prefere to receive notification before HO, but how to achieve that?
    - For PoC purposes we could have some external entity, which would fetch notifications from AMF and compare cell ID. When the cell changes, this external entity could
      notify eventListener endpoint to invoke the relocation.

---
---

## Further considerations

### B. Determine whether to relocate the application.

#### When application shouldn't be relocated, after relocation notification was received?
- Application doesn't support relocation -> How to define if application supports the relocation? On the K8s level labels (allowRelocation=true) but can EMCO use this information?
- Situation that there is no better cluster should be considered in the step C.
- Others?

### C. Determine the ‘best’ target MEC cluster based on many criteria, as stated in the problem description.

- Let's say that we allow to relocate different microservices to different clusters. Do we "determine best T-MEH" for each microservice?
- There should be some controller which is aware of MEC topology and links between them?
- 5GFF EDS?

### D. Perform the relocation

- RelocationWorkflow implementation should be further discussed, due to earlier considerations

