#!/usr/bin/env bash
# gcloud config set project {PROJECTID}
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
HOSTNAME=asia.gcr.io
PROJECTID=${LINE_PROJECT}
IMAGE=cvbot
TAG=`git rev-parse --short=7 HEAD`

cd "${DIR}/.."

docker build -t ${HOSTNAME}/${PROJECTID}/${IMAGE}:${TAG} -f deploy/Dockerfile .
if [[ ${?} != 0 ]]; then
	# build error
    exit $?
fi
docker push ${HOSTNAME}/${PROJECTID}/${IMAGE}:${TAG}
gcloud beta run deploy ${IMAGE} --project ${PROJECTID} --image ${HOSTNAME}/${PROJECTID}/${IMAGE}:${TAG} --region us-central1 --platform managed
