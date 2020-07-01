#!/bin/bash
 
cd /home/kettle/kettle_job
for dir in $(find . -type f -name *.kjb | sort -n)
do
  echo $dir
  #echo ${dir##*/} 
  /home/kettle/data-integration/kitchen.sh -file $dir /lever:basic>>./logs/${dir##*/}.log
done

