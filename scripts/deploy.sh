#!/bin/bash

# Date: 10/10/2017
# Author(s): Spicer Matthews (spicer@options.cafe)
# Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
#
# Deploy the entire app (frontend and backend). We do some compiling locally and then deploy.

# cd to backend
cd ../backend

# Build backend
echo "########## Building app.skyclerk.com ##########"
env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o builds/app.skyclerk.com
# upx builds/app.skyclerk.com
cd ../scripts

# Build frontend
cd ../frontend
echo "########## Building Frontend ##########"
npm install
npm run ng build --prod
cd ../scripts

# Build centcom
cd ../centcom
echo "########## Building Centcom ##########"
npm install
cd src
export NODE_ENV=production
npx tailwindcss build tailwind.css -o styles.css
cd ..
npm run ng build --prod --base-href="/centcom/"
cd ../scripts

# Deploy to backend with Ansible
cd ../ansible
ansible-playbook app.yml
cd ../scripts
