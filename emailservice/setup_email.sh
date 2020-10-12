#|/bin/bash
read -p "Enter name: " name
read -p "Enter email server ip: " ip
read -p "Enter users port: " userPort
read -p "Enter MSA port: " msaPort
read -p "Enter MTA port: " mtaPort

cd users

cat << EOF > Dockerfile
FROM golang:latest
MAINTAINER 114765
RUN mkdir /email
ADD . /email/
WORKDIR /email
RUN go get github.com/gorilla/mux
EXPOSE 8888
CMD [ "go", "run", "/email/users.go", "${name}" ] 
EOF

docker build --tag "${name}usersimage" .
docker run --name "${name}userscontainer" --net emailnet --ip "${ip}1" --detach --publish "${userPort}:8888" --security-opt apparmor=unconfined "${name}usersimage"

cd ../MSA

cat << EOF > Dockerfile
FROM golang:latest
MAINTAINER 114765
RUN mkdir /email
ADD . /email/
WORKDIR /email
RUN go get github.com/gorilla/mux
EXPOSE 8888
CMD [ "go", "run", "/email/MSA.go", "${ip}" ] 
EOF

docker build --tag "${name}msaimage" .
docker run --name "${name}msacontainer" --net emailnet --ip "${ip}2" --detach --publish "${msaPort}:8888" --security-opt apparmor=unconfined "${name}msaimage"

cd ../MTA

cat << EOF > Dockerfile
FROM golang:latest
MAINTAINER 114765
RUN mkdir /email
ADD . /email/
WORKDIR /email
RUN go get github.com/gorilla/mux
EXPOSE 8888
CMD [ "go", "run", "/email/MTA.go", "${ip}" ] 
EOF

docker build --tag "${name}mtaimage" .
docker run --name "${name}mtacontainer" --net emailnet --ip "${ip}3" --detach --publish "${mtaPort}:8888" --security-opt apparmor=unconfined "${name}mtaimage"

curl -X POST -d "${ip}" localhost:3000/bluebook/${name}