package argparser

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Parser struct {
	generate       *bool
	name           *string
	commands       []string
	CurrentCommand string
}

func NewParser() *Parser {
	return &Parser{
		generate:       nil,
		name:           nil,
		commands:       make([]string, 0),
		CurrentCommand: "",
	}
}

func (p *Parser) Arg(i int) string {
	return flag.Arg(i)
}

func (p *Parser) NArg() int {
	return flag.NArg()
}

func (p *Parser) addCommonFlags() {
	p.generate = flag.Bool("generate", false, "Generating default config")
	p.name = flag.String("config", "config.json", "Config name")
}

func (p *Parser) String(name string, value string, usage string) *string {
	return flag.String(name, value, usage)
}

func (p *Parser) Bool(name string, value bool, usage string) *bool {
	return flag.Bool(name, value, usage)
}

func (p *Parser) Int(name string, value int, usage string) *int {
	return flag.Int(name, value, usage)
}

func (p *Parser) parseConfig(config interface{}) (interface{}, error) {
	file, err := os.Open(*p.name)
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

func (p *Parser) generateConfig(config interface{}) error {
	_, err := os.Lstat(*p.name)
	if !os.IsNotExist(err) {
		return errors.New("File already exist")
	}

	file, err := os.Create(*p.name)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "    ")
	return enc.Encode(config)
}

func (p *Parser) AddCommands(commands []string) {
	p.commands = commands
}

func (p *Parser) isCommandsValid() bool {
	return len(p.commands) > 0
}

func (p *Parser) getCurrentCommand() error {
	if !p.isCommandsValid() {
		return nil
	}

	if p.NArg() <= 0 {
		return errors.New("No command specified")
	}

	p.CurrentCommand = p.Arg(0)

	for _, c := range p.commands {
		if c == p.CurrentCommand {
			return nil
		}
	}

	return fmt.Errorf("Invalid command was specified %s. Available commands: %s", p.CurrentCommand, strings.Join(p.commands, ", "))
}

func (p *Parser) Arguments() []string {
	startIndex := 0
	if p.isCommandsValid() {
		startIndex = 1
	}

	args := make([]string, 0, p.NArg())
	for i := startIndex; i < p.NArg(); i++ {
		args = append(args, p.Arg(i))
	}

	return args
}

func (p *Parser) Parse(config interface{}) (interface{}, error, bool) {
	p.addCommonFlags()
	flag.Parse()

	if *p.generate {
		err := p.generateConfig(config)
		return config, err, true
	}

	err := p.getCurrentCommand()
	if err != nil {
		return config, err, false
	}

	config, err = p.parseConfig(config)
	return config, err, false
}
