FROM golang:latest
MAINTAINER LamdaLamdaLamda
LABEL version="0.1"
LABEL name="githubUpdates"

WORKDIR /app
RUN apt install -y git
CMD command -v go || echo "Build failed..."; exit 1;
RUN git clone https://github.com/Lisiadito/githubUpdates
WORKDIR /app/githubUpdates
RUN export GOPATH=$GOPATH:/app/githubUpdates
ENV TELEGRAM_API_TOKEN SOME_API_KEY
ENV GITHUB_API_TOKEN SOME_API_KEY
RUN go get -d ./...
RUN go build main.go
ENTRYPOINT bash
