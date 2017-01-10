#!/usr/bin/env bash
set -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

if [ -z "$1" ]
then
  echo "No version specified. Aborting."
  exit 1
fi

VERSION=$1

[[ -n "$2" ]] \
  && DOCKER_USER="$2" \
  || DOCKER_USER="magneticio"

echo "pushing docker images..."
docker push ${DOCKER_USER}/sava:1.0.${VERSION}
docker push ${DOCKER_USER}/sava:1.1.${VERSION}
docker push ${DOCKER_USER}/sava-frontend:1.2.${VERSION}
docker push ${DOCKER_USER}/sava-backend1:1.2.${VERSION}
docker push ${DOCKER_USER}/sava-backend2:1.2.${VERSION}
docker push ${DOCKER_USER}/sava-frontend:1.3.${VERSION}
docker push ${DOCKER_USER}/sava-backend:1.3.${VERSION}
