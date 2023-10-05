package module

import (
	"10.254.188.33/matyspi5/erd/pkg/erc/src/config"
	log "10.254.188.33/matyspi5/erd/pkg/erc/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/errs"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/model"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/topology"
	"github.com/pkg/errors"
	"math"
	"strconv"
)

// SmartPlacementIntentClient implements the SmartPlacementIntentManager.
// It will also be used to maintain some localized state (Not important, inherited from sample emco controller)
type SmartPlacementIntentClient struct {
	dbInfo DBInfo
}

// SearchParams is needed to find best MEC Host.
type SearchParams struct {
	currentMECs   []model.MecHost // MEC Hosts which are being checked to be candidate
	evalNeighMECs []model.MecHost // MECs from currentMECs which met latency. If needed, consider their neighbours later
	checkedMECs   []model.MecHost // Remember which clusters have already been checked
	candidateMECs []model.MecHost // If cluster meets all the requirements, save it as candidate MEC
}

// NewIntentClient creates and returns new SmartPlacementIntentClient
func NewIntentClient() *SmartPlacementIntentClient {
	return &SmartPlacementIntentClient{
		dbInfo: DBInfo{
			collection: "resources", // should remain the same
			tag:        "data",      // should remain the same
		},
	}
}

// SmartPlacementIntentManager a manager is an interface for exposing the client's functionalities.
type SmartPlacementIntentManager interface {
	ServeSmartPlacementIntentHeuristic(checkIfClause bool, intent model.SmartPlacementIntent) (model.MecHost, error)
	ServeSmartPlacementIntentOptimal(intent model.SmartPlacementIntent) (model.MecHost, error)
	ServeSmartPlacementIntentML(checkIfMasked bool, intent model.SmartPlacementIntent) (model.MecHost, error)
}

// ServeSmartPlacementIntentHeuristic based on SmartPlacementIntent tries to find and return the best MEC Host.
// This method uses Heuristic algorithm, with or without checking if cluster is good enough (checkIfClause).
// When checkIfClause is true, search algorithm is skipped each time the current MEC Host is "good enough" (meets all requirements)
// Returns best MEC Host (model.MecHost) if found, error otherwise.
func (i *SmartPlacementIntentClient) ServeSmartPlacementIntentHeuristic(checkIfClause bool, intent model.SmartPlacementIntent) (model.MecHost, error) {
	var err error
	var sp SearchParams
	var bestMec model.MecHost
	var checkFurther bool

	tc := topology.NewTopologyClient()
	tc.CurrentCell = intent.Spec.SmartPlacementIntentData.TargetCell

	log.Infof("Smart Placement Intent: %+v", intent)

	if checkIfClause {
		skip, err := checkIfCurrentClusterOk(*tc, intent)
		if err != nil {
			log.Warnf("Can't skip current cluster. Reason: %v", err)
		}

		if skip {
			return model.MecHost{}, errs.ERR_CLUSTER_OK
		}
	}

	// Get all MEC Hosts associated with given Cell ID
	sp.currentMECs, err = tc.GetMecHostsByCellId(tc.CurrentCell)

	// Check if any cluster is a valid candidate, if yes -> try to find optimal cluster
	sp, err = FindCandidates(tc, sp, intent)
	if err != nil {
		log.Errorf("Could not serve Smart Placement Intent: %v", err)
		return model.MecHost{}, err
	}

	// Try to find the optimal cluster.
	bestMec, err = FindOptimalCluster(sp.candidateMECs, intent)
	if err != nil {
		// save information about checked clusters -> no need to check single cluster twice
		sp.checkedMECs = append(sp.checkedMECs, sp.currentMECs...)
		sp.currentMECs = []model.MecHost{}
		sp.candidateMECs = []model.MecHost{}
		// flag which indicates, if we should search further
		//	* search until best cluster is found OR
		//	* until there are any clusters on the evalNeighMECs list
		checkFurther = true
	}

	for checkFurther {
		if len(sp.evalNeighMECs) == 0 {
			// If evalNeighMECs list is empty -> there are not any clusters to check (search failed).
			reason := "no valid neighbours to consider"
			err := errors.New(reason)
			log.Warnf("[RESULT] Could not find optimal cluster: %v", reason)
			return model.MecHost{}, err
		}

		log.Infof("EvalNeighList: %v", sp.evalNeighMECs)
		for _, mec := range sp.evalNeighMECs {
			mec, err = tc.GetMecNeighbours(mec)
			if err != nil {
				log.Warnf("Could not proceed %v neighbours. Reason: %v", mec.BuildClusterEmcoFQDN(), err)
				continue
			}

			for _, neigh := range mec.Neighbours {
				if skip := checkIfSkip(*neigh, sp.checkedMECs); skip {
					// If MEC Host already checked -> just skip
					continue
				}
				if skip := checkIfSkip(*neigh, sp.currentMECs); skip {
					// If MEC Host is in the current search space -> just skip
					continue
				}
				if skip := checkCoverageZone(*neigh, sp.checkedMECs[0]); skip {
					// If we are there, checkedMECs[0] exists. Skip neighbour if it's outside coverage zone [region]
					continue
				}
				sp.currentMECs = append(sp.currentMECs, *neigh)
			}
		}
		// Delete old clusters from evalNeighMECs list -> later add new if any exists
		sp.evalNeighMECs = []model.MecHost{}

		// Try to find new candidates, based on the new set of clusters
		sp, err = FindCandidates(tc, sp, intent)
		if err != nil {
			log.Errorf("Could not serve Smart Placement Intent: %v", err)
			return model.MecHost{}, err
		}

		// Try to find optimal cluster, based on the new set of clusters
		bestMec, err = FindOptimalCluster(sp.candidateMECs, intent)
		if err != nil {
			sp.checkedMECs = append(sp.checkedMECs, sp.currentMECs...)
			sp.currentMECs = []model.MecHost{}
			sp.candidateMECs = []model.MecHost{}
		} else {
			checkFurther = false
			return bestMec, err
		}
	}

	return bestMec, err
}

