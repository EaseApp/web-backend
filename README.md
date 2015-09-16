# Ease Web Backend

![Build Status](https://travis-ci.org/EaseApp/web-backend.svg?branch=master)

## Set up your dev enviroment
Follow the instructions on the README.md in the vagrant directory.

#### Quick Start
ssh into the vagrant box.

`vagrant ssh`

Navigate to the project root.

`cd /home/vagrant/go/src/github.com/easeapp/web-backend/`

Start the server.

`go run server.go`

If you wish to have the server restart on file change, run the server with the `run_dev_server.sh` script in the `web-backend` directory:

`./run_dev_server.sh`


#### Go say hello

Visit `localhost:3000` to see the hello world.


## API Specification

#### Login Methods

TODO - I'm not sure how oauth with a frontend framework will work.

#### Data Methods

**POST /data**

Parameters:
  - string path (ex: `'/users/joe'`)
  - JSON data (ex: `'{"name":"Guilliams"}'`)

Returns:
  - Success or failure.
  - Possibly the updated JSON.

This API method saves the given data to the given path.

**GET /data**

Parameters:
  - string path

Returns:
  - Success or failure.
  - JSON of the data at the given path.
  
This API method gets the data at the given path.
