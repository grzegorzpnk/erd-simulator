package mec_topology

import (
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/model"

	"encoding/json"
	"io"
	"os"
)

// readConfigFile reads the specified smsConfig file to setup some env variables
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