// ServeSmartPlacementIntentOptimal based on SmartPlacementIntent tries to find and return the best MEC Host.
// This method uses Optimal algorithm, it means that MEC Hosts from entire Coverage Zone are take into account at once.
// When Coverage Zone increases in size, this algorithm gets slow and may not be recommended for critical applications.
// Returns best MEC Host (model.MecHost) if found, error otherwise.
func (i *SmartPlacementIntentClient) ServeSmartPlacementIntentOptimal(intent model.SmartPlacementIntent) (model.MecHost, error) {
	var err error
	var sp SearchParams
	var bestMec model.MecHost

	tc := topology.NewTopologyClient()
	tc.CurrentCell = intent.Spec.SmartPlacementIntentData.TargetCell

	log.Infof("\n\nReceived request about finding new cluster for app: %v located at Cluster: %v, that moves towards cell: %v\n", intent.Spec.AppName, intent.CurrentPlacement.Cluster, intent.Spec.SmartPlacementIntentData.TargetCell)
	log.Infof("Searching Type: Optimal")
	log.Infof("Smart Placement Intent: %+v", intent)

	// Evaluate what coverage zone we are considering (based on region)
	firstMEC, err := tc.GetMecHostsByCellId(tc.CurrentCell)

	sp.currentMECs, err = tc.GetMecHostsByRegion(firstMEC[0].Identity.Location.Region)

	// Check if any cluster is a valid candidate, if yes -> try to find optimal cluster
	sp, err = FindCandidates(tc, sp, intent)
	if err != nil {
		log.Errorf("Could not serve Smart Placement Intent: %v", err)
		return model.MecHost{}, err
	}

	// Try to find the optimal cluster.
	bestMec, err = FindOptimalCluster(sp.candidateMECs, intent)
	if err != nil {
		log.Warnf(" Could not find optimal cluster for given APP[%v]", intent.Spec.AppName)
		return model.MecHost{}, err
	} else {
		log.Infof("Found optimal cluster[%v] for given APP[%v]", bestMec.Identity.Cluster, intent.Spec.AppName)
		return bestMec, nil
	}

}

// checkIfCurrentClusterOk checks if current MEC Host meets all requirements specified in the SmartPlacementIntent.
// Returns `true` if current MEC Host meets all the requirements, `false` otherwise.
func checkIfCurrentClusterOk(tc topology.Client, i model.SmartPlacementIntent) (bool, error) {

	//todo: currently old cluster is taken from intent, that means it should be properly updated at simu part ( hanlder.go), but it was not tested, so maybe more safe would be to take this data from NMT since, NMT already knows cluster per app
	mec, err := tc.GetMecHost(i.CurrentPlacement.Provider, i.CurrentPlacement.Cluster)
	if err != nil {
		return false, err
	}

	mec, err = tc.CollectResourcesInfo(mec)
	if err != nil {
		return false, err
	}

	latency, err := tc.GetShortestPath(i.Spec.SmartPlacementIntentData.TargetCell, mec)
	if err != nil {
		return false, err
	}
	mec.SetLatency(latency)

	if resourcesOk(i, mec) && latencyOk(i, mec) {
		log.Infof("Current MEC[%v+%v] is OK. Skipping", i.CurrentPlacement.Provider, i.CurrentPlacement.Cluster)
		return true, nil
	}

	return false, err
}

