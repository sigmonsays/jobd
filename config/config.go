package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/sigmonsays/jobd/job"
	yaml "gopkg.in/yaml.v2"
)

// main configuration structure
type AppConfig struct {
	HttpAddr      string         `yaml:"http_addr"`
	DataDir       string         `yaml:"datadir"`
	Defaults      *Defaults      `yaml:"defaults"`
	ShellDefaults *job.ShellSpec `yaml:"shell_defaults"`
	Jobs          []*job.JobSpec
}

type Defaults struct {
	Keep int
}

func (c *AppConfig) LoadYaml(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	b := bytes.NewBuffer(nil)
	_, err = b.ReadFrom(f)
	if err != nil {
		return err
	}

	if err := c.LoadYamlBuffer(b.Bytes()); err != nil {
		return err
	}

	if err := c.FixupConfig(); err != nil {
		return err
	}

	return nil
}

func (c *AppConfig) LoadYamlBuffer(buf []byte) error {
	err := yaml.Unmarshal(buf, c)
	if err != nil {
		return err
	}
	return nil
}

func (me *AppConfig) PrintConfig() {
	d, err := yaml.Marshal(me)
	if err != nil {
		fmt.Println("Marshal error", err)
		return
	}
	fmt.Println("-- Configuration --")
	fmt.Println(string(d))
}

func GetDefaultConfig() *AppConfig {
	cfg := &AppConfig{}
	cfg.DataDir = "/srv/jobd"
	cfg.HttpAddr = ":8093"
	cfg.ShellDefaults = &job.ShellSpec{
		Timeout: 3600,
	}
	return cfg
}

func (c *AppConfig) LoadDefault() {
	*c = *GetDefaultConfig()
}

// after loading configuration this gives us a spot to "fix up" any configuration
// or abort the loading process
func (c *AppConfig) FixupConfig() error {
	// var emptyConfig AppConfig

	return nil
}

func PrintDefaultConfig() {
	conf := GetDefaultConfig()
	conf.PrintConfig()
}
