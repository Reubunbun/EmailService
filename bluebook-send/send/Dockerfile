FROM golang:latest
MAINTAINER 114765
RUN mkdir /email
ADD . /email/
WORKDIR /email
RUN go get github.com/gorilla/mux
EXPOSE 8888
CMD [ "go", "run", "/email/send.go" ] 
