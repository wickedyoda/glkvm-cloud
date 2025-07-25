FROM alpine:latest
WORKDIR /home
COPY ./rttys /usr/bin/rttys
ENTRYPOINT ["/usr/bin/rttys"]
