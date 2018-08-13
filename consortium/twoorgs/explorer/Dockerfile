FROM node:8.11.3-alpine

ADD . /opt/explorer

RUN apk add --no-cache --virtual npm-deps python make g++ && \
    python -m ensurepip && \
    rm -r /usr/lib/python*/ensurepip && \
    pip install --upgrade pip setuptools && \
	rm -r /root/.cache

RUN npm install -g yarn
RUN cd /opt/explorer && npm install

RUN apk del npm-deps

RUN cd /opt/explorer/client && npm install && npm run build

EXPOSE 8080

CMD node /opt/explorer/main.js