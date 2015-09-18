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

Install the dependencies.
`make dependencies`

 Your dev environment should be ready to go!
