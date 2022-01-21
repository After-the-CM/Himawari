#!/bin/bash
trap "kill \`pgrep -f Himawari\`" 1 2 3 15

cd ./api/
make setup
nohup make run &

cd ../webui/
npm run dev