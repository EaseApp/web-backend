sudo apt-get -y install git;
sudo apt-get -y install mercurial;

# Install Go 1.5
wget --quiet https://storage.googleapis.com/golang/go1.5.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.5.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> /home/vagrant/.profile;
export PATH=$PATH:/usr/local/go/bin;

# Set up GOPATH.
echo "export GOPATH=/home/vagrant/go" >> /home/vagrant/.profile;

# Install inotify-tools for live code reload.
sudo apt-get -y install inotify-tools

# Make the dev server script executable.
chmod a+x /home/vagrant/go/src/github.com/EaseApp/web-backend/run_dev_server.sh;

source /etc/lsb-release && echo "deb http://download.rethinkdb.com/apt $DISTRIB_CODENAME main" | sudo tee /etc/apt/sources.list.d/rethinkdb.list;
wget -qO- http://download.rethinkdb.com/apt/pubkey.gpg | sudo apt-key add -;
sudo apt-get update;
sudo apt-get -y install rethinkdb;

