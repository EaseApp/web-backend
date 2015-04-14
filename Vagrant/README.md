Vagrant setup
-------------
##Prerequisites
* Install [Virtual Box](https://www.virtualbox.org/wiki/Downloads)
* Install [Vagrant](https://www.vagrantup.com/downloads.html)

##Spin up the VM
Run `vagrant up` within this directory.
Initial installation and provisioning will take a while.

##Spin up the server
One provisioning is complete, ssh into the vagrant box.

`vagrant ssh`

Navigate to the projects root.

`cd /home/vagrant/go/src/github.com/easeapp/web-backend/`

Start the server.

`go run server.go`

##Go say hello

Visit `localhost:3000` to see the hello world.
