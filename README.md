# Prometheus example

This project is a project in order to make the first steps with
Prometheus. It is a simple golang http-server which exposes some
simple Prometheus metrics. The exposed metrics can be scraped and
queried by a Prometheus server.


## Prom Queries

The metrics can be visualised in a grafana dashboard. It is available
as `graphana_dashboards/dashboard.json` in this repository. To get all
sample queries use `jq`:

``` bash
cat dashboard.json| jq ".panels | map(.title, .targets)"
```


### Generate metrics

* You can generate metrics with the provided script
  `insert_posts.sh`. All provided http-endpoints contain a random
  component. I.e. the http-responses vary in response time and success
  rate. Errors happened with a with a specific chance. This makes the
  sample Grafana dashboard "more" realistic and you can differentiate
  between errors, good requests and timeouts etc.

* There is also a go-client which loops endlessly over the provided
  endpoints

```bash
# builds the server and the client
make

# run the server on localhost:8080
./server

# run the client, you stop it with CTRL-c
./client
```


## Development

The project comes with a simple `Makefile`.

``` bash
# build server and client
make

# lint source code
make lint

# run tests
make test

# clean up
make clean
```
