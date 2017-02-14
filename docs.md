Use the rancher plugin to execute commands inside running service service in [rancher](http://rancher.com).

The following parameters are used to configure this plugin:

- `url` - url to your rancher server, including protocol and port
- `rancher_access_key` - rancher api access key
- `rancher_secret_key` - rancher api secret key
- `service` - name of rancher service to act on
- `cmd` - command to be executed inside the container
- `expect` - string to search for inside of the command execution output
- `exec-timeout` - timeout for execution to take place before failure

The following is a sample Rancher configuration in your `.drone.yml` file:

```yaml
deploy:
  rancher:
    url: https://example.rancher.com
    access_key: 1234567abcdefg
    secret_key: abcdefg1234567
    service: drone/drone
    cmd: 'ls -la'
```

if you want to add secrets for the access_key and secret it's RANCHER_ACCESS_KEY and RANCHER_SECRET_KEY


Note that if your `service` is part of a stack, you should use the notation `stackname/servicename` as this will make sure that the found service is part of the correct stack. If no stack is specified, this plugin will update the first service with a matching name which may not be what you want.

Note that if the service contains multiple containers the container the command will be executed in is arbitrary.