// checkIfSkip checks if mec in already present in mecList (based on provider+cluster identifier)
// Returns `true` if it is, `false` otherwise.
func checkIfSkip(mec model.MecHost, mecList []model.MecHost) bool {
	for _, candidate := range mecList {
		if mec.BuildClusterEmcoFQDN() == candidate.BuildClusterEmcoFQDN() {
			return true
		}
	}
	return false
}

// checkCoverageZone checks if source MEC Host (sMec) and target MEC Host (tMec) are located in the same Coverage Zone.
// Returns `true` if it is, `false` otherwise.
func checkCoverageZone(sMec, tMec model.MecHost) bool {
	if sMec.Identity.Location.Region != tMec.Identity.Location.Region {
		return true
	}
	return false
}

// FindCandidates based on SmartPlacementIntent and current SearchParams, iterates over MEC Hosts and checks the constraints
// If latency constraint for given MEC Host is met, it's appended to the SearchParams.evalNeighMECs for further search
// If all constraints for given MEC Host are met, it's appended to the SearchParams.candidateMECs list
// Returns (updated SearchParams struct, error)
func FindCandidates(tc *topology.Client, sp SearchParams, i model.SmartPlacementIntent) (SearchParams, error) {
	log.Infof("Among all clusters in Searching Area let's identify those that fulfill latency, resources and max load constraints")
	for _, mec := range sp.currentMECs {
		// Shortest path is a sum of latency on all the links -> from targetCell to the candidate MEC Host
		latency, err := tc.GetShortestPath(tc.CurrentCell, mec)
		if err != nil {
			log.Warnf("Could not get shortest path: %v. Skipping.", err)
			continue
		}

		mec.SetLatency(latency)

		// If latency is not met for this MEC Host -> his neighbours can't be a valid candidates
		if !latencyOk(i, mec) {
			continue
		}

		// If considered mec meets latency requirements, we can consider his neighbours as candidates later
		sp.evalNeighMECs = append(sp.evalNeighMECs, mec)

		mec, err = tc.CollectResourcesInfo(mec)
		if err != nil {
			log.Warnf("Could not collect MEC resources. Error: %v", err)
		}

		if !resourcesOk(i, mec) {
			log.Warnf("Resources condition for cluster [%v] not met. Skipping. CPU UTIL: %v, MEM util: %v", mec.BuildClusterEmcoFQDN(), mec.Resources.Cpu.Utilization, mec.Resources.Memory.Utilization)
			continue
		}

		log.Infof("Identified possible MEC: %+v: %v", mec.Identity.Cluster, mec.Resources)

		// TODO: we can specify more restrictions for the candidate MEC Hosts: for example consider only specific level/region
		sp.candidateMECs = append(sp.candidateMECs, mec)
	}
	//log.Infof("Candidates list: %v", sp.candidateMECs)
	return sp, nil
}

// latencyOk checks if latency constraints specified in SmartPlacementIntent (i) are met
// Returns true if yes, false otherwise
func latencyOk(i model.SmartPlacementIntent, mec model.MecHost) bool {
	latency := mec.GetLatency()
	latencyMax := i.Spec.SmartPlacementIntentData.ConstraintsList.LatencyMax

	if mec.GetLatency() < 0 {
		return false
	} else if latencyMax > latency {
		return true
	} else {
		return false
	}
}

