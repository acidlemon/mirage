mirage - reverse proxy frontend for docker
===========================================

mirage is reverse proxy for docker container and container manager.

mirage can launch and terminate docker container and serve http request
with specified subdomain. Additionaly, mirage passes variable to Dockerfile
using environment variables.

Usage
------

1. Setup mirage and edit configuration (see Setup section for detail.)
2. Run mirage.

Following instructions use below settings.

```
host:
  webapi: docker.dev.example.net
  reverse_proxy_suffix: .dev.example.net
listen:
  HTTP: 80
```

Prerequisite: you should resolve `*.dev.example.net` to your docker host.

### Using CLI

3. Launch docker container using curl.
```
curl http://docker.dev.example.net/api/launch \
  -d subdomain=cool-feature \
  -d branch=feature/cool \
  -d image=myapp:latest
```
4. Now, you can access to container using "http://cool-feature.dev.exmaple.net/".

5. Terminate container using curl.
```
curl http://docker.dev.example.net/api/terminate \
  -d subdomain=cool-feature
```

### Using Web Interface

3. Access to mirage web interface via "http://docker.dev.example.net/".
4. Press "Launch New Container".
5. Fill launch options.
  - subdomain: cool-feature
  - branch: feature/cool
  - image: myapp:latest
6. Now, you can access to container using "http://cool-feature.dev.exmaple.net/".
7. Press "Terminate" button.


Currently, v0.1 only support "branch" variable using `GIT_BRANCH`
environment variables.


Setup
------

(write in progress :)


License
--------

The MIT License (MIT)

(c) 2014 acidlemon. (c) 2014 KAYAC Inc.



