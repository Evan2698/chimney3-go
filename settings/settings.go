package settings

import (
	"encoding/json"
	"io"
	"os"
)

// Settings represents the flat configuration from configs/setting.json.
// Fields match the JSON keys exactly so Parse can unmarshal directly.
type Settings struct {
	Listen       string `json:"listen"`
	RemoteListen string `json:"remote_listen"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Method       string `json:"method"`
	Which        string `json:"which"`
	Udplisten    string `json:"udplisten"`
	Httpurl      string `json:"httpurl"`
	Mode         string `json:"mode"`
}

// Parse loads the flat settings from a JSON file at path and returns a Settings.
func Parse(path string) (*Settings, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	cfg := &Settings{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
