D:
cd D:\Kettle\data-integration
title demojob
::  -param:"IsStartPageNo=true" -param:"StartPageNo=4"
::kitchen.bat /norep -file=C:\Users\zhongduzhi\Desktop\KettleStudy\��ҳdemo\��ҳdemo.kjb /lever:basic >D:/kettle/data-integration/logs/job.log
kitchen.bat /norep -file=C:\Users\zhongduzhi\Desktop\KettleStudy\��ҳdemo\��ҳdemo.kjb /lever:basic>>D:/kettle/data-integration/logs/job.log -param:"IsStartPageNo=false" -param:"StartPageNo=3"