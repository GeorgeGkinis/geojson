# geojson

To recreate the development environment run :
vagrant up --provision

Service A,B and C will start as containers.

To start the services in your own environment run :

```
docker run -d --rm --name service-C --net="host" -p 8091:8091 -p 8090:8090  georgegkinis/geojson:c

docker run -d --rm --name service-B --net="host" -p 8080:8080  georgegkinis/geojson:b

docker run -d --rm --name service-A --net="host" georgegkinis/geojson:a
```

The utility script src/geogson/scripts/docker.run.sh does the same.

To see a plot of the GeoJSON data visit [localhost:8091](http://localhost:8091) in your browser.

NJoy and thanks!