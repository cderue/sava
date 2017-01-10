#!/usr/bin/env bash
set -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

if [ -z "$1" ]
then
  echo "No version specified. Aborting."
  exit 1
fi

VERSION=$1

echo "cleaning..."

rm -Rf ${DIR}/target 2> /dev/null
mkdir ${DIR}/target

echo "installing dependencies..."
go get -v github.com/gin-gonic/gin

echo "building sava executables..."

for SERVICE in sava_1.0 sava_1.1 sava_frontend_1.2 sava_backend_1.2 sava_frontend_1.3
  do
    echo "building: ${SERVICE}"
    cd ${DIR}/../${SERVICE}
    # build the app for linux/i386 an build the docker container
    GOOS=linux GOARCH=386 go build
  done


function createDockerFileFor {
  sed -e 's/_SERVICE_/'${1}'/g' ${DIR}/Dockerfile > ${DIR}/target/$1/Dockerfile
}

function buildDockerImageFor {
  docker build -t magneticio/${1} ${DIR}/target/${2}
}

echo "prepping Docker files and dependencies"

for SERVICE in sava_1.0 sava_1.1 sava_frontend_1.2 sava_backend_1.2 sava_frontend_1.3
  do
    echo "building Docker image: ${SERVICE}"
    mkdir ${DIR}/target/${SERVICE}
    if [ -d ${DIR}/../${SERVICE}/public ]; then
      cp -r ${DIR}/../${SERVICE}/public ${DIR}/target/${SERVICE}/public
    fi
    if [ -s ${DIR}/../${SERVICE}/index.tmpl ]; then
      cp -r ${DIR}/../${SERVICE}/index.tmpl ${DIR}/target/${SERVICE}/public/index.tmpl
    fi
    mv ${DIR}/../${SERVICE}/${SERVICE} ${DIR}/target/${SERVICE}/
    createDockerFileFor ${SERVICE}
  done

echo "building docker image: sava_1.0".${VERSION}
buildDockerImageFor sava:1.0.${VERSION} sava_1.0

echo "building docker image: sava_1.1".${VERSION}
buildDockerImageFor sava:1.1.${VERSION} sava_1.1

echo "building docker image: sava_frontend_1.2".${VERSION}
buildDockerImageFor sava-frontend:1.2.${VERSION} sava_frontend_1.2

echo "building docker image: sava-backend1:1.2".${VERSION}
buildDockerImageFor sava-backend1:1.2.${VERSION} sava_backend_1.2

echo "building docker image: sava-backend2:1.2".${VERSION}
buildDockerImageFor sava-backend2:1.2.${VERSION} sava_backend_1.2

echo "building docker image: sava_frontend_1.3".${VERSION}
buildDockerImageFor sava-frontend:1.3.${VERSION} sava_frontend_1.3

echo "building docker image: sava-backend:1.3".${VERSION}
buildDockerImageFor sava-backend:1.3.${VERSION} sava_backend_1.2


# echo "pushing docker images..."
# docker push magneticio/sava:1.0.${VERSION}
# docker push magneticio/sava:1.1.${VERSION}
# docker push magneticio/sava-frontend:1.2.${VERSION}
# docker push magneticio/sava-backend1:1.2.${VERSION}
# docker push magneticio/sava-backend2:1.2.${VERSION}
# docker push magneticio/sava-frontend:1.3.${VERSION}
# docker push magneticio/sava-backend:1.3.${VERSION}
