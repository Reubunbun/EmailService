#|/bin/bash

docker network create --subnet 192.168.1.0/24 emailnet

cd bluebook

cat << EOF > Dockerfile
FROM golang:latest
MAINTAINER 114765
RUN mkdir /email
ADD . /email/
WORKDIR /email
RUN go get github.com/gorilla/mux
EXPOSE 8888
CMD [ "go", "run", "/email/bluebook.go" ] 
EOF

docker build --tag bluebookimage .
docker run --name bluebookcontainer --net emailnet --ip 192.168.1.2 --detach --publish 3000:8888 --security-opt apparmor=unconfined bluebookimage

cd ../send

cat << EOF > Dockerfile
FROM golang:latest
MAINTAINER 114765
RUN mkdir /email
ADD . /email/
WORKDIR /email
RUN go get github.com/gorilla/mux
EXPOSE 8888
CMD [ "go", "run", "/email/send.go" ] 
EOF

docker build --tag sendimage .
docker run --name sendcontainer --net emailnet --ip 192.168.1.3 --detach --publish 3001:8888 --security-opt apparmor=unconfined sendimage