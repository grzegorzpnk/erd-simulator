package module

import (
	log "10.254.188.33/matyspi5/erd/pkg/erc/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/errs"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/model"
	"10.254.188.33/matyspi5/erd/pkg/erc/src/pkg/topology"
	"github.com/pkg/errors"
	"math"
)

// SmartPlacementIntentClient implements the SmartPlacementIntentManager.
// It will also be used to maintain some localized state (Not important, inherited from sample emco controller)
type SmartPlacementIntentClient struct {
	dbInfo DBInfo
}

type SearchParams struct {
	currentMECs   []model.MecHost
	evalNeighMECs []model.MecHost
	checkedMECs   []model.MecHost
	candidateMECs []model.MecHost
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
}

func (i *SmartPlacementIntentClient) ServeSmartPlacementIntentHeuristic(checkIfClause bool, intent model.SmartPlacementIntent) (model.MecHost, error) {
	var err error
	var sp SearchParams
	var bestMecHost model.MecHost
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
	bestMecHost, err = FindOptimalCluster(sp.candidateMECs, intent)
	if err != nil {
		// save information about checked clusters -> no need to check single cluster twice
		sp.checkedMECs = append(sp.checkedMECs, sp.currentMECs...)
		// currentMECs list will be evaluated based on valid neighbours
		sp.currentMECs = []model.MecHost{}
		// candidateMecs list will be evaluated based on new currentMECs list
		sp.candidateMECs = []model.MecHost{}
		// flag which indicates, if we should search further
		//	* search until best cluster is found OR
		//	* until there are any clusters on the evalNeighMECs list
		checkFurther = true
	}

	for checkFurther {
		checkFurther = false
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

		sp, err = FindCandidates(tc, sp, intent)
		if err != nil {
			log.Errorf("Could not serve Smart Placement Intent: %v", err)
			return model.MecHost{}, err
		}

		bestMecHost, err = FindOptimalCluster(sp.candidateMECs, intent)
		if err != nil {
			sp.checkedMECs = append(sp.checkedMECs, sp.currentMECs...)
			sp.currentMECs = []model.MecHost{}
			sp.candidateMECs = []model.MecHost{}
			checkFurther = true
		} else {
			// if error is nil -> best MEC Host is found. Just return.
			return bestMecHost, err
		}
	}

	return bestMecHost, err
}

func (i *SmartPlacementIntentClient) ServeSmartPlacementIntentOptimal(intent model.SmartPlacementIntent) (model.MecHost, error) {
	var err error
	var sp SearchParams
	var bestMecHost model.MecHost

	tc := topology.NewTopologyClient()
	tc.CurrentCell = intent.Spec.SmartPlacementIntentData.TargetCell

	log.Infof("Smart Placement Intent: %+v", intent)

	// Evaluate what coverage zone we are considering (based on region)
	temp, err := tc.GetMecHostsByCellId(tc.CurrentCell)

	sp.currentMECs, err = tc.GetMecHostsByRegion(temp[0].Identity.Location.Region)

	// Check if any cluster is a valid candidate, if yes -> try to find optimal cluster
	sp, err = FindCandidates(tc, sp, intent)
	if err != nil {
		log.Errorf("Could not serve Smart Placement Intent: %v", err)
		return model.MecHost{}, err
	}

	// Try to find the optimal cluster.
	bestMecHost, err = FindOptimalCluster(sp.candidateMECs, intent)
	if err != nil {
		log.Warnf(" Could not find optimal cluster for given APP[%v]", intent.Spec.AppName)
		return model.MecHost{}, err
	} else {
		log.Infof("Found optimal cluster[%v] for given APP[%v]", bestMecHost.Identity.Cluster, intent.Spec.AppName)
		return bestMecHost, nil
	}

}

func checkIfCurrentClusterOk(tc topology.Client, i model.SmartPlacementIntent) (bool, error) {
	mecHost, err := tc.GetMecHost(i.CurrentPlacement.Provider, i.CurrentPlacement.Cluster)
	if err != nil {
		return false, err
	}

	mecHost, err = tc.CollectResourcesInfo(mecHost)
	if err != nil {
		return false, err
	}

	latency, err := tc.GetShortestPath(i.Spec.SmartPlacementIntentData.TargetCell, mecHost)
	mecHost.Resources.Latency = latency

	if resourcesOk(i, mecHost) && latencyOk(i, mecHost) {
		log.Infof("Current mecHost[%v+%v] is OK. Skipping", i.CurrentPlacement.Provider, i.CurrentPlacement.Cluster)
		return true, nil
	}

	return false, err
}

func checkIfSkip(mec model.MecHost, mecList []model.MecHost) bool {
	for _, candidate := range mecList {
		if mec.BuildClusterEmcoFQDN() == candidate.BuildClusterEmcoFQDN() {
			//log.Infof("Skipping: %v", candidate.BuildClusterEmcoFQDN())
			return true
		}
	}
	return false
}

