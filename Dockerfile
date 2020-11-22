# Dockerfile
FROM ubuntu:latest
MAINTAINER ish <ish@innogrid.com>

RUN mkdir -p /harp/
WORKDIR /harp/

ADD harp /harp/
RUN chmod 755 /harp/harp

EXPOSE 7000

CMD ["/harp/harp"]
