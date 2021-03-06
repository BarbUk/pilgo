package main

import (
	"errors"
	"fmt"

	"github.com/gbrlsnchs/cli"
	"github.com/gbrlsnchs/pilgo/config"
	"github.com/gbrlsnchs/pilgo/fs"
	"github.com/gbrlsnchs/pilgo/linker"
	"github.com/gbrlsnchs/pilgo/parser"
	"gopkg.in/yaml.v3"
)

type checkCmd struct {
	fail bool
}

func (cmd *checkCmd) register(getcfg func() appConfig) cli.ExecFunc {
	return func(prg cli.Program) error {
		appcfg := getcfg()
		fs := fs.New(appcfg.fs)
		b, err := fs.ReadFile(appcfg.conf)
		if err != nil {
			return err
		}
		var c config.Config
		if yaml.Unmarshal(b, &c); err != nil {
			return err
		}
		baseDir, err := appcfg.userConfigDir()
		if err != nil {
			return err
		}
		cwd, err := appcfg.getwd()
		if err != nil {
			return err
		}
		var p parser.Parser
		tr, err := p.Parse(c, parser.BaseDir(baseDir), parser.Cwd(cwd), parser.Envsubst)
		if err != nil {
			return err
		}
		ln := linker.New(fs)
		if err = ln.Resolve(tr); err != nil {
			var cft *linker.ConflictError
			if errors.As(err, &cft) {
				if !cmd.fail {
					goto printtree
				}
				name := prg.Name()
				errw := prg.Stderr()
				for _, err := range cft.Errs {
					fmt.Fprintf(errw, "%s: %v\n", name, err)
				}
			}
			return err
		}
	printtree:
		fmt.Fprint(prg.Stdout(), tr)
		return nil
	}
}
