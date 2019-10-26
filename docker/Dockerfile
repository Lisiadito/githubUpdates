FROM golang:latest
LABEL maintainer "LamdaLamdaLamda, Lisiadito"
LABEL version="0.1"
LABEL name="githubUpdates"

RUN command -v go
#RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN apt update
RUN apt-get install go-dep
RUN go get github.com/Lisiadito/githubUpdates
WORKDIR /go/src/github.com/Lisiadito/githubUpdates
RUN dep ensure
RUN go build main.go
#ENTRYPOINT bash
CMD ["./main"]