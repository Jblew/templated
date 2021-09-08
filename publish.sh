#!/usr/bin/env bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${DIR}"
set -e

BASE_TAG="jedrzejlewandowski/templated"
VERSION="1.2.0"

git tag "v${VERSION}"
git push origin "v${VERSION}"

docker build -t "${BASE_TAG}" .
docker tag "${BASE_TAG}" "${BASE_TAG}:latest"
docker tag "${BASE_TAG}" "${BASE_TAG}:${VERSION}"
docker tag "${BASE_TAG}" "${BASE_TAG}:${VERSION}-alpine"
docker push "${BASE_TAG}:latest"
docker push "${BASE_TAG}:${VERSION}"
docker push "${BASE_TAG}:${VERSION}-alpine"