#!/bin/bash

set -eu

echo ${INPUT_CREDENTIALS} > credentials.json
gcloud auth activate-service-account --key-file credentials.json
gcloud config set project ${INPUT_PROJECT_ID}

gcloud app deploy --quiet
old_versions=$(gcloud app versions list --filter=traffic_split=0.00 --format='value(id)' | tr '\n' ' ')
gcloud app versions delete --quiet $old_versions
