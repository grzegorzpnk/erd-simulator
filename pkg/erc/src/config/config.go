// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package config

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"strings"
)

// Configuration loads up all the values that are used to configure
// backend implementations
type Configuration struct {
	PluginDir        string `json:"plugin-dir"`
	ServicePort      string `json:"service-port"`
	TopologyEndpoint string `json:"mec-topology-endpoint"`
	Tau              string `json:"tau"`
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
	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Error getting cwd. Using .")
		cwd = "."
	}

	return &Configuration{
		PluginDir:        cwd,
		ServicePort:      "8686",
		TopologyEndpoint: "dupa",
		Tau:              "80",
	}
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
				log.Printf("%s must be positive, not %d.\n",
					varName, intValue)
				valid = false
			}
		}
	}
	return valid
}

// GetConfiguration returns the configuration for the app.
// It will try to load it if it is not already loaded.
func GetConfiguration() *Configuration {
	if gConfig == nil {
		conf, err := readConfigFile("config.json")
		if err != nil {
			log.Println("Error loading config file: ", err)
			log.Println("Using defaults...")
		}
		gConfig = conf

		if !isValidConfig(gConfig) {
			log.Fatalln("Bad data in config. Exiting.")
			return nil
		}
	}

	return gConfig
}

// SetConfigValue sets a value in the configuration
// This is mostly used to customize the application and
// should be used carefully.
func SetConfigValue(key string, value string) *Configuration {
	c := GetConfiguration()
	if value == "" || key == "" {
		return c
	}

	v := reflect.ValueOf(c).Elem()
	if v.Kind() == reflect.Struct {
		f := v.FieldByName(key)
		if f.IsValid() {
			if f.CanSet() {
				if f.Kind() == reflect.String {
					f.SetString(value)
				}
			}
		}
	}
	return c
}
