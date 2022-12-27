package api

import (
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"strconv"
)

func (h *apiHandler) InstantiateApplication(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var mecApp model.MECApp
	_ = json.NewDecoder(r.Body).Decode(&mecApp)
	log.Infof("Client tries to deploy new mecApp ID: %v at MEC Host: %v \n", mecApp.Id, params["cluster"])

	var mecHost *model.MecHost
	mecHost = h.graphClient.GetMecHost(params["cluster"], "orange")

	if mecHost == nil {
		err := fmt.Errorf("MecHost %v does not exists", params["cluster"])
		log.Errorf(err.Error())
		w.WriteHeader(http.StatusConflict)
	}
	if !(mecHost.CheckEnoughResources(mecApp)) {
		err := fmt.Errorf("Mec App %v,cannot be instantiated beacuse of not enough resources", mecApp.Id)
		log.Errorf(err.Error())
		w.WriteHeader(http.StatusConflict)
	} else {
		mecHost.InstantiateApp(mecApp)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *apiHandler) DeleteApplication(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var mecApp model.MECApp
	_ = json.NewDecoder(r.Body).Decode(&mecApp)
	log.Infof("Client tries to uninstall mecApp ID: %v from MEC Host: %v \n", mecApp.Id, params["cluster"])

	var mecHost *model.MecHost
	mecHost = h.graphClient.GetMecHost(params["cluster"], "orange")

	if !(mecHost.CheckAppExists(mecApp)) {
		err := fmt.Errorf("Mec App %v,cannot be deleted beacuse it does not exists", mecApp.Id)
		log.Errorf(err.Error())
		w.WriteHeader(http.StatusConflict)
	} else {
		mecHost.UninstallApp(mecApp)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *apiHandler) RelocateApplication(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var oldMecHost, newMecHost *model.MecHost
	oldMecHost = h.graphClient.GetMecHost(params["old-cluster"], "orange")
	newMecHost = h.graphClient.GetMecHost(params["new-cluster"], "orange")

	if oldMecHost.Identity.Cluster == newMecHost.Identity.Cluster {
		log.Infof("Application already on cluster - redundant relocation")
		w.WriteHeader(http.StatusOK)
		return
	}

	var mecApp model.MECApp
	_ = json.NewDecoder(r.Body).Decode(&mecApp)
	log.Infof("Client tries to RELOCATE mecApp ID: %v from MEC Host: %v to new MEC Host: %v \n", mecApp.Id, oldMecHost.Identity.Cluster, newMecHost.Identity.Cluster)

	if !(oldMecHost.CheckAppExists(mecApp)) {
		err := fmt.Errorf("Mec App %v,cannot be relocated beacuse it does not exists at source cluster", mecApp.Id)
		log.Errorf(err.Error())
		w.WriteHeader(http.StatusConflict)
		return
	}

	//instantiate
	if !(newMecHost.CheckEnoughResources(mecApp)) {
		err := fmt.Errorf("Mec App %v,cannot be instantiated beacuse of not enough resources", mecApp.Id)
		log.Errorf(err.Error())
		w.WriteHeader(http.StatusConflict)
		return
	} else {
		newMecHost.InstantiateApp(mecApp)
		oldMecHost.UninstallApp(mecApp)
		w.WriteHeader(http.StatusOK)
		log.Infof("App relocated succesfully ! \n")
	}
}

func (h *apiHandler) GenerateInitialClusters(w http.ResponseWriter, r *http.Request) {

	var mhs []*model.MecHost

	mhs = h.graphClient.MecHosts

	var mhsT = make([]*model.MecHost, len(mhs))
	var cells = map[int]int{}
	search := true
	for search {
		fmt.Println("[DEBUG] Starting search.")
		copy(mhsT, mhs)
		cells = generateRandomCells()

		search = false
		for index, edgeApp := range h.graphClient.Application {
			cmh, err := findCanidateMec(edgeApp, cells[index+1], mhsT)
			if err != nil {
				search = true
				fmt.Printf("Could not find candidate mec for App[%v]. Search failed.\n", edgeApp.Id)
				break
			} else {
				fmt.Printf("Found candidate mec: %v\n", cmh)
			}
			mhsT = updateMecResourcesInfo(mhsT, cmh, edgeApp.Apps[0].Workflows[0].Params)
			in.Deployments[index].Apps[0].PlacementClusters[0].Clusters = []string{cmh.Identity.Cluster}
		}
	}
	printCellsInfo(cells)

	//instantiate
}

func findCanidateMec(app *model.MECApp, cell int, mhs []*model.MecHost) (model.MecHost, error) {
	candidates := []model.MecHost{}

	for index, mh := range mhs {
		latency, _ := ShortestPath(model.CellId(strconv.Itoa(cell)), mh)

		mhs[index].Resources.Latency = latency
	}

	for _, mh := range mhs {
		if resourcesOk(app, mh) && latencyOk(app, mh) {
			candidates = append(candidates, mh)
			// We need to prioritize N+1 and N+2 level clusters
			if mh.Identity.Cluster == "mec1" {
				candidates = append(candidates, mh)
				candidates = append(candidates, mh)
				candidates = append(candidates, mh)
			} else if mh.Identity.Cluster == "mec3" || mh.Identity.Cluster == "mec4" || mh.Identity.Cluster == "mec5" ||
				mh.Identity.Cluster == "mec6" || mh.Identity.Cluster == "mec7" {
				candidates = append(candidates, mh)
				candidates = append(candidates, mh)
			}
		} else {
		}
	}
	if len(candidates) <= 0 {
		return model.MecHost{}, errors.New("no candidates found")
	}
	return candidates[rand.Intn(len(candidates))], nil
}

func generateRandomCells() map[int]int {
	var cells map[int]int = map[int]int{}
	for i := 1; i <= 50; i++ {
		cells[i] = rand.Intn(42) + 1
	}
	return cells
}

/*// latencyOk checks if latency constraints specified in intent (i) are met
func latencyOk(eap model.EdgeAppParams, mec model.MecHost) bool {
	latency := mec.GetLatency()
	latencyMax := eap.LatencyMax

	if mec.GetLatency() < 0 {
		return false
	} else if latencyMax > latency {
		return true
	} else {
		return false
	}
}

// resourcesOk checks if resource constraints specified in intent (i) are met
func resourcesOk(eap model.EdgeAppParams, mec model.MecHost) bool {
	var cpuUtilization, memUtilization float64

	cpuUtilization = 100 * (mec.GetCpuUsed() + eap.AppCPUReq) / mec.GetCpuAllocatable()
	memUtilization = 100 * (mec.GetMemUsed() + eap.AppMemReq) / mec.GetMemAllocatable()

	cpuMax := eap.CPUUtilMax // 80
	memMax := eap.MemUtilMax // 80
	cpuMecAvaliable := mec.GetCpuAllocatable() - mec.GetCpuUsed()
	memMecAvaliable := mec.GetMemAllocatable() - mec.GetMemUsed()

	if cpuUtilization < 0 || memUtilization < 0 {
		return false
	} else if cpuMecAvaliable < eap.AppCPUReq {
		return false
	} else if memMecAvaliable < eap.AppMemReq {
		return false
	} else if cpuMax >= cpuUtilization && memMax >= memUtilization {
		//log.Warnf("[RES-CHECK][DEBUG] Resources OK!")
		return true
	} else {
		//log.Warnf("[RES-CHECK][DEBUG] Resources not OK :/")
		return false
	}
}
*/
func printCellsInfo(val interface{}) {
	jsonCells, err := json.Marshal(val)
	if err != nil {
		log.Errorf("Marshal err: %v", err)
	}

	fmt.Println("----- CELLS -----")
	fmt.Println(string(jsonCells))
}
