{CATCHY NAME} Framework
---

Introduction
---

Components
---

### Edge Relocation Controller (ERC)

#### Overview

ERC is designed to play a role of Placement Controller, which for now supports only one type on Intent (`Smart Placement Intent`), and
based on that can interact with NMT component to find the best cluster based on user location (based on CELL_ID). When we 
say best cluster, we mean in terms of avaliable resources (CPU, MEMORY requests) and e2e latency.

#### Smart Placement Intent



#### APIs Endpoints

### LCM-Workflow

`LCM-Workflow` is designed as Temporal Workflow. It is composed of several `Temporal Activities`, but is easy extensible.
`LCM-Workflow` makes use of `EMCO orchestrator` and it's designed to be executed as a POST-INSTALL Hook (EMCO TAC Intent), 
which means that EMCO automatically starts this Workflow after the application is instantiated.

#### Activities:

- `SubCellChangedNotification` - activity which is run as the first one. For now, it is responsible to subscribe for `LOCATION_REPORT`
notifications from AMF. But can be extended to subscribe other types of notifications.
- `GetCellChangedNotification` - activity in which HTTP listener is started until mobility notification is sent by AMF.
- `DiscoverCurrentCluster` - activity which allows to collect information (from EMCO) about current placement of given app.
    This information is saved and sent in the `Smart Placement Intent (SPI)`.
- `GenerateSmartPlacementIntent` - activity which generates `SPI` based on: data provided in the `LCM Workflow Intent`, new
    CELL_ID received from AMF in mobility notifications and current placement.
- `CallPlacementController` - in this activity the `SPI` is sent to ERC. If Placement Controller returns valid cluster
  (in form supported by EMCO `Provider+Cluster`) then `Relocate Workflow Intent` is generated, if not we go back to the 
    `GetCellChangedNotification` activity, until the next notification arrives.
- `GenerateRelocateWfIntent` - activity which generates `Relocate Workflow Intent` based on: data provided in the `LCM Workflow Intent`,
    and new cluster info provided by the `ERC`.
- `CallTemporalWfController` - based on the generated `Relocate Workflow Intent` the `Relocate Workflow` is created and then started.
    After completion, we go back to the `GetCellChangedNotification` activity.

#### Known Issues:

- When the new `LCM-Workflow` is started, the pair of ports (PORT, NODE_PORT) is selected and reserved exclusively for this
Workflow. It's related to the design of how the notifications are received - for each `LCM-Workflow` separate HTTP listener is run.
Due to that, a set of 200 ports are exposed as NodePorts when the LCM-Workflow Worker is run.
- In the current design, the listener is terminated when the notification arrives, so it's possible that some notifications will
be rejected.

### Intermediate Notifier (Innot)

### Network + MEC Topology (NMT)

### Observability Controller (OBS)

### Relocate-Workflow

### (optional) free5gc

### (optional) Ueransim 

Prerequisites
---

Deployment
---

etc...
---