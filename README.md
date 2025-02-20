# go-rpc-distance-app

## Project description

The project implements interaction between golang services. The project represents application to calculate distance between geo points in unlimited quantity. Frontend application on localhost:8000 shows form for geo points. JS AJAX sends it to server on same port that sends data to RPC-server via tcp-connection. Then RPC microservice consumes these geo data and calculates distance between geo points, then it returnes distance data back to frontend application. In another goroutine it writes the geo points and distance to the table in mysql database. All the results are represented in localhost:8000/results/. Open page on localhost:8000 and calculate distance between geo points.

## Project setup

In Windows open C:/Windows/System32/drivers/etc/hosts or in Linux /etc/hosts and add lines for docker ip-addresses to escape mysql connection failures:

```
172.18.0.2 localhost
172.18.0.3 localhost
172.18.0.4 localhost
```

```
docker-compose up -d
```

Wait for several seconds to start database and services

#### Open project
```
localhost:8000

localhost:8000/results/
```

Open page and type geo points to calculate distance.
