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

#### Store a geolocation point

Send a POST request to the endpoint http://localhost:8080/store with a json like this:

```js
{
    "location": [-79.38066843, 43.65483486],
     "user_id": "demo@demo.com",
     "info": "info to share" 
}
```

Returns status 200 if saves the geolocation point

#### Store and retrieve near points

Send a POST request to the endpoint http://localhost:8080/near with a json like this:

```js
{
    "time_range": 5, // seconds
    "max_distance": 5, // meters
    "stroke": {
        "location": [-79.38066843, 43.65483486],
        "user_id": "demo@demo.com",
        "info": "info to share"
    } 
}
```

Returns status 200 if saves the geolocation point and the points found near. In the time_range you set a time limit to find users if you don't want to use this filter set it to 0, the max distance is a parameter to query the nearest users.

A response can look like this:

```js
[
    {
        "location": [-79.38066853, 43.65483586],
        "user_id": "another@demo.com",
        "info": "info to share"
    }
]
```

