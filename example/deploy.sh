#!/bin/sh -eux

cd `dirname $0`

go mod vendor

set +x

if [ ! -f ./service_account.json ]; then
  echo $GCLOUD_KEY | base64 --decode > service_account.json
  gcloud auth activate-service-account --key-file ./service_account.json
fi

gcloud builds submit --tag gcr.io/${APPLICATION}/fosite-datastore-storage
gcloud beta run deploy fosite-datastore-storage \
  --allow-unauthenticated \
  --image gcr.io/${APPLICATION}/fosite-datastore-storage \
  --region=us-central1 \
  --set-env-vars=DATASTORE_PROJECT_ID=${APPLICATION},BASE_URL=${BASE_URL}
