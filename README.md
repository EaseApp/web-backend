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

The data API is how users interact with their application data.  The REST API can be accessed from any language or runtime.  Ease also has client libraries for [JavaScript](https://github.com/EaseApp/javascript-client), [Java (Android)](https://github.com/EaseApp/android-client), and [Swift (iOS)](https://github.com/EaseApp/ios-client) to make accessing these endpoints easier.

Application data is stored in JSON format. It can be accessed using JSON paths in the form `/users/ryan/messages`.  Each section of the URL corresponds to a key in the JSON object.  

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

## Sync Endpoint for Realtime Updates

The sync service allow a user to get real time updates via websockets. The Javascript library has abstracted the websocket connection, so you can use the connect() function. If you want to manually connect to Sync, first open make a TCP websocket connection to ws://sync.easeapp.co:8000. After successful connection, send a string representation of a JSON document with the following three attributes. 

```
{
  "username":"easetestuser@example.com",
  "appName":"testapp",
  "authorization":"JzlaHSLCZdFDqKjYLonmyjFkhXFkYY"
}
```
| Attribute | Description |
|:--|:--|
| username | The username that the application is registered under |
| appName | The name of the application that you want to get notifications for through the sync service |
| authorization | The application token for the associated application, this is used to authenticate the connection |

On successful authentication, the connection server will send back:
```
{
  "status":"successful"
}
```

The sync service will begin to pass back the sync information back as things are updated on the application in the following format:
```
{
  "action": "SAVE",
  "data": yourData, 
  "path": "/"
}
```

| Attribute | Description |
|:--|:--|
| action | The operation that was performed. Valid actions are "SAVE" and "DELETE" |
| path | The location in which the data was affected |
| data | The data that is now stored in that path |

# Example Application

To see an example of Ease in action, check out the Ease [web chat repo](https://github.com/EaseApp/chat-repo).  To see it running live, check it out [here](http://easeapp.github.io/chat-demo/).
