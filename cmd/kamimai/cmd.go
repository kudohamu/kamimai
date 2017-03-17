package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kudohamu/kamimai/core"
)

type (
	// A Cmd executes a command
	Cmd struct {
		Name  string
		Usage string
		Run   func(*Cmd, ...string) error

		flag flag.FlagSet
	}
)

var (
	version uint64
)

// Exec executes a command with arguments.
func (c *Cmd) Exec(args []string) error {
	c.flag.Uint64Var(&version, "version", 0, "")
	c.flag.Parse(args)

	// Load config
	config, err := mustReadConfig(*dirPath)
	if err != nil {
		panic(err)
	}
	config.WithEnv(*env)
	return c.Run(c, c.flag.Args()...)
}

func mustReadConfig(dir string) (*core.Config, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		fpath := filepath.Join(dir, f.Name())

		switch filepath.Ext(fpath) {
		case ".yml", ".yaml":
			fi, err := os.Open(fpath)
			if err != nil {
				return nil, err
			}
			defer fi.Close()
			return core.NewConfigFromYML(fi, dir)
		case ".tml", ".toml":
			fi, err := os.Open(fpath)
			if err != nil {
				return nil, err
			}
			defer fi.Close()
			return core.NewConfigFromToml(fi, dir)
		}
	}
	return nil, errors.New("config file not found")
}
