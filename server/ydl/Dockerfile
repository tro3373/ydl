FROM python:3.11-rc-bullseye

RUN set -eux && \
    apt-get update && \
    apt-get install -y --no-install-recommends ffmpeg && \
    rm -rf /var/lib/apt/lists/*
