#!/bin/bash
cd ./api/
make setup
make build

cd ../webui/
npm install
npm build