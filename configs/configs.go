package configs

import (
	"github.com/danmaina/logger"
	"github.com/go-yaml/yaml"
	"io"
	"messaging/constants"
	"os"
)

// ReadConfigs Reads Configs from file or create default Configs
func ReadConfigs() (*Config, error) {
	logger.DEBUG("Reading Default Config File or Creating Config File if not exists")

	// Fetch/ Create Yaml config file
	configFile, errFetchFile := os.OpenFile(constants.Path, os.O_RDWR|os.O_CREATE, os.ModePerm)

	var ConfigStruct *Config

	if errFetchFile != nil {
		logger.ERR("An Error Occurred while initializing configs: ", errFetchFile)
		return ConfigStruct, errFetchFile
	}

	// Read Contents of config file
	configFileByteArr, errReadingByteArr := io.ReadAll(configFile)

	if errReadingByteArr != nil {
		logger.ERR("An Error Occurred while reading contents of config file: ", errReadingByteArr)
		return ConfigStruct, nil
	}

	// Get config from yaml
	// Get Configuration from file yaml
	errDecodingFileYaml := yaml.Unmarshal(configFileByteArr, &ConfigStruct)

	if errDecodingFileYaml != nil {
		logger.ERR("An Error Occurred while converting yaml to Config Struct: ", errDecodingFileYaml)
		return ConfigStruct, errDecodingFileYaml
	}

	defer configFile.Close()

	// Check if config file is empty? write default configs to file: Return configs from file
	if ConfigStruct == nil || *ConfigStruct == (Config{}) {
		logger.ERR("Config File Does Not Contain any information, Loading Default Configs")

		errDecodingDefaultYaml := yaml.Unmarshal([]byte(constants.DefaultConfigs), &ConfigStruct)

		lenConfigs, errWritingDefaultConfigs := configFile.WriteString(constants.DefaultConfigs)

		if errWritingDefaultConfigs != nil {
			logger.ERR("Could not write default configs to config file")
		} else {
			logger.INFO("Wrote Default Configs to file. Bytes Written: ", lenConfigs)
		}

		if errDecodingDefaultYaml != nil {
			logger.ERR("Error Decoding Default Configs Yaml: ", errDecodingDefaultYaml)
		}
	}

	return ConfigStruct, nil
}
