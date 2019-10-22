# Dockerfile
FROM ubuntu:latest
MAINTAINER ish <ish@innogrid.com>

RUN mkdir -p /GraphQL_harp/
WORKDIR /GraphQL_harp/

ADD GraphQL_harp /GraphQL_harp/
RUN chmod 755 /GraphQL_harp/GraphQL_harp

EXPOSE 8001

CMD ["/GraphQL_harp/GraphQL_harp"]
