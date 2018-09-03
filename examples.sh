manageAMQ bulkinsert -u admin -p admin -q "track4j.TrackingJSV1" -d parsed/sorted/advInit/
manageAMQ -u admin -p admin bulkinsert -q "track4j.TrackingJSV1" -m 10 -f ~/tracking/parsed/sorted/advInit/2018-08-10.log
