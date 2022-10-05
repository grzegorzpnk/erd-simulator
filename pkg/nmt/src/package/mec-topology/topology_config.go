package mec_topology

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// readConfigFile reads the specified smsConfig file to setup some env variables
func (g *Graph) ReadTopologyConfigFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()

	for {
		var mec []*MecHost
		err := dec.Decode(&mec)
		if err == io.EOF {
			// all done
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		g.MecHosts = mec

	}

}

func (g *Graph) ReadNetworkTopologyConfigFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	var cells []Cell

	for {
		err := dec.Decode(&cells)
		if err == io.EOF {
			// all done
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(cells)
	}

	for _, v := range cells {
		for _, n := range g.MecHosts {
			if n.Identity.Location.LocalZone == v.LocalZone {
				n.SupportingCells = append(n.SupportingCells, v)
			}

		}
	}

}
