#!/usr/bin/env sh
apk update && apk add --no-cache git mercurial openssh
export REPOSITORY_PRIVATE_KEY=$(echo $REPOSITORY_PRIVATE_KEY | base64 -d)
export CGO_ENABLED=0
export GO111MODULE=on
export GOPRIVATE=github.com/Confialink
mkdir -p ~/.ssh && umask 0077 && echo "${REPOSITORY_PRIVATE_KEY}" > ~/.ssh/id_rsa \
&& git config --global url."git@github.com:Confialink".insteadOf https://github.com/Confialink \
&& ssh-keyscan bitbucket.org >> ~/.ssh/known_hosts \
&& ssh-keyscan github.com >> ~/.ssh/known_hosts