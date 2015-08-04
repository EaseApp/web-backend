Vagrant setup
-------------
##Prerequisites
* Install [Virtual Box](https://www.virtualbox.org/wiki/Downloads)
* Install [Vagrant](https://www.vagrantup.com/downloads.html)

##Spin up the VM
Run `vagrant up` within this directory.
Initial installation and provisioning will take a while.

##Spin up the server
Once provisioning is complete, ssh into the vagrant box.

`vagrant ssh`

Navigate to the project root.

`cd /home/vagrant/go/src/github.com/EaseApp/web-backend/`

Start the server.

`go run server.go`

If you wish to have the server restart on file change, run the server with the `run_dev_server.sh` script in the `web-backend` directory:

`./run_dev_server.sh`


##Go say hello

Visit `localhost:3000` to see the hello world.

