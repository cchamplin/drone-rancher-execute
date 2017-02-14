package main

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
        "fmt"
        "encoding/json"
        "bytes"
)

type Rancher struct {
	Url            string `json:"url"`
	AccessKey      string `json:"access_key"`
	SecretKey      string `json:"secret_key"`
	Service        string `json:"service"`
        Command        string `json:"command"`
}

var version string // build number set at compile-time

func standardArgs(args []string) {
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
                        Name:   "access-key,access_key",
                        Usage:  "rancher access key",
                        EnvVar: "PLUGIN_ACCESS_KEY, RANCHER_ACCESS_KEY",
                },
                cli.StringFlag{
                        Name:   "secret-key,secret_key",
                        Usage:  "rancher secret key",
                        EnvVar: "PLUGIN_SECRET_KEY, RANCHER_SECRET_KEY",
                },
                cli.StringFlag{
                        Name:   "service",
                        Usage:  "Service to act on",
                        EnvVar: "PLUGIN_SERVICE",
                },
                cli.StringFlag{
                        Name:   "command,cmd",
                        Usage:  "command to execute",
                        EnvVar: "PLUGIN_COMMAND",
                },
                cli.StringFlag{
                        Name:   "expect",
                        Usage:  "string to search for in the returned response",
                        EnvVar: "PLUGIN_EXPECT",
                },
                cli.IntFlag{
                        Name:   "exec-timeout,exec_timeout",
                        Usage:  "Timeout for command to execute",
                        EnvVar: "PLUGIN_EXEC_TIMEOUT",
                },
                cli.BoolTFlag{
                        Name:   "yaml-verified",
                        Usage:  "Ensure the yaml was signed",
                        EnvVar: "DRONE_YAML_VERIFIED",
                },
        }

        if err := app.Run(args); err != nil {
                log.Fatal(err)
        }

}

func legacyArgs() ([]string,bool) {
        var buf *bytes.Buffer
        hasArgs := false
        for i, argv := range os.Args {
		if argv == "--" {
			arg := os.Args[i+1]
			buf = bytes.NewBufferString(arg)
                        hasArgs = true
			break
		}
	}
        if hasArgs {
          var data map[string]interface{}
          if err := json.Unmarshal(buf.Bytes(),&data); err != nil {
            return nil, false
          }
          if _, ok := data["vargs"]; !ok {
            return nil, false
          }
          var vargs map[string]interface{}
          vargs = data["vargs"].(map[string]interface{})
          var newArgs []string
          newArgs = make([]string,(len(vargs)*2)+1)
          newArgs[0] = os.Args[0]
          idx := 1
          for k,v := range vargs {
            newArgs[idx] = fmt.Sprintf("%s%s","--",k)
            switch val := v.(type) {
              default:
                newArgs[idx+1] = fmt.Sprintf("%#v",val)
              case string:
                newArgs[idx+1] = fmt.Sprintf("%s",val)
            }
            idx = idx + 2
          }
          return newArgs,true
        }
        return nil, false
       
}

func main() {
   newArgs,ok := legacyArgs()
   if ok {
    standardArgs(newArgs)
   } else {
    standardArgs(os.Args)
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
