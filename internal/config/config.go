package config

import (
	"encoding/json"
	"os"

	"github.com/adrg/xdg"
	"github.com/spf13/viper"
)

var configPath = xdg.ConfigHome + "/ticktask"
var dataPath = xdg.DataHome + "/ticktask/"
var configName = "ticktask"
var configType = "json"

// LoadConfig loads the configuration
func LoadConfig() error {
	viper.BindEnv("task_path", "TT_TASK_PATH")
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)
	viper.SetDefault("task_path", dataPath)
	viper.SetDefault("ssh_key", os.Getenv("HOME")+"/.ssh/id_rsa")
	viper.SetDefault("profiles", []map[string]interface{}{
		{
			"name":     "default",
			"contexts": []string{},
			"projects": []string{},
			"active":   false,
		},
	})

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil

}

func GetConfigPath(defaultConfig bool) string {
	if defaultConfig {
		return configPath + "/" + configName + "." + configType
	}
	return viper.ConfigFileUsed()
}

func GetConfigJSON() (*string, error) {
	configJSON, err := json.MarshalIndent(viper.AllSettings(), "", "  ")
	if err != nil {
		return nil, err
	}
	configString := string(configJSON)
	return &configString, nil

}

func GetTaskPath() string {
	return viper.GetString("task_path")
}

func GetTodoTxtPath() string {
	return viper.GetString("task_path") + "/todo.txt"
}

// SaveConfig saves the configuration
func SaveConfig() error {
	//viper.Set("server", c.Server)
	//viper.Set("auth", c.Auth)
	//viper.Set("db_path", c.DbPath)

	err := viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func GetSSHPublicKeyPath() string {
	return viper.GetString("ssh_key")
}
