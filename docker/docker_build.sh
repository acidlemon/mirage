#!/bin/sh

# download mirage
LATEST_VERSION=$(curl https://api.github.com/repos/acidlemon/mirage/releases | python -c 'import sys; import json; print(json.loads(sys.stdin.read()))[0]["tag_name"]');
curl -L -o mirage.zip https://github.com/acidlemon/mirage/releases/download/${LATEST_VERSION}/mirage-${LATEST_VERSION}-linux-amd64.zip
unzip mirage.zip

# build
docker build --rm -t mirage .

# clean up
rm -rf mirage mirage.zip
