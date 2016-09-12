#!/bin/bash 

yid=$(cat work/yipid)
d=$(date)

ps cax | grep $yid > /dev/null
if [ $? -eq 0 ]; then
  echo "Process is running."
else
  bash start.sh
  echo $d" :: restarting">>work/restart.log
fi
