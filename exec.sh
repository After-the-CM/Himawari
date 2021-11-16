#!/bin/bash
trap "echo 'Himawari is Stopping...'; kill \`pgrep -f Himawari\`; rm -rf ../nohup.out" 1 2 3 15

nohup ./api/Himawari &
cd ./webui/
# 完成版では、npm run startに変更が入ります。
npm run dev 