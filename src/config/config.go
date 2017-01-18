package config

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	MongodbHostname string `json:"mongodb_hostname"`
	Secret          string `json:"secret"`
	JWTIssuer       string `json:"jwt_issuer"`
}

const configFile = "conf.json"

var config *Configuration // Global configuration instance

// GetConfig returns the global config object or if it
// does not exist it creates it and returns it.
func GetConfig() Configuration {
	if config != nil {
		return *config
	}

	config = &Configuration{}

	file, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)

	err = decoder.Decode(config)
	if err != nil {
		log.Fatal(err)
	}

	return *config
}

// TODO: Implement
func UpdateConfig() {

}
