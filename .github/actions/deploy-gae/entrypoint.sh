#!/bin/bash

set -eu

echo ${INPUT_CREDENTIALS} > credentials.json
gcloud auth activate-service-account --key-file credentials.json
gcloud app deploy --project ${INPUT_PROJECT_ID} --quiet
