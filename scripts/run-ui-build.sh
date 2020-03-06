#!/bin/bash
cd ui
npm install
npm run build
rm -rf ../server/app/dist
mv dist ../server/app
