package argparser

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
)

type Parser struct {
	Config     interface{}
	Name       *string
	isGenerate *bool
}

func (p *Parser) addCommonFlags() {
	p.isGenerate = flag.Bool("generate", false, "Generating default config")
	p.Name = flag.String("config", "config.log", "Config name")
}

func (p *Parser) parseConfig() error {
	file, err := os.Open(*p.Name)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, p.Config)
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) Parse() (interface{}, error) {
	p.addCommonFlags()
	flag.Parse()
	if err := p.parseConfig(); err != nil {
		return nil, err
	}

	return p.Config, nil
}
