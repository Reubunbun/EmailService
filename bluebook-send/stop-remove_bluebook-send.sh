#|/bin/bash

docker stop bluebookcontainer
docker rm bluebookcontainer
docker rmi bluebookimage

docker stop sendcontainer
docker rm sendcontainer
docker rmi sendimage