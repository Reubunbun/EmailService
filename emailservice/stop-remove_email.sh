#|/bin/bash
read -p "Enter name: " name

docker stop "${name}userscontainer"
docker rm "${name}userscontainer"
docker rmi "${name}usersimage"

docker stop "${name}msacontainer"
docker rm "${name}msacontainer"
docker rmi "${name}msaimage"

docker stop "${name}mtacontainer"
docker rm "${name}mtacontainer"
docker rmi "${name}mtaimage"