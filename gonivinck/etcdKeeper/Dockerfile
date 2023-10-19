FROM alpine:latest

ENV HOST="0.0.0.0"
ENV PORT="8080"

WORKDIR /opt
RUN apk add git
RUN mkdir etcdkeeper

WORKDIR /opt/etcdkeeper
RUN git clone https://github.com/evildecay/etcdkeeper.git
EXPOSE ${PORT}

WORKDIR /opt/etcdkeeper
CMD ./etcdkeeper -h $HOST -p $PORT -atuh=true
