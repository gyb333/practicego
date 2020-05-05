D:
cd D:\Kettle\data-integration
title job
#多个参数需要 -param:"ETLId=1" -param:"args=test"
kitchen.bat /norep -file=C:\Users\zhongduzhi\Desktop\KettleStudy\job.kjb /lever:basic >D:/kettle/data-integration/logs/job.log -param:"ETLId=2" -param:"args=test"
