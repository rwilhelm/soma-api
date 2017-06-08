CONST_DEVICE_ID = 0eb0bed2-4068-11e7-b5ec-1a16a99bea8b
CONST_TRIP_UUID = 0b5841b2-4bdd-11e7-9b19-1a16a99bea8b

DEVICE_ID = $(shell uuid)
TRIP_UUID = $(shell uuid)

new-trip:
	echo '{"device_id":"'$(CONST_DEVICE_ID)'","locationData":[{"latitude":50.36353176091952,"timestamp":1495461615.936066,"speed":-1,"accuracy":10,"altitude":83.53457641601562,"bearing":-1,"longitude":7.558246505952352}],"uuid":"'$(TRIP_UUID)'"}' | curl -k -H "Content-Type: application/json" -d @- https://soma.uni-koblenz.de:5000/upload

new-device:
	echo '{"device_id":"'$(DEVICE_ID)'","locationData":[{"latitude":50.36353176091952,"timestamp":1495461615.936066,"speed":-1,"accuracy":10,"altitude":83.53457641601562,"bearing":-1,"longitude":7.558246505952352}],"uuid":"'$(TRIP_UUID)'"}' | curl -k -H "Content-Type: application/json" -d @- https://soma.uni-koblenz.de:5000/upload

duplicate:
	echo '{"device_id":"'$(CONST_DEVICE_ID)'","locationData":[{"latitude":50.36353176091952,"timestamp":1495461615.936066,"speed":-1,"accuracy":10,"altitude":83.53457641601562,"bearing":-1,"longitude":7.558246505952352}],"uuid":"'$(CONST_TRIP_UUID)'"}' | curl -k -H "Content-Type: application/json" -d @- https://soma.uni-koblenz.de:5000/upload

duplicate-trip:
	echo '{"device_id":"'$(DEVICE_ID)'","locationData":[{"latitude":50.36353176091952,"timestamp":1495461615.936066,"speed":-1,"accuracy":10,"altitude":83.53457641601562,"bearing":-1,"longitude":7.558246505952352}],"uuid":"'$(CONST_TRIP_UUID)'"}' | curl -k -H "Content-Type: application/json" -d @- https://soma.uni-koblenz.de:5000/upload
