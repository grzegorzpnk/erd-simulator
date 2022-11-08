### MEC Host structure

1. Mec Host is represented as struct below.

```go
type MecHost struct {
	Identity   MecIdentity  `json:"identity"`
	Resources  MecResources `json:"resources,omitempty"`
	Neighbours []*MecHost
}
```

```go
const (
	// We will have 3 types of MEC Hosts 
	MecLocal    MecType = iota
	MecRegional
	MecCentral
)

// MecIdentity represents single Mec Host
type MecIdentity struct {
	Provider string      `json:"provider"`
	Cluster  string      `json:"cluster"`
	Location MecLocation `json:"location"`
}

// MecResInfo contains information about MEC Host (eg. Cpu, Memory)
type MecResInfo struct {
	Used        float64 `json:"used"`        // How many cpu/memory used (value)
	Allocatable float64 `json:"allocatable"` // How many cpu/memory available (value)
	Utilization float64 `json:"utilization"` // How much is cpu/memory utilized (percentage)
}

type MecResources struct {
	Latency float64    `json:"latency"`
	Cpu     MecResInfo `json:"cpu"`
	Memory  MecResInfo `json:"memory"`
}

// Type is necessary for the algorithm. Region is needed to specify if specific MEC neighbours should
// be taken into account as candidate MEC Hosts. TODO: Please consider different names
type MecLocation struct {
	Type      MecType `json:"type"`
	Region    string  `json:"region"`               // eg. poland
	Zone      string  `json:"zone,omitempty"`       // eg. west, "" if type different from MecLocal and MecRegional
	LocalZone string  `json:"local-zone,omitempty"` // eg. wroclaw, "" if type different from MecLocal
}
```

2. When we call for `list of MEC Hosts` or `list of MEC Hosts neighbours`, Topology server should return a list of []MecIdentity
    and then we can fetch other information: `MEC Resources` or `MEC Neighbours`.

```go
var mecList = []MecIdentity{
{
    Provider: "provider1",
    Cluster:  "mec1",
    Location: MecLocation{
        Type:      MecLocal,
        Region:    "poland",
        Zone:      "west",
        LocalZone: "wroclaw",
        },
    },
    {
    Provider: "provider1",
    Cluster:  "mec2",
    Location: MecLocation{
        Type:      MecLocal,
        Region:    "poland",
        Zone:      "west",
        LocalZone: "katowice",
        },
    },
}
```
3. MEC `latency`, `cpu` & `memory` should be collected independently.

 ~~Note that there are two types of latency: between `MEC Host` and `targetCell` & between `two MEC Hosts`.~~
Latency between MecHosts isn't important -> we will check the shortest path instead and consider as E2E latency.

5. MEC `Neighbours` should be collected independently

### Topology server Endpoints

If we suppose that in the config file we provide topology endpoint like:

```http request
http://localhost:8787/v2
```

Then we can introduce few endpoints

a. [GET] `MEC Hosts` (MecIdentity) list associated with given `cell-id` (as []MecIdentity)

```http request
/topology/cell/{cell-id}/mec-hosts
```

~~b. [GET] `Latency` between given `cell-id` and `MEC Hosts`~~

```go
// eg. cell-id = 00000020; mec-name=provider1+mec1
```

```http request
/topology/cell/{cell-id}/mec/{mec-name}/latency
```

~~c. [GET] `Latency` between given MEC Hosts and his neighbour~~

```go
// eg. mec-name=provider1+mec1; mec-neighbour=provider1+mec2
```

```http request
/topology/mec/{mec-name}/neighbour/{mec-neighbour}/latency
```

d. [GET] `CPU` struct for given `MEC Hosts` (used, allocatable, utilization)

~~Single endpoint for all CPU information~~ (or consider endpoint for each: used, allocatable, utilization):

```go
// eg. mec-name=provider1+mec1
```

```http request
/topology/mec/{mec-name}/cpu
```

e. [GET] `MEMORY` struct for given `MEC Host` (used, allocatable, utilization)

~~Single endpoint for all MEMORY information~~ (or consider endpoint for each: used, allocatable, utilization):

```go
// eg. mec-name=provider1+mec1
```

```http request
/topology/mec/{mec-name}/memory
```

f. [GET] neighbour list of given `MEC Host`

Returned neighbours suppose to be a list of []MecIdentity

```go
// eg. mec-name=provider1+mec1
```

```http request
/topology/mec/{mec-name}/neighbours
```

g. [GET] Shortest Path between two MEC Hosts as latency between `start-mec` and `end-mec`

Returns value of The Shortest Path which represents E2E latency

```http request
/topology/start/{start-node}/stop/{stop-node}/shortest-path
```

*Note: we always consider cell <-> mec path?? or we also consider mec <-> mec path?*