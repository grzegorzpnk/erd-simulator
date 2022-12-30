// Based on: https://gitlab.com/project-emco/core/emco-base/-/blob/main/src/orchestrator/pkg/infra/config/config.go

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// Configuration loads up all the values that are used to configure
// backend implementations
type Configuration struct {
	ServicePort               string `json:"service-port"`
	ClusterControllerEndpoint string `json:"cluster-controller-endpoint"`
	EdgeProvider              string `json:"edge-provider"`
	Tau                       string `json:"tau"`

	//values for initial placement controller
	ResMin        string `json:"res-min"`
	ResMax        string `json:"res-max"`
	MaxCellNumber string `json:"max-cell-number"` //this value specify the range to generate initial cells between 0 - MaxCellNumber
}

// Config is the structure that stores the configuration
var gConfig *Configuration

// readConfigFile reads the specified smsConfig file to setup some env variables
func readConfigFile(file string) (*Configuration, error) {
	f, err := os.Open(file)
	if err != nil {
		return defaultConfiguration(), err
	}
	defer f.Close()

	// Setup some defaults here
	// If the json file has values in it, the defaults will be overwritten
	conf := defaultConfiguration()

	// Read the configuration from json file
	decoder := json.NewDecoder(f)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}

func defaultConfiguration() *Configuration {

	return &Configuration{
		ServicePort:               "8080",
		ClusterControllerEndpoint: "http://10.254.185.50:32138/",
		EdgeProvider:              "orange",
		Tau:                       "0.8",
		ResMin:                    "500",
		ResMax:                    "1000",
		MaxCellNumber:             "42",
	}

}

// GetConfiguration returns the configuration for the app.
// It will try to load it if it is not already loaded.
func GetConfiguration() *Configuration {
	if gConfig == nil {
		conf, err := readConfigFile("config.json")
		if err != nil {
			fmt.Println("Error loading config file: \n", err)
			fmt.Println("Using defaults...\n")
		}
		gConfig = conf

		if !isValidConfig(gConfig) {
			fmt.Println("Bad data in config. Exiting.")
			return nil
		}
	}

	return gConfig
}

func isValidConfig(cfg *Configuration) bool {
	valid := true
	members := reflect.ValueOf(cfg).Elem()

	// If a config param has "Time" in its name, and is type int,
	// ensure its value is positive.
	for i := 0; i < members.NumField(); i++ {
		varName := members.Type().Field(i).Name
		varValue := members.Field(i).Interface()
		if strings.Contains(varName, "Time") {
			intValue, ok := varValue.(int)
			if ok && intValue <= 0 {
				fmt.Println("%s must be positive, not %d.\n",
					varName, intValue)
				valid = false
			}
		}
	}
	return valid
}
