package main

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/client"
        "github.com/gorilla/websocket"
	"strings"
        "encoding/base64"
	"time"
)

type Plugin struct {
	URL            string
	Key            string
	Secret         string
	Service        string
	Command        string
        ExecTimeout    int
        Expect         string
	YamlVerified   bool
}

func (p *Plugin) Exec() error {
	log.Info("Drone Rancher Execute Plugin built")

	if p.URL == "" || p.Key == "" || p.Secret == "" || p.Command == "" {
		return errors.New("Eek: Must have url, key, secret, command, and service definied")
	}

	var wantedService, wantedStack string
	if strings.Contains(p.Service, "/") {
		parts := strings.SplitN(p.Service, "/", 2)
		wantedStack = parts[0]
		wantedService = parts[1]
	} else {
		wantedService = p.Service
	}

	rancher, err := client.NewRancherClient(&client.ClientOpts{
		Url:       p.URL,
		AccessKey: p.Key,
		SecretKey: p.Secret,
	})
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create rancher client: %s\n :(", err))
	}

	var stackId string
	if wantedStack != "" {
		environments, err := rancher.Environment.List(&client.ListOpts{})
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to list rancher environments: %s\n", err))
		}
		for _, env := range environments.Data {
			if env.Name == wantedStack {
				stackId = env.Id
			}
		}
		if stackId == "" {
			return errors.New(fmt.Sprintf("Unable to find stack %s\n", wantedStack))
		}
	}
	services, err := rancher.Service.List(&client.ListOpts{})
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to list rancher services: %s\n", err))
	}
	found := false
	var service client.Service
	for _, svc := range services.Data {
		if svc.Name == wantedService && ((wantedStack != "" && svc.EnvironmentId == stackId) || wantedStack == "") {
			service = svc
			found = true
		}
	}
	if !found {
		return errors.New(fmt.Sprintf("Unable to find service %s\n", p.Service))
	}

        var instances client.ContainerCollection

        err = rancher.GetLink(service.Resource,"instances",&instances)

        if err != nil {
                 return errors.New(fmt.Sprintf("Failed to list rancher containers: %s\n",err))
        }
        instanceFound := false
        var container client.Container
        for _,cntr := range instances.Data {
          //fmt.Printf("%s\n",cntr.Name)
          instanceFound = true
          container = cntr
          break
        }
        if !instanceFound {
                return errors.New(fmt.Sprintf("Failed to locate instance of service container"))
        }
        containerExec := &client.ContainerExec {
          AttachStdin: true,
          AttachStdout: true,
          Tty: false,
          Command: []string{"sh","-c",p.Command},
        }
        access,err := rancher.Container.ActionExecute(&container,containerExec)
        if err != nil {
                 return errors.New(fmt.Sprintf("Failed to exec in container: %s\n",err))
        }
        //fmt.Printf("%s?token=%s\n",access.Url,access.Token)
        c, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s?token=%s",access.Url, access.Token), nil)
	if err != nil {
		fmt.Printf("Could not connect to container via api: %s\n", err)
	}
	defer c.Close()

	done := make(chan string)
        timeout := p.ExecTimeout
        if timeout <= 0 {
          timeout = 3
        }
        var completeMessage string
	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
                                if websocket.IsUnexpectedCloseError(err,websocket.CloseGoingAway,websocket.CloseNormalClosure) {
				  done <- fmt.Sprintf("Failed to read container response: %s\n", err)
                                }
				return
			}
                        byteData, err := base64.StdEncoding.DecodeString(string(message))
                        if err != nil {
                          done <- fmt.Sprintf("Failed to decode container response: %s\n", err)
                          return
                        }
                        completeMessage = completeMessage + fmt.Sprintf("%q",byteData)   
		}
	}()
        wsComplete := false
        for {
          select {
            case err, ok := <-done:
              if !ok {
                //return errors.New(err)
                wsComplete = true
              } else { 
                return errors.New(err)
              }
            case <-time.After(time.Duration(timeout)*time.Second):
              return errors.New(fmt.Sprintf("Execeeded timeout for command execution"))
          }
          if wsComplete {
            break;
          }
        }
        if p.Expect != "" {
          if !strings.Contains(completeMessage,p.Expect) {
            return errors.New(fmt.Sprintf("Command result did not contain expected string: \"%s\"",completeMessage))
          }
        }
        c.Close()
        log.Info(fmt.Sprintf("Executed command \"%s\" inside container %s for service %s\n",p.Command,container.Name,p.Service))
	return nil
}

