#!/bin/bash
trap "echo 'Himawari is Stopping...'; kill \`pgrep -f Himawari\`" 1 2 3 15

cd ./api/
nohup ./Himawari &

cd ../webui/
npm run start