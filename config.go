package ogpapp

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

// Config ogp.app config
type Config struct {
	BaseURL            string  `toml:"base_url"`
	APIServerPort      string  `toml:"api_server_port"`
	TLS                bool    `toml:"tls"`
	KoruriBoldFontPath string  `toml:"koruri_bold_font_path"`
	DefaultImageWidth  int     `toml:"default_image_width"`
	DefaultImageHeight int     `toml:"default_image_height"`
	DefaultFontSize    float64 `toml:"default_font_size"`
	LocalDev           bool    `toml:"local_dev"`
	ServerCertPath     string  `toml:"server_cert_path"`
	ServerKeyPath      string  `toml:"server_key_path"`
}

// NewConfig create app config
func NewConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to open path %v: %w", path, err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file %v: %w", path, err)
	}
	var cfg Config
	if err := toml.Unmarshal(buf, &cfg); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal toml data: %w", err)
	}
	return &cfg, nil
}