// resourcesOk checks if resource constraints specified in SmartPlacementIntent (i) are met
// Returns true if yes, false otherwise
func resourcesOk(i model.SmartPlacementIntent, mec model.MecHost) bool {
	var cpuUtilization, memUtilization float64
	if i.CurrentPlacement.Provider == mec.Identity.Provider && i.CurrentPlacement.Cluster == mec.Identity.Cluster {
		// If we consider the current cluster, we don't need to assume app requirements
		// because this application is already placed on this mec.
		cpuUtilization = 100 * (mec.GetCpuUsed()) / mec.GetCpuCapacity()
		memUtilization = 100 * (mec.GetMemUsed()) / mec.GetMemCapacity()
	} else {
		cpuUtilization = 100 * (mec.GetCpuUsed() + i.Spec.SmartPlacementIntentData.AppCpuReq) / mec.GetCpuCapacity()
		memUtilization = 100 * (mec.GetMemUsed() + i.Spec.SmartPlacementIntentData.AppMemReq) / mec.GetMemCapacity()
	}

	cpuMax, _ := strconv.ParseFloat(config.GetConfiguration().Tau, 64) //80
	memMax, _ := strconv.ParseFloat(config.GetConfiguration().Tau, 64) // 80
	cpuMecAvaliable := mec.GetCpuCapacity() - mec.GetCpuUsed()
	memMecAvaliable := mec.GetMemCapacity() - mec.GetMemUsed()

	if cpuUtilization < 0 || memUtilization < 0 {
		//log.Warnf("[RES-CHECK][DEBUG] cpuUtilization[%v], memUtilization[%v]", cpuUtilization, memUtilization)
		return false
	} else if cpuMecAvaliable < i.Spec.SmartPlacementIntentData.AppCpuReq {
		//log.Warnf("[RES-CHECK][DEBUG] cpuMecAvaliable[%v] < appCpuReq[%v] = true", cpuMecAvaliable, i.Spec.SmartPlacementIntentData.AppCpuReq)
		return false
	} else if memMecAvaliable < i.Spec.SmartPlacementIntentData.AppMemReq {
		//log.Warnf("[RES-CHECK][DEBUG] memMecAvaliable[%v] < appMemReq[%v] = true", memMecAvaliable, i.Spec.SmartPlacementIntentData.AppMemReq)
		return false
	} else if cpuMax >= cpuUtilization && memMax >= memUtilization {
		//log.Warnf("[RES-CHECK][DEBUG] Resources OK!")
		return true
	} else {
		//log.Warnf("[RES-CHECK][DEBUG] Resources not OK :/")
		return false
	}
}

// FindOptimalCluster iterates over candidate MEC Hosts (mecs), and based on the computed cost, it selects the best one
// Returns the best (cheapest) MEC Host if found, error otherwise (mecs list is empty)
func FindOptimalCluster(mecs []model.MecHost, intent model.SmartPlacementIntent) (model.MecHost, error) {
	log.Infof("Looking for optimal cluster among clusters that fulfilled constaraints... Cost is:")

	if len(mecs) <= 0 {
		reason := "candidate MECs list is empty"
		err := errors.New(reason)
		log.Warnf("[RESULT] Could not find optimal cluster: %v", reason)
		return model.MecHost{}, err
	}

	var bestMec model.MecHost
	var bestCost float64 = math.Inf(1)
	for _, mec := range mecs {
		currentCost := ComputeObjectiveValue(intent, mec)
		log.Infof("MEC: %v, cost: %v", mec.Identity.Cluster, currentCost)
		if currentCost < bestCost {
			bestMec = mec
			bestCost = currentCost
		}
	}

	log.Infof("Found best cluster: %+v", bestMec)
	return bestMec, nil
}

// ComputeObjectiveValue based on SmartPlacementIntent evaluate MEC Host overall cost
// Returns MEC Host cost value as float64
func ComputeObjectiveValue(i model.SmartPlacementIntent, mec model.MecHost) float64 {
	var staticCost float64
	switch mec.Identity.Location.Type {
	case 0:
		staticCost = 1.0
	case 1:
		staticCost = 0.66667
	case 2:
		staticCost = 0.33333
	default:
		log.Warnf("[INTENT] MEC Type[%v] should not be considered: static-cost[%v].", mec.Identity.Location.Type, 404)
		staticCost = 404
	}
	pw := i.Spec.SmartPlacementIntentData.ParametersWeights

	nLat, nCpu, nMem := NormalizeMecParameters(mec)

	log.Infof("Calculating Objective Function: Lateny weight: %v, [%v] | Res wight: %v [CPU: weights %v [Util: %v], MEM: : weights %v [Util: %v] | Static cost: %v]",
		pw.LatencyWeight, nLat, pw.ResourcesWeight, pw.CpuUtilizationWeight, nCpu, pw.MemUtilizationWeight, nMem, staticCost)

	return pw.LatencyWeight*nLat + pw.ResourcesWeight*(pw.CpuUtilizationWeight*nCpu+pw.MemUtilizationWeight*nMem)*staticCost
}

// NormalizeMecParameters returns latency/resources parameters as values from range(0, 1)
// TODO: For now consider 30ms as the maximum latency to normalize the value.
func NormalizeMecParameters(mec model.MecHost) (float64, float64, float64) {
	log.Infof("Mem utilization before normalization: %v", mec.GetMemUtilization())
	return mec.GetLatency() / 30, mec.GetCpuUtilization(), mec.GetMemUtilization()
}
