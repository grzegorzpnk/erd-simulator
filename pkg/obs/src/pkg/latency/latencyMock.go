package latency

import (
	"10.254.188.33/matyspi5/erd/pkg/obs/src/config"
	log "10.254.188.33/matyspi5/erd/pkg/obs/src/logger"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type MockClient struct {
	Cells []CellInfo
	Mehs  []MehInfo
}

type CellInfo struct {
	cellId string
	// leave it for now
	latitude  float64
	longitude float64
}

type MehInfo struct {
	mehName string
	// leave it for now
	latitude  float64
	longitude float64
}

func (mc *MockClient) InitializeLatencyMock() {
	log.Info("[LTC] Initializing LatencyMock...")
	rand.Seed(time.Now().UnixNano())
	for _, cell := range config.GetConfiguration().Cells {
		cl := CellInfo{
			cellId: cell,
		}
		mc.Cells = append(mc.Cells, cl)
	}

	for _, meh := range config.GetConfiguration().Clusters {
		provider, clusters := meh.Provider, meh.Clusters
		for _, cluster := range clusters {
			mh := MehInfo{
				mehName: provider + "+" + cluster,
			}
			mc.Mehs = append(mc.Mehs, mh)
		}
	}
}

func (mc *MockClient) GetCellById(cellId string) CellInfo {
	for _, cell := range mc.Cells {
		if cell.cellId == cellId {
			return cell
		}
	}
	return CellInfo{}
}

func (mc *MockClient) GetMehByName(mehName string) MehInfo {
	for _, meh := range mc.Mehs {
		if meh.mehName == mehName {
			return meh
		}
	}
	return MehInfo{}
}

func (mc *MockClient) CellExists(cellId string) bool {
	cell := mc.GetCellById(cellId)
	if cell.cellId != "" {
		return true
	}
	return false
}

func (mc *MockClient) MehExists(mehName string) bool {
	meh := mc.GetMehByName(mehName)
	if meh.mehName != "" {
		return true
	}
	return false
}

func (mc *MockClient) GetMockedLatencyMs(cell, meh string) (float64, error) {
	if !mc.CellExists(cell) || !mc.MehExists(meh) {
		err := errors.New(fmt.Sprintf("could not get latency. reason: cell-id: %v or mec-name: %v, doesn't exist", cell, meh))
		log.Errorf("[LTC] error: %v", err)
		return -1, err
	}

	probability := getRandFloat(0, 1)

	if probability < 0.93 {
		return math.Round(getRandFloat(1, 10)*100) / 100, nil
	} else if probability >= 0.93 && probability < 0.98 {
		return math.Round(getRandFloat(100, 300)*100) / 100, nil
	} else {
		return math.Round(getRandFloat(500, 1500)*100) / 100, nil
	}

}

func getRandFloat(min, max float64) float64 {
	res := min + rand.Float64()*(max-min)
	return res
}
