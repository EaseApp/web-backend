cp /etc/rethinkdb/default.conf.sample /etc/rethinkdb/instances.d/instance1.conf

echo "" >> /etc/rethinkdb/instances.d/instance1.conf
# Make port-forwarding for the rethinkdb admin work.
echo "bind=0.0.0.0" >> /etc/rethinkdb/instances.d/instance1.conf

sudo /etc/init.d/rethinkdb restart

