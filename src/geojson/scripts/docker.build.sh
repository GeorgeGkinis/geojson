docker build --rm -f $(pwd)/../service-A/Dockerfile -t georgegkinis/geojson:a $(pwd)/../service-A/
docker build --rm -f $(pwd)/../service-B/Dockerfile -t georgegkinis/geojson:b $(pwd)/../service-B/
docker build --rm -f $(pwd)/../service-C/Dockerfile -t georgegkinis/geojson:c $(pwd)/../service-C/