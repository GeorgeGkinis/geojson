# geojson

#### Run in Vagrant
To recreate the development environment run :
```
vagrant up --provision
```

After a while the services should be running.

#### Run in Docker

To start the services in your own environment run :

```
docker run -d --rm --name service-C --net="host" -p 8091:8091 -p 8090:8090  georgegkinis/geojson:c

docker run -d --rm --name service-B --net="host" -p 8080:8080  georgegkinis/geojson:b

docker run -d --rm --name service-A --net="host" georgegkinis/geojson:a
```

The utility script src/geogson/scripts/docker.run.sh does the same.


#### Run as executables

```
cd src/geojson/scripts

./go.build.sh

cd ../service-C && SERVER_C_URL=0.0.0.0:8091 ./service-C &
cd ../service-B && nohup ./service-B &
cd ../service-A && nohup ./service-A &
# Also run this comment, needs an extra \r ;)

```

To see a plot of the GeoJSON data visit [localhost:8091](http://localhost:8091) in your browser.

The country with the largest population density is Singapore:
 
 **64589 People/KM^2**

Next steps:

- run as a Kubernetes Cluster
- multiples of service-B
- calculation of AREA based from multipolygons when AREA==0

N-Joy and thanks!