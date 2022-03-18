FROM golang:1.17

RUN mkdir app
COPY . /app

WORKDIR /app
RUN mkdir build

RUN  apt update -y
RUN  apt install ffmpeg -y
RUN apt install jq -y
RUN bash -c "url=$(curl https://api.github.com/repos/yt-dlp/yt-dlp/releases | jq -r '.[0] | .assets | .[2] |.browser_download_url'); curl -L -o yt-dlp \$url"

RUN mkdir build/tool
RUN mv ./yt-dlp build/tool
RUN chmod +x /app/build/tool/yt-dlp

ENV YTD_PATH=/app/build/tool/yt-dlp


RUN go build -o /build/byd -ldflags="-s -w"  .

ENTRYPOINT ["/build/byd"]
