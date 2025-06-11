#!/bin/bash

# Uses the defualt creds in ~/.aws/credentials

export AWS_ACCESS_KEY_ID=$(op item get "AWS Access Key (cloudmanic)" --fields "access key id" --reveal)
export AWS_SECRET_ACCESS_KEY=$(op item get "AWS Access Key (cloudmanic)" --fields "secret access key" --reveal)
rm ../backend/cache/sqlite/skyclerk-prod.sqlite*
litestream restore -o ../backend/cache/sqlite/skyclerk-prod.sqlite s3://app.skyclerk.com-db