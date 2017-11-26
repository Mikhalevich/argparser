package argparser

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
)

var (
	generate *bool
	name     *string
)

func Arg(i int) string {
	return flag.Arg(i)
}

func NArg() int {
	return flag.NArg()
}

func addCommonFlags() {
	generate = flag.Bool("generate", false, "Generating default config")
	name = flag.String("config", "config.json", "Config name")
}

func String(name string, value string, usage string) *string {
	return flag.String(name, value, usage)
}

func Bool(name string, value bool, usage string) *bool {
	return flag.Bool(name, value, usage)
}

func Int(name string, value int, usage string) *int {
	return flag.Int(name, value, usage)
}

func parseConfig(config interface{}) (interface{}, error) {
	file, err := os.Open(*name)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return nil, err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	err = dec.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func generateConfig(config interface{}) error {
	_, err := os.Lstat(*name)
	if !os.IsNotExist(err) {
		return errors.New("File already exist")
	}

	file, err := os.Create(*name)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	return enc.Encode(config)
}

func Parse(config interface{}) (interface{}, error, bool) {
	addCommonFlags()
	flag.Parse()

	if *generate {
		err := generateConfig(config)
		return config, err, true
	}

	config, err := parseConfig(config)
	return config, err, false
}
