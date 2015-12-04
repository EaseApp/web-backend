# Ease Web Backend

![Build Status](https://travis-ci.org/EaseApp/web-backend.svg?branch=master)

##Set up your dev enviroment
Follow the instructions on the README.md in the vagrant directory.

##Quick Start
ssh into the vagrant box.

`vagrant ssh`

Navigate to the project root.

`cd /home/vagrant/go/src/github.com/EaseApp/web-backend/`

Start the server.

`go run server.go`

If you wish to have the server restart on file change, run the server with `gin` in the `web-backend` directory:  

`gin`


##Go say hello

Visit `localhost:3000` to see the hello world.

# Data API Documentation

The data API is how users interact with their application data.  The REST API can be accessed from any language or runtime.  Ease also has client libraries for [JavaScript](https://github.com/EaseApp/javascript-client), [Java (Android)](https://github.com/EaseApp/java-client), and [Swift (iOS)](https://github.com/EaseApp/ios-client) to make accessing these endpoints easier.

Application data is stored in a JSON format.  
The API consists of three main endpoints: read, save, and delete.  

## Read Endpoint

### GET /data/{username}/{app_name}

#### Parameters:

`path`:  This is the path of where to read the data.

Example:
```
curl -H "Authorization: JzlaHSLCZdFDqKjYLonmyjFkhXFkYY" \
'http://api.easeapp.co/data/easetestuser@example.com/testapp?path=/'

```


## Write Endpoint

### POST /data/{username}/{app_name}

#### Parameters:

The endpoint takes two JSON parameters:

`path`: This is the path of where to save the data.  
`data`: This is the data to save.  It can be any JSON data.

Example:
```
curl -H 'Authorization: JzlaHSLCZdFDqKjYLonmyjFkhXFkYY' \
'http://api.easeapp.co/data/easetestuser@example.com/testapp' -X POST \
-d '{"path": "/messages","data": ["Hello, world!", "Welcome to Ease!"]}'
```

## Delete Endpoint

### POST /data/{username}/{app_name}

#### Parameters:

The delete endpoint takes one JSON parameter.

`path`: This is the JSON path of the data to be deleted.

Example:
```
curl -H 'Authorization: JzlaHSLCZdFDqKjYLonmyjFkhXFkYY' \
'http://api.easeapp.co/data/easetestuser@example.com/testapp' -X DELETE \
-d '{"path": "/messages"}'
```
