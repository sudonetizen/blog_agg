package config

import (
    "os"
    "fmt"
    "encoding/json"
)

const configFile = ".gatorconfig.json"

type Config struct {
    Db_url            string `json:"db_url"`
    Current_user_name string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
    val, err := os.UserHomeDir()
    if err != nil {return "", err}
    configFilePath := val+"/"+configFile
    return configFilePath, nil
}

func Read() Config {
    path, err := getConfigFilePath()
    if err != nil {fmt.Errorf("error %v", err)} 
    
    data, err := os.ReadFile(path)
    if err != nil {fmt.Errorf("error %v", err)} 

    configStruct := Config{}
    err = json.Unmarshal(data, &configStruct)  
    if err != nil {fmt.Errorf("error %v", err)} 

    return configStruct 
}

func SetUser(cfg Config, user_name string) {
    path, err := getConfigFilePath()
    if err != nil {fmt.Errorf("error %v", err)} 

    cfg.Current_user_name = user_name
    data, err := json.Marshal(cfg)
    if err != nil {fmt.Errorf("error %v", err)} 

    err = os.WriteFile(path, data, 0666)
    if err != nil {fmt.Errorf("error %v", err)} 
}


