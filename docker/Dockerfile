FROM ubuntu
MAINTAINER acidlemon

RUN apt-get update
RUN apt-get -y install wget jq curl unzip

WORKDIR /opt

RUN export MIRAGE_VERSION=`curl https://api.github.com/repos/acidlemon/mirage/releases | jq -r ".[0].tag_name"` && wget https://github.com/acidlemon/mirage/releases/download/${MIRAGE_VERSION}/mirage-${MIRAGE_VERSION}-linux-amd64.zip -O /opt/mirage.zip
RUN unzip mirage.zip
WORKDIR /opt/mirage

ADD config.yml /opt/mirage/config.yml

CMD ./mirage

