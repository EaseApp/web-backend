# Ease Web Backend

##Set up your dev enviroment
Follow the instructions on the README.md in the vagrant directory.

##Quick Start
ssh into the vagrant box.

`vagrant ssh`

Navigate to the project root.

`cd /home/vagrant/go/src/github.com/easeapp/web-backend/`

Start the server.

`go run server.go`

If you wish to have the server restart on file change, run the server with the `run_dev_server.sh` script in the `web-backend` directory:

`./run_dev_server.sh`


##Go say hello

Visit `localhost:3000` to see the hello world.
