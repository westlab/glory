package glory

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Conf struct {
	UpdateIntervalMin int             `json:"update_interval_min"`
	WorkingGroups     []*WorkingGroup `json:"working_groups"`
}

type WorkingGroup struct {
	Title    string   `json:"title"`
	Describe string   `json:"describe"`
	Deadline string   `json:"Deadline"`
	DirID    string   `json:"dir_id"`
	Members  []string `json:"members"`
}

func LoadConfig(filePath string) (*Conf, error) {
	var config Conf
	configSrc, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read config error: %w", err)
	}

	if err = json.Unmarshal(configSrc, &config); err != nil {
		return nil, fmt.Errorf("parse config error: %w", err)
	}

	if config.UpdateIntervalMin <= 0 {
		config.UpdateIntervalMin = 1 //リソースの浪費対策
	}

	return &config, nil
}
