docker run -d --rm --name service-C --net="host" -p 8091:8091 -p 8090:8090  georgegkinis/geojson:c
docker run -d --rm --name service-B --net="host" -p 8080:8080  georgegkinis/geojson:b
docker run -d --rm --name service-A --net="host" georgegkinis/geojson:a