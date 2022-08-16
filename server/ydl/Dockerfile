FROM python:3.11-rc-bullseye

RUN set -eux && \
    apt-get update && \
    apt-get install -y --no-install-recommends ffmpeg && \
    rm -rf /var/lib/apt/lists/*


ENV TZ=Asia/Tokyo
COPY ydl /root/ydl
WORKDIR /root
EXPOSE 3000 3000
ENTRYPOINT ["/root/ydl"]
# CMD ["/root/ydl"]
