FROM ubuntu:16.04

RUN apt-get update -qq && apt-get upgrade -yq\
    && apt-get install -yq curl git nano \
    && curl -sL https://deb.nodesource.com/setup_8.x | bash \
    && apt-get install -y nodejs build-essential

ADD . /opt
RUN cd /opt && npm install

WORKDIR /opt

EXPOSE 8080

CMD ["npm", "start"]