package main

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
)

type Rancher struct {
	Url            string `json:"url"`
	AccessKey      string `json:"access_key"`
	SecretKey      string `json:"secret_key"`
	Service        string `json:"service"`
        Command        string `json:"command"`
}

var version string // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "rancher execute"
	app.Usage = "rancher execute"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{

		cli.StringFlag{
			Name:   "url",
			Usage:  "url to the rancher api",
			EnvVar: "PLUGIN_URL",
		},
		cli.StringFlag{
			Name:   "access-key",
			Usage:  "rancher access key",
			EnvVar: "PLUGIN_ACCESS_KEY, RANCHER_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "secret-key",
			Usage:  "rancher secret key",
			EnvVar: "PLUGIN_SECRET_KEY, RANCHER_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   "service",
			Usage:  "Service to act on",
			EnvVar: "PLUGIN_SERVICE",
		},
		cli.StringFlag{
			Name:   "command",
			Usage:  "command to execute",
			EnvVar: "PLUGIN_COMMAND",
		},
                cli.StringFlag{
                        Name:   "expect",
                        Usage:  "string to search for in the returned response",
                        EnvVar: "PLUGIN_EXPECT",
                },
                cli.IntFlag{
                        Name:   "exec-timeout",
                        Usage:  "Timeout for command to execute",
                        EnvVar: "PLUGIN_EXEC_TIMEOUT",
                },
		cli.BoolTFlag{
			Name:   "yaml-verified",
			Usage:  "Ensure the yaml was signed",
			EnvVar: "DRONE_YAML_VERIFIED",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		URL:            c.String("url"),
		Key:            c.String("access-key"),
		Secret:         c.String("secret-key"),
		Service:        c.String("service"),
		Command:        c.String("command"),
                Expect:         c.String("expect"),
                ExecTimeout:    c.Int("exec-timeout"),
		YamlVerified:   c.BoolT("yaml-verified"),
	}
	return plugin.Exec()
}
