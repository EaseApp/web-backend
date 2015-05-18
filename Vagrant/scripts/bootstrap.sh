sudo apt-get -y install git;
sudo apt-get -y install mercurial;

# Install Go 1.4.2
wget --quiet https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.4.2.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> /home/vagrant/.profile;
export PATH=$PATH:/usr/local/go/bin;

source /etc/lsb-release && echo "deb http://download.rethinkdb.com/apt $DISTRIB_CODENAME main" | sudo tee /etc/apt/sources.list.d/rethinkdb.list;
wget -qO- http://download.rethinkdb.com/apt/pubkey.gpg | sudo apt-key add -;
sudo apt-get update;
sudo apt-get -y install rethinkdb;
