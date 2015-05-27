export GOPATH=/home/vagrant/go;
export PATH=$PATH:/usr/local/go/bin;

# Server dependencies.
go get -d github.com/dancannon/gorethink;
go get -d github.com/codegangsta/negroni;
go get -d github.com/gorilla/mux;

# Go command line tools.
echo "export PATH=$PATH:$GOPATH/bin;" >> /home/vagrant/.profile;
go get github.com/codegangsta/gin;
