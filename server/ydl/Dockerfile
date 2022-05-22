FROM alpine:latest

ARG PKGS

RUN deps="ffmpeg python3 libc6-compat" && \
    apk add -q --progress --update --no-cache ${deps} && \
    ln -s /usr/bin/python3 /usr/bin/python

COPY ydl /usr/local/bin/

WORKDIR /root

ENTRYPOINT ["ydl"]
