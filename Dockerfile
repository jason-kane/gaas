from gliderlabs/alpine
# FROM golang:alpine

RUN apk add --no-cache fortune
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN mkdir /images
ADD gaas /usr/local/sbin/gaas
ADD templates /templates
ADD guidry.png /

CMD ["gaas"]

EXPOSE 8080