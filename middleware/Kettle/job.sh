#!/bin/bash
 
cd /home/kettle/kettle_job

#echo $#
currTime=$(date "+%Y-%m-%d,%H:%M:%S")
for dir in $(find . -type f -name *_$1_*.kjb | sort -n)
do
  #echo $dir
  echo ${dir##*/}
  /home/kettle/data-integration/kitchen.sh -file $dir /lever:basic>>./logs/${dir##*/}${currTime}.log
done


