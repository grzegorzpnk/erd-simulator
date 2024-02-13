package mec_topology

import (
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/model"

	"encoding/json"
	"io"
	"os"
)

func (g *Graph) ReadTopologyConfigFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Errorf(err.Error())
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()

	for {
		var mec []*model.MecHost
		err := dec.Decode(&mec)
		if err == io.EOF {
			// all done
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		g.MecHosts = mec
	}
}

// readConfigFile reads the specified smsConfig file to setup some env variables
func (g *Graph) ReadMECConnectionFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Errorf(err.Error())
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()

	for {
		var links []*model.Edge
		err := dec.Decode(&links)
		if err == io.EOF {
			// all done
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		for _, v := range links {
			g.AddLink(*v)
		}
	}

}

func (g *Graph) ReadNetworkTopologyConfigFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Errorf(err.Error())
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	var cells []*model.Cell

	for {
		err := dec.Decode(&cells)
		if err == io.EOF {
			// all done
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		log.Infof("cells: %v", cells)
	}

	g.NetworkCells = cells

	for _, v := range cells {
		for _, n := range g.MecHosts {
			if n.Identity.Location.LocalZone == v.LocalZone {
				n.SupportingCells = append(n.SupportingCells, *v)
			}

		}
	}
}

func (g *Graph) AssigneCapacityToClusters() {

	//for _, v := range g.MecHosts {
	//
	//	if v.Identity.Location.Level == 0 {
	//		v.MemoryResources.Capacity = 4000
	//		v.CpuResources.Capacity = 4000
	//	}
	//	if v.Identity.Location.Level == 1 {
	//		v.MemoryResources.Capacity = 8000
	//		v.CpuResources.Capacity = 8000
	//	}
	//	if v.Identity.Location.Level == 2 {
	//		v.MemoryResources.Capacity = 12000
	//		v.CpuResources.Capacity = 12000
	//	}
	//
	//	v.CpuResources.Utilization = 0
	//	v.CpuResources.Used = 0
	//	v.MemoryResources.Utilization = 0
	//	v.MemoryResources.Used = 0
	//
	//}

	for _, v := range g.MecHosts {

		//initial resources at nodes, cause it imitates the overhead of Operating System ( the same as for demonstrator)
		if v.Identity.Location.Level == 0 {

			v.MemoryResources.Capacity = 4000
			v.MemoryResources.Used = 1112
			v.MemoryResources.Utilization = v.MemoryResources.Used / v.MemoryResources.Capacity

			v.CpuResources.Capacity = 4000
			v.CpuResources.Used = 1552
			v.CpuResources.Utilization = v.CpuResources.Used / v.CpuResources.Capacity

		}
		if v.Identity.Location.Level == 1 {

			v.MemoryResources.Capacity = 8000
			v.MemoryResources.Used = 1080
			v.MemoryResources.Utilization = v.MemoryResources.Used / v.MemoryResources.Capacity

			v.CpuResources.Capacity = 8000
			v.CpuResources.Used = 1200
			v.CpuResources.Utilization = v.CpuResources.Used / v.CpuResources.Capacity
		}
		if v.Identity.Location.Level == 2 {

			v.MemoryResources.Capacity = 12000
			v.MemoryResources.Used = 1080
			v.MemoryResources.Utilization = v.MemoryResources.Used / v.MemoryResources.Capacity

			v.CpuResources.Capacity = 12000
			v.CpuResources.Used = 1548
			v.CpuResources.Utilization = v.CpuResources.Used / v.CpuResources.Capacity
		}

		log.Infof("[Cluster: %v] Initial: Mem util: %v procent, CPU util:  %v procent", v.Identity.Cluster, v.MemoryResources.Utilization, v.CpuResources.Utilization)

	}

}
