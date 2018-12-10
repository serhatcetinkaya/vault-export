package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type authConfig struct {
	Token     string `yaml:"token"`
	VaultAddr string `yaml:"vault_addr"`
}

func main() {
	key := flag.String("k", "", "key of secret to be retrieved")
	flag.Parse()

	if *key == "" {
		fmt.Println("Must provide key")
		return
	}

	cfg, err := readConfig()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v", err)
		return
	}

	c, err := api.NewClient(&api.Config{
		Address: cfg.VaultAddr,
	})
	if err != nil {
		fmt.Printf("Failed to create Vault client: %v", err)
		return
	}

	c.SetToken(cfg.Token)

	readSecrets(c, *key)
}

func readSecrets(cli *api.Client, path string) {
	if rune(path[len(path)-1]) == rune('/') {
		sec, err := cli.Logical().List(path)
		if err != nil {
			fmt.Printf("Failed to get secret: %v", err)
			return
		}
		for _, v := range sec.Data {
			convertedData := convert(v)
			for _, i := range convertedData {
				if len(i) > 0 {
					defer readSecrets(cli, path+i)
				}
			}
		}
	} else {
		sec, err := cli.Logical().Read(path)
		if err != nil {
			fmt.Printf("Failed to get secret: %v", err)
			return
		}
		fmt.Printf("vault write " + path)
		for k, v := range sec.Data {
			fmt.Printf(" %s=%q", k, v)
		}
		fmt.Printf("\n")
	}
}

func readConfig() (authConfig, error) {
	configPath := os.Getenv("VAULT_EXPORTER_CONFIG_FILE")
	if configPath == "" {
		configPath = ".auth.yaml"
	}

	bs, err := ioutil.ReadFile(configPath)
	if err != nil {
		return authConfig{}, errors.Wrap(err, "failed to read configuration file")
	}

	var cfg authConfig
	if err := yaml.Unmarshal(bs, &cfg); err != nil {
		return authConfig{}, errors.Wrap(err, "failed to parse configuration file")
	}

	return cfg, nil
}

func convert(t interface{}) []string {
	res := []string{}
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)
		for i := 0; i < s.Len(); i++ {
			res = append(res, fmt.Sprintf("%v", s.Index(i)))
		}
	}
	return res
}
