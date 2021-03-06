package main

import (
	"os"

	"github.com/gbrlsnchs/cli"
	"github.com/gbrlsnchs/pilgo/cmd/internal"
	"github.com/gbrlsnchs/pilgo/config"
	"github.com/gbrlsnchs/pilgo/fs"
	"github.com/gbrlsnchs/pilgo/fs/fsutil"
)

type appConfig struct {
	conf          string
	fs            fs.Driver
	getwd         func() (string, error)
	userConfigDir func() (string, error)
	version       string
}

func (cfg *appConfig) copy() appConfig { return *cfg }

type rootCmd struct {
	// store
	check   checkCmd
	config  configCmd
	init    initCmd
	link    linkCmd
	show    showCmd
	version versionCmd
}

func main() {
	os.Exit(run())
}

func run() int {
	var (
		root   rootCmd
		appcfg = appConfig{
			fs:            fsutil.OSDriver{},
			getwd:         os.Getwd,
			userConfigDir: os.UserConfigDir,
			version:       internal.Version(),
		}
	)
	cli := cli.New(&cli.Command{
		Options: map[string]cli.Option{
			"config": cli.StringOption{
				OptionDetails: cli.OptionDetails{
					Description: "Use a different configuration file.",
				},
				DefValue:  config.DefaultName,
				Recipient: &appcfg.conf,
			},
		},
		Subcommands: map[string]*cli.Command{
			"check": {
				Description: "Check the status of your dotfiles.",
				Options: map[string]cli.Option{
					"fail": cli.BoolOption{
						OptionDetails: cli.OptionDetails{
							Short:       'f',
							Description: "Fail when there are one or more conflicts.",
						},
						DefValue:  false,
						Recipient: &root.check.fail,
					},
				},
				Exec: root.check.register(appcfg.copy),
			},
			"config": {
				Description: "Configure a dotfile in the configuration file.",
				Options: map[string]cli.Option{
					"basedir": cli.StringOption{
						OptionDetails: cli.OptionDetails{
							Description: "Set the file's base directory.",
							ArgLabel:    "DIR",
						},
						Recipient: &root.config.baseDir,
					},
					"link": cli.VarOption{
						OptionDetails: cli.OptionDetails{
							Description: "Set the file's link name. An empty string skips the file.",
							ArgLabel:    "NAME",
						},
						Recipient: &root.config.link,
					},
					"targets": cli.VarOption{
						OptionDetails: cli.OptionDetails{
							Description: "Comma-separated list of the file's targets.",
							ArgLabel:    "TARGET 1,...,TARGET n",
						},
						Recipient: &root.config.targets,
					},
				},
				Arg: cli.StringArg{
					Label:     "TARGET",
					Required:  false,
					Recipient: &root.config.file,
				},
				Exec: root.config.register(appcfg.copy),
			},
			"init": {
				Description: "Initialize a configuration file.",
				Options: map[string]cli.Option{
					"force": cli.BoolOption{
						OptionDetails: cli.OptionDetails{
							Description: "Override an already existing configuration file.",
							Short:       'f',
						},
						Recipient: &root.init.force,
					},
					"include": cli.VarOption{
						OptionDetails: cli.OptionDetails{
							Description: "Comma-separated list of targets to be included.",
							ArgLabel:    "TARGET 1,...,TARGET n",
						},
						Recipient: &root.init.include,
					},
					"exclude": cli.VarOption{
						OptionDetails: cli.OptionDetails{
							Description: "Comma-separated list of targets to be excluded.",
							ArgLabel:    "TARGET 1,...,TARGET n",
						},
						Recipient: &root.init.exclude,
					},
					"hidden": cli.BoolOption{
						OptionDetails: cli.OptionDetails{
							Description: "Include hidden files.",
							Short:       'H',
						},
						Recipient: &root.init.hidden,
					},
				},
				Exec: root.init.register(appcfg.copy),
			},
			"link": {
				Description: "Link your dotfiles as set in the configuration file.",
				Exec:        root.link.register(appcfg.copy),
			},
			"show": {
				Description: "Show your dotfiles in a tree view.",
				Exec:        root.show.register(appcfg.copy),
			},
			"version": {
				Description: "Print version.",
				Exec:        root.version.register(appcfg.copy),
			},
		},
	})
	return cli.ParseAndRun(os.Args)
}
