package latency

import (
	"10.254.188.33/matyspi5/erd/pkg/obs/src/config"
	log "10.254.188.33/matyspi5/erd/pkg/obs/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/obs/src/pkg/model"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type MockClient struct {
	Cells []model.Cell
	Mecs  []model.MecHost
}

func (mc *MockClient) InitializeLatencyMock() {
	log.Info("[LTC] Initializing LatencyMock...")
	rand.Seed(time.Now().UnixNano())

	var cells []model.Cell
	cells = config.ReadNetworkTopologyConfigFile("networkTopology.json")

	mc.Cells = cells

	var mecHosts []model.MecHost
	mecHosts = config.ReadTopologyConfigFile("mecTopology.json")

	mc.Mecs = mecHosts

}

func (mc *MockClient) GetCellById(cellId string) (model.Cell, error) {
	for _, cell := range mc.Cells {
		if cell.Id == cellId {
			return cell, nil
		}
	}
	err := errors.New(fmt.Sprintf("couldn't find cell[%v]", cellId))
	return model.Cell{}, err
}

func (mc *MockClient) GetMehByFqdn(mecFqdn string) (model.MecHost, error) {
	for _, mec := range mc.Mecs {
		if mec.BuildClusterEmcoFQDN() == mecFqdn {
			return mec, nil
		}
	}
	err := errors.New(fmt.Sprintf("couldn't find mecHost[%v]", mecFqdn))
	return model.MecHost{}, err
}

func (mc *MockClient) CellExists(cellId string) bool {
	_, err := mc.GetCellById(cellId)
	if err != nil {
		return false
	}
	return true
}

func (mc *MockClient) MehExists(mecFqdn string) bool {
	_, err := mc.GetMehByFqdn(mecFqdn)
	if err != nil {
		return false
	}
	return true
}

func (mc *MockClient) GetMockedLatencyMs(sNode, tNode interface{}) (float64, error) {

	//latency := generateInitialLatency()
	var latency float64
	switch source := sNode.(type) {
	case model.Cell:
		switch target := tNode.(type) {
		case model.Cell:
			// source -> cell; target -> cell
			// not allowed
		case model.MecHost:
			// source -> cell; target -> mecHost
			if !mc.CellExists(source.Id) || !mc.MehExists(target.BuildClusterEmcoFQDN()) {
				err := errors.New(fmt.Sprintf("couldn't get latency. reason: cell-id: %v or mec-name: %v, doesn't exist", source.Id, target.BuildClusterEmcoFQDN()))

				log.Errorf("[LTC] error: %v", err)
				return -1, err
			}
			if target.Identity.Location.Level == 0 {
				latency = 0.5
			} else {
				latency = -404
			}

			// leave latency as it is cell<-> mec latency should be low
		}
	case model.MecHost:
		switch target := tNode.(type) {
		case model.Cell:
			// source -> mecHost; target -> cell
			if !mc.CellExists(target.Id) || !mc.MehExists(source.BuildClusterEmcoFQDN()) {
				err := errors.New(fmt.Sprintf("couldn't get latency. reason: cell-id: %v or mec-name: %v, doesn't exist", target.Id, source.BuildClusterEmcoFQDN()))

				log.Errorf("[LTC] error: %v", err)
				return -1, err
			}
			// leave latency as it is cell<-> mec latency should be low
			if source.Identity.Location.Level == 0 {
				latency = 0.5
			} else {
				latency = -404
			}

		case model.MecHost:
			if source.BuildClusterEmcoFQDN() == target.BuildClusterEmcoFQDN() {
				return 0, nil
			}
			// source -> mecHost; target -> mecHost
			if !mc.MehExists(source.BuildClusterEmcoFQDN()) || !mc.MehExists(target.BuildClusterEmcoFQDN()) {
				err := errors.New(fmt.Sprintf("couldn't get latency. reason: mec-name: %v or mec-name: %v, doesn't exist",
					source.BuildClusterEmcoFQDN(), target.BuildClusterEmcoFQDN()))

				log.Errorf("[LTC] error: %v", err)
				return -1, err
			}

			levelDiff := float64(source.Identity.Location.Level - target.Identity.Location.Level)

			switch levelDiff {
			case 0:
				if source.Identity.Location.Level == 0 {
					if source.Identity.Location.LocalZone == target.Identity.Location.LocalZone {
						latency = 1
					} else {
						latency = 2
					}
				} else if source.Identity.Location.Level == 1 {
					if source.Identity.Location.Zone == target.Identity.Location.Zone {
						latency = 2
					} else {
						latency = 3
					}
				} else if source.Identity.Location.Level == 2 {
					if source.Identity.Location.LocalZone == target.Identity.Location.LocalZone {
						latency = 3
					} else {
						latency = 4
					}
				}
			case 1:
				if source.Identity.Location.Level == 2 {
					latency = 6
				} else if source.Identity.Location.Level == 1 {
					latency = 4
				}
			case -1:
				if source.Identity.Location.Level == 1 {
					latency = 6
				} else if source.Identity.Location.Level == 0 {
					latency = 4
				}
			case 2, -2:
				// unsupported, no direct links
				latency = -404
			}
		}
	}
	return latency, nil
}

func getRandFloat(min, max float64) float64 {
	res := min + rand.Float64()*(max-min)
	return res
}

func generateInitialLatency() float64 {
	var latency float64

	probability := getRandFloat(0, 1)

	if probability < 0.93 {
		latency = math.Round(getRandFloat(1, 10)*100) / 100
	} else if probability >= 0.93 && probability < 0.98 {
		latency = math.Round(getRandFloat(20, 50)*100) / 100
	} else {
		latency = math.Round(getRandFloat(50, 150)*100) / 100
	}
	return latency
}
