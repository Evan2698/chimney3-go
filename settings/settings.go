package settings

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"chimney3-go/privacy"
)

func (s *Settings) validateAddress(addr, field string) error {
	if addr == "" {
		return nil
	}

	if host, port, err := net.SplitHostPort(addr); err == nil {
		if host == "" {
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

	if _, err := strconv.Atoi(addr); err == nil {
		return nil
	}

	return fmt.Errorf("%s: invalid address %q (expected host:port or port)", field, addr)
}

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
	cfg.Normalize()
	return cfg, nil
}

func validateServiceName(name string) error {
	n := strings.TrimSpace(name)
	if n == "" {
		return nil
	}
	switch strings.ToLower(n) {
	case "socks5", "proxy", "kcp":
		return nil
	default:
		return fmt.Errorf("which: unsupported service %q (expected socks5, proxy, or kcp)", name)
	}
}

func validateRuntimeMode(mode string) error {
	n := strings.TrimSpace(mode)
	if n == "" {
		return nil
	}
	switch strings.ToLower(n) {
	case "server", "client":
		return nil
	default:
		return fmt.Errorf("mode: unsupported mode %q (expected server or client)", mode)
	}
}

func validateMethodName(name string) error {
	n := strings.TrimSpace(name)
	if n == "" {
		return nil
	}
	if privacy.NewMethodWithName(n) == nil {
		return fmt.Errorf("method: unsupported method %q", name)
	}
	return nil
}

func normalizeField(value string) string {
	return strings.TrimSpace(value)
}

func canonicalMethodName(name string) string {
	trimmed := normalizeField(name)
	if trimmed == "" {
		return ""
	}
	return strings.ToUpper(trimmed)
}

// Normalize trims whitespace and canonicalizes values so the runtime treats
// equivalent inputs consistently regardless of caller formatting.
func (s *Settings) Normalize() {
	if s == nil {
		return
	}
	s.Listen = normalizeField(s.Listen)
	s.RemoteListen = normalizeField(s.RemoteListen)
	s.Username = normalizeField(s.Username)
	s.Password = normalizeField(s.Password)
	s.Method = canonicalMethodName(s.Method)
	s.Which = strings.ToLower(normalizeField(s.Which))
	s.Udplisten = normalizeField(s.Udplisten)
	s.Httpurl = normalizeField(s.Httpurl)
	s.Mode = strings.ToLower(normalizeField(s.Mode))
}

// Validate checks common address fields are in a sensible host:port or port-only format.
func (s *Settings) Validate() error {
	if s == nil {
		return fmt.Errorf("settings: nil")
	}
	if err := s.validateAddress(s.Listen, "Listen"); err != nil {
		return err
	}
	if err := s.validateAddress(s.RemoteListen, "RemoteListen"); err != nil {
		return err
	}
	if err := s.validateAddress(s.Udplisten, "Udplisten"); err != nil {
		return err
	}
	if err := s.validateAddress(s.Httpurl, "Httpurl"); err != nil {
		return err
	}
	if err := validateServiceName(s.Which); err != nil {
		return err
	}
	if err := validateRuntimeMode(s.Mode); err != nil {
		return err
	}
	if err := validateMethodName(s.Method); err != nil {
		return err
	}
	return nil
}
