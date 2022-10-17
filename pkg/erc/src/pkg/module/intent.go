package module

import (
	log "10.254.188.33/matyspi5/erd/pkg/erc/src/logger"
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
	ServeSmartPlacementIntent(intent model.SmartPlacementIntent) (model.MecHost, error)
}

// ServeSmartPlacementIntent TODO: not ready.
// TODO: Consider that mecs need to be in the same region
func (i *SmartPlacementIntentClient) ServeSmartPlacementIntent(intent model.SmartPlacementIntent) (model.MecHost, error) {
	var err error
	var sp SearchParams
	var bestMecHost model.MecHost
	var checkFurther bool

	tc := topology.NewTopologyClient()
	tc.CurrentCell = intent.Spec.SmartPlacementIntentData.TargetCell

	log.Infof("Smart Placement Intent: %+v", intent)

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
			return model.MecHost{}, errors.New("Could not find optimal cluster!")
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

func checkIfSkip(mec model.MecHost, mecList []model.MecHost) bool {
	for _, candidate := range mecList {
		if mec.BuildClusterEmcoFQDN() == candidate.BuildClusterEmcoFQDN() {
			//log.Infof("Skipping: %v", candidate.BuildClusterEmcoFQDN())
			return true
		}
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
// TODO: consider also current application requests
func resourcesOk(i model.SmartPlacementIntent, mec model.MecHost) bool {
	cpu := mec.GetCpuUtilization()
	mem := mec.GetMemUtilization()
	cpuMax := i.Spec.SmartPlacementIntentData.ConstraintsList.CpuUtilizationMax
	memMax := i.Spec.SmartPlacementIntentData.ConstraintsList.MemUtilizationMax

	if cpu < 0 || mem < 0 {
		return false
	} else if cpuMax > cpu && memMax > mem {
		return true
	} else {
		return false
	}
}

// FindOptimalCluster TODO: implement algorithm to find the optimal cluster among mecHosts (candidates)
// This is old dummy implementation.. Here we check constraints for the second time (its redundant -> remove)
func FindOptimalCluster(mecHosts []model.MecHost, intent model.SmartPlacementIntent) (model.MecHost, error) {
	log.Infof("Looking for optimal cluster...")

	if len(mecHosts) <= 0 {
		reason := "mec host list is empty"
		err := errors.New(reason)
		log.Warnf("Could not find optimal cluster: %v", reason)
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

// ComputeObjectiveValue TODO: implement function to calculate the objective function
// TODO: Remember that MEC cost should be considered
func ComputeObjectiveValue(i model.SmartPlacementIntent, mec model.MecHost) float64 {
	var staticCost int
	switch mec.Identity.Location.Type {
	case 0:
		staticCost = 1
	case 1:
		staticCost = 2
	case 2:
		staticCost = 3
	default:
		log.Warnf("[INTENT] MEC Type[%v] should not be considered: static-cost[%v].", mec.Identity.Location.Type, 404)
		staticCost = 404
	}
	cl := i.Spec.SmartPlacementIntentData.ConstraintsList
	pw := i.Spec.SmartPlacementIntentData.ParametersWeights

	nLat, nCpu, nMem := NormalizeMecParameters(cl, mec)

	return float64(staticCost) * (pw.LatencyWeight*nLat + pw.CpuUtilizationWeight*nCpu + pw.MemUtilizationWeight*nMem)
}

// NormalizeMecParameters TODO: find the best way to normalize all the values
func NormalizeMecParameters(cl model.Constraints, mec model.MecHost) (float64, float64, float64) {
	return mec.GetLatency() / cl.LatencyMax, mec.GetCpuUtilization() / cl.CpuUtilizationMax, mec.GetMemUtilization() / cl.MemUtilizationMax
}
