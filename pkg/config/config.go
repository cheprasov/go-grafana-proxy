package config

import (
    "encoding/json"
    "io/ioutil"
)

type Config struct {
    Listen   string `json:"listen"`
    Compress bool   `json:"compress"`

    PostUrl string `json:"postUrl"`
    ApiKey  string `json:"apiKey"`

    Authorization string `json:"authorization"`
}

func ReadConfig(filename string) (Config, error) {
    var config Config
    jsonData, err := ioutil.ReadFile(filename)
    if err != nil {
        return config, err
    }

    err = json.Unmarshal(jsonData, &config)
    if err != nil {
        return config, err
    }

    return config, nil
}
