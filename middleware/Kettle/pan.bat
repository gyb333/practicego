D:
cd D:\Kettle\data-integration
title pan
#多个参数需要 -param:"ETLId=1" -param:"args=test"
pan.bat /norep -file=C:\Users\zhongduzhi\Desktop\KettleStudy\pan.ktr /lever:basic >D:/kettle/data-integration/logs/pan.log -param:"ETLId=1" -param:"args=test"
