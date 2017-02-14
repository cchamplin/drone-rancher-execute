# drone-rancher-execute


Drone plugin to execute commands inside of a Rancher service. For the usage information and a listing of the available options please take a look at [the docs](DOCS.md).

## Binary

Build the binary using `make`:

```
make deps build
```


## Docker

Build the container using `make`:

```
make deps docker
```

### Example

## Usage

Build and deploy from your current working directory:

```
docker run --rm                          \
  -e PLUGIN_URL=<source>                 \
  -e PLUGIN_ACCESS_KEY=<key>     \
  -e PLUGIN_SECRET_KEY=<secret>  \
  -e PLUGIN_SERVICE=<service>            \  
  -e PLUGIN_COMMAND=<command>         \
  -v $(pwd):$(pwd)                       \
  -w $(pwd)                              \
  plugins/drone-rancher-execute 
```
