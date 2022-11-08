package model

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func (c *InputYaml) GetYamlFile(filePath string) *InputYaml {

	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func (c *InputYaml) SaveYamlFile(fileName string) {

	yamlData, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf("Marshal: %v", err)
	}

	//fmt.Println(string(yamlData))

	err = ioutil.WriteFile(fileName, yamlData, 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}
}

type InputYaml struct {
	HostIP             string `yaml:"HostIP"`
	RsyncPort          int    `yaml:"RsyncPort"`
	TemporalServerIP   string `yaml:"TemporalServerIP"`
	TemporalServerPort int    `yaml:"TemporalServerPort"`
	ProjectName        string `yaml:"ProjectName"`
	AdminCloud         string `yaml:"AdminCloud"`
	Controllers        []struct {
		Name     string `yaml:"Name"`
		Port     int    `yaml:"Port"`
		Anchor   string `yaml:"Anchor"`
		Priority int    `yaml:"Priority"`
		Type     string `yaml:"Type"`
	} `yaml:"Controllers"`
	Deployments []struct {
		CaName      string `yaml:"CaName"`
		ProfileName string `yaml:"ProfileName"`
		DigName     string `yaml:"DigName"`
		GpiName     string `yaml:"GpiName"`
		TacIntent   string `yaml:"tacIntent"`
		Apps        []struct {
			Name              string `yaml:"Name"`
			AppFile           string `yaml:"AppFile"`
			ProfileFile       string `yaml:"ProfileFile"`
			PlacementClusters []struct {
				Provider string   `yaml:"Provider"`
				Clusters []string `yaml:"Clusters"`
			} `yaml:"PlacementClusters"`
			Workflows []struct {
				WfIntentName         string        `yaml:"WfIntentName"`
				WfDescription        string        `yaml:"WfDescription"`
				WfType               string        `yaml:"WfType"`
				HookType             string        `yaml:"HookType"`
				WfClientEndpointName string        `yaml:"WfClientEndpointName"`
				WfClientEndpointPort int           `yaml:"WfClientEndpointPort"`
				WfClientName         string        `yaml:"WfClientName"`
				WfID                 string        `yaml:"WfID"`
				WfTaskQueue          string        `yaml:"WfTaskQueue"`
				WfEmcoClm            string        `yaml:"WfEmcoClm"`
				WfEmcoOrch           string        `yaml:"WfEmcoOrch"`
				WfEmcoOrchStatus     string        `yaml:"WfEmcoOrchStatus"`
				WfEmcoMgr            string        `yaml:"WfEmcoMgr"`
				Params               EdgeAppParams `yaml:"Params"`
			} `yaml:"Workflows"`
		} `yaml:"Apps"`
	} `yaml:"Deployments"`
	Providers []struct {
		Name     string `yaml:"Name"`
		Clusters []struct {
			Name       string `yaml:"Name"`
			Label      string `yaml:"Label"`
			Reference  string `yaml:"Reference"`
			KubeConfig string `yaml:"KubeConfig"`
		} `yaml:"Clusters"`
	} `yaml:"Providers"`
}

type EdgeAppParams struct {
	LatencyMax         float64 `yaml:"LatencyMax"`
	CPUUtilMax         float64 `yaml:"CpuUtilMax"`
	MemUtilMax         float64 `yaml:"MemUtilMax"`
	LtcWeight          float64 `yaml:"LtcWeight"`
	ResWeight          float64 `yaml:"ResWeight"`
	CPUWeight          float64 `yaml:"CpuWeight"`
	MemWeight          float64 `yaml:"MemWeight"`
	AppCPUReq          float64 `yaml:"AppCpuReq"`
	AppMemReq          float64 `yaml:"AppMemReq"`
	InnotURL           string  `yaml:"InnotUrl"`
	PlcControllerURL   string  `yaml:"PlcControllerUrl"`
	RelocateClientName string  `yaml:"RelocateClientName"`
	RelocateClientPort int     `yaml:"RelocateClientPort"`
	RelocateWfName     string  `yaml:"RelocateWfName"`
}
