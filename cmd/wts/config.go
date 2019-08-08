package main

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path/filepath"
)

const (
	// Config_V1 - first config version
	Config_V1 = "V1"
)

// Config file structure
type Config struct {
	Version string `yaml:"version"`
	Wts     struct {
		Token string `yaml:"token"`
		Debug bool   `yaml:"debug"`
		Pull  bool   `yaml:"pull"`
	} `yaml:"wts"`
}

// ParseReader filling config struct
func (c *Config) ParseReader(r io.Reader) error {
	if err := yaml.NewDecoder(r).Decode(c); err != nil {
		return err
	}
	if c.Version != Config_V1 {
		return fmt.Errorf("Unsupported version %s", c.Version)
	}
	return nil
}

// ParseFile filling config struct
func (c *Config) ParseFile(name string) error {
	exname := os.ExpandEnv(name)
	dir := filepath.Dir(exname)
	_ = os.MkdirAll(dir, os.ModePerm)
	if _, err := os.Stat(exname); os.IsNotExist(err) {
		def := new(Config)
		def.Version = Config_V1
		def.Wts.Debug = false
		def.Wts.Pull = false
		f, err := os.Create(exname)
		if err != nil {
			return err
		}
		w := bufio.NewWriter(f)
		if err := yaml.NewEncoder(w).Encode(def); err != nil {
			return err
		}
		if err := w.Flush(); err != nil {
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
	}
	f, err := os.Open(exname)
	if err != nil {
		return err
	}
	if err := c.ParseReader(bufio.NewReader(f)); err != nil {
		return err
	}
	return f.Close()
}
