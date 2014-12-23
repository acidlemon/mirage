FROM ubuntu
MAINTAINER acidlemon

RUN apt-get update
RUN apt-get -y install git docker.io mercurial wget

# install go 1.4
RUN wget https://storage.googleapis.com/golang/go1.4.linux-amd64.tar.gz 
RUN tar -C /usr/local -xzf go1.4.linux-amd64.tar.gz

RUN mkdir /opt/mirage
RUN mkdir /opt/gopath
ENV GOPATH  /opt/gopath
ENV PATH /opt/gopath/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
WORKDIR /opt/mirage
RUN git clone https://github.com/acidlemon/mirage .
RUN go get -v github.com/acidlemon/mirage
RUN go build
ADD config.yml /opt/mirage/config.yml

CMD ./mirage


