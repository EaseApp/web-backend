sudo apt-get -y install git;
sudo apt-get -y install mercurial;

sudo apt-get -y install golang;

source /etc/lsb-release && echo "deb http://download.rethinkdb.com/apt $DISTRIB_CODENAME main" | sudo tee /etc/apt/sources.list.d/rethinkdb.list;
wget -qO- http://download.rethinkdb.com/apt/pubkey.gpg | sudo apt-key add -;
sudo apt-get update;
sudo apt-get -y install rethinkdb;
