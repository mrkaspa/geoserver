# GeoServer

The Geoserver allows to do geolocation operations as a service through a REST API interface. Geoserver is a stand alone application that can be deployed into a server.

## Build

The Geoserver is built in Go so you need to make sure to have all the Go tooling installed and then download the package:

```
go get github.com/mrkaspa/geoserver
```

Go to the package through the console and use go build or go install commands to generate the binary.

## Configuration

Create a file called .env and place it in the same directory of the geoserver binaray. The content of this file must be:

```
MONGO_URI=mongodb://localhost:27017/geoserver
MONGO_DB=geoserver
PORT=8080
HOST=localhost
MODE=production
```

We need the mongo URI, the name of the database, the port and host for the server and the environment for the application.

## API

The API is divided by in two, a tipical REST API and a WEBSOCKET service.

### REST



### Websocket
