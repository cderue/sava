#!/usr/bin/env bash

dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

reset=`tput sgr0`
green=`tput setaf 2`

version="1.0.0"

cd ${dir}
rm -Rf ${dir}/target 2> /dev/null && mkdir -p ${dir}/target
find . -name 'Dockerfile' | cpio -pdm ${dir}/target 2> /dev/null

export GOOS='linux'
export GOARCH='amd64'
echo "${green}building binary for ${GOOS}:${GOARCH}${reset}"
CGO_ENABLED=0 go build -a -installsuffix cgo
chmod +x sava_runner
mv sava_runner target

cd ${dir}/target
echo "${green}building docker image: magneticio/sava_runner:${version} ${reset}"
docker build -t magneticio/sava_runner:${version} .

echo "${green}done.${reset}"
