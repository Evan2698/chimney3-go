package settings

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
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

// Validate checks common address fields are in a sensible host:port or port-only format.
func (s *Settings) Validate() error {
	// helper to validate addresses
	validate := func(addr, field string) error {
		if addr == "" {
			return nil
		}
		// Try host:port first
		if host, port, err := net.SplitHostPort(addr); err == nil {
			if host == "" {
				// allow empty host (means all interfaces) but port must be numeric
				if _, err := strconv.Atoi(port); err != nil {
					return fmt.Errorf("%s: invalid port %q", field, port)
				}
				return nil
			}
			if _, err := strconv.Atoi(port); err != nil {
				return fmt.Errorf("%s: invalid port %q", field, port)
			}
			return nil
		}
		// If no colon, allow a numeric port
		if _, err := strconv.Atoi(addr); err == nil {
			return nil
		}
		return fmt.Errorf("%s: invalid address %q (expected host:port or port)", field, addr)
	}

	if err := validate(s.Listen, "Listen"); err != nil {
		return err
	}
	if err := validate(s.RemoteListen, "RemoteListen"); err != nil {
		return err
	}
	if err := validate(s.Udplisten, "Udplisten"); err != nil {
		return err
	}
	// Httpurl can be host:port or port
	if err := validate(s.Httpurl, "Httpurl"); err != nil {
		return err
	}
	return nil
}
