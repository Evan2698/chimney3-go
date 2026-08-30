package settings

import "testing"

func TestSettingsValidateAllowsPortOnlyValue(t *testing.T) {
	cfg := &Settings{Listen: "1080", RemoteListen: "2080", Udplisten: "3080", Httpurl: "8080"}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() returned error for valid port-only config: %v", err)
	}
}

func TestSettingsValidateRejectsMalformedAddress(t *testing.T) {
	cfg := &Settings{Listen: "bad-address"}
	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() expected error for malformed address")
	}
}

func TestSettingsValidateRejectsUnknownServiceMetadata(t *testing.T) {
	cfg := &Settings{
		Listen:       "127.0.0.1:1080",
		RemoteListen: "127.0.0.1:1081",
		Udplisten:    "127.0.0.1:1082",
		Httpurl:      "127.0.0.1:8080",
		Method:       "UNKNOWN",
		Which:        "unknown-service",
		Mode:         "invalid-mode",
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() expected error for unsupported method/service/mode")
	}
}

func TestSettingsValidateAcceptsCaseInsensitiveMethod(t *testing.T) {
	cfg := &Settings{
		Listen:       "127.0.0.1:1080",
		RemoteListen: "127.0.0.1:1081",
		Udplisten:    "127.0.0.1:1082",
		Httpurl:      "127.0.0.1:8080",
		Method:       "  chacha-20  ",
		Which:        "socks5",
		Mode:         "server",
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() should accept case/space variations in method names: %v", err)
	}
}

func TestSettingsNormalizeCanonicalizesInputs(t *testing.T) {
	cfg := &Settings{
		Listen:       "  127.0.0.1:1080  ",
		RemoteListen: "  127.0.0.1:1081  ",
		Udplisten:    "  127.0.0.1:1082  ",
		Httpurl:      "  127.0.0.1:8080  ",
		Method:       "  CHACHA-20  ",
		Which:        "  SOCKS5  ",
		Mode:         "  Server  ",
	}

	cfg.Normalize()

	if cfg.Method != "CHACHA-20" {
		t.Fatalf("Normalize() should canonicalize method case without altering the supported value, got %q", cfg.Method)
	}
	if cfg.Which != "socks5" {
		t.Fatalf("Normalize() should lowercase service names, got %q", cfg.Which)
	}
	if cfg.Mode != "server" {
		t.Fatalf("Normalize() should lowercase runtime modes, got %q", cfg.Mode)
	}
	if cfg.Listen != "127.0.0.1:1080" || cfg.RemoteListen != "127.0.0.1:1081" || cfg.Udplisten != "127.0.0.1:1082" || cfg.Httpurl != "127.0.0.1:8080" {
		t.Fatalf("Normalize() should trim whitespace from address fields, got listen=%q remote=%q udp=%q http=%q", cfg.Listen, cfg.RemoteListen, cfg.Udplisten, cfg.Httpurl)
	}
}

func TestSettingsValidateRejectsNilPointer(t *testing.T) {
	var cfg *Settings
	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() on nil Settings should return an error")
	}
}