func checkCoverageZone(sMec, tMec model.MecHost) bool {
	if sMec.Identity.Location.Region != tMec.Identity.Location.Region {
		return true
	}
	return false
}

func FindCandidates(tc *topology.Client, sp SearchParams, i model.SmartPlacementIntent) (SearchParams, error) {
	log.Infof("Looking for candidates...")
	for _, mec := range sp.currentMECs {
		// Shortest path is a sum of latency on all the links -> from targetCell to the candidate MEC Host
		latency, err := tc.GetShortestPath(tc.CurrentCell, mec)
		if err != nil {
			log.Warnf("Could not get shortest path: %v. Skipping.", err)
			continue
		}

		mec.SetLatency(latency)

		if !latencyOk(i, mec) {
			//log.Warnf("Latency condition for cluster [%v] not met. Skipping.", mec.BuildClusterEmcoFQDN())
			continue
		}

		// If considered mec meets latency requirements, we can consider his neighbours as candidates later
		// If latency is not met for this MEC Host -> his neighbours can't be a valid candidates
		sp.evalNeighMECs = append(sp.evalNeighMECs, mec)

		mec, err = tc.CollectResourcesInfo(mec)
		if err != nil {
			log.Warnf("Could not collect MEC resources. Error: %v", err)
		}

		log.Infof("Current MEC is: %+v", mec)

		if !resourcesOk(i, mec) {
			//log.Warnf("Resources condition for cluster [%v] not met. Skipping.", mec.BuildClusterEmcoFQDN())
			continue
		}

		// TODO: we can specify more restrictions for the candidate MEC Hosts: for example consider only specific level/region
		sp.candidateMECs = append(sp.candidateMECs, mec)
	}
	log.Infof("Candidates list: %v", sp.candidateMECs)
	return sp, nil
}

// latencyOk checks if latency constraints specified in intent (i) are met
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

// resourcesOk checks if resource constraints specified in intent (i) are met
func resourcesOk(i model.SmartPlacementIntent, mec model.MecHost) bool {
	var cpuUtilization, memUtilization float64
	if i.CurrentPlacement.Provider == mec.Identity.Provider && i.CurrentPlacement.Cluster == mec.Identity.Cluster {
		// If we consider the current cluster, we don't need to assume app requirements
		// because this application is already placed on this mec.
		cpuUtilization = 100 * (mec.GetCpuUsed()) / mec.GetCpuAllocatable()
		memUtilization = 100 * (mec.GetMemUsed()) / mec.GetMemAllocatable()
	} else {
		cpuUtilization = 100 * (mec.GetCpuUsed() + i.Spec.SmartPlacementIntentData.AppCpuReq) / mec.GetCpuAllocatable()
		memUtilization = 100 * (mec.GetMemUsed() + i.Spec.SmartPlacementIntentData.AppMemReq) / mec.GetMemAllocatable()
	}

	cpuMax := i.Spec.SmartPlacementIntentData.ConstraintsList.CpuUtilizationMax // 80
	memMax := i.Spec.SmartPlacementIntentData.ConstraintsList.MemUtilizationMax // 80
	cpuMecAvaliable := mec.GetCpuAllocatable() - mec.GetCpuUsed()
	memMecAvaliable := mec.GetMemAllocatable() - mec.GetMemUsed()

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

func FindOptimalCluster(mecHosts []model.MecHost, intent model.SmartPlacementIntent) (model.MecHost, error) {
	log.Infof("Looking for optimal cluster...")

	if len(mecHosts) <= 0 {
		reason := "mec host list is empty"
		err := errors.New(reason)
		log.Warnf("[RESULT] Could not find optimal cluster: %v", reason)
		return model.MecHost{}, err
	}

	var bestMecHost model.MecHost
	var bestCost float64 = math.Inf(1)
	for _, mecHost := range mecHosts {
		currentCost := ComputeObjectiveValue(intent, mecHost)
		if currentCost < bestCost {
			bestMecHost = mecHost
			bestCost = currentCost
		}
	}

	log.Infof("Found best cluster: %+v", bestMecHost)
	return bestMecHost, nil
}

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
	cl := i.Spec.SmartPlacementIntentData.ConstraintsList
	pw := i.Spec.SmartPlacementIntentData.ParametersWeights

	nLat, nCpu, nMem := NormalizeMecParameters(cl, mec)

	return pw.LatencyWeight*nLat + pw.ResourcesWeight*(pw.CpuUtilizationWeight*nCpu+pw.MemUtilizationWeight*nMem)*float64(staticCost)
}

// NormalizeMecParameters TODO: find the best way to normalize all the values
// TODO: For now consider 30ms as the maximum latency to normalize the value.
// TODO: For now consider 100% as maximum Utilization
func NormalizeMecParameters(cl model.Constraints, mec model.MecHost) (float64, float64, float64) {
	//			[0.33 0.66 1]	0.8									0.8
	return mec.GetLatency() / 30, mec.GetCpuUtilization() / 100, mec.GetMemUtilization() / 100
}
