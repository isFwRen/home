#!/bin/sh
while : 
 do 
 echo "1端口  2二进制文件  3项目编码 4下载后缀为这些的单 5下载路径"
 echo $1 $1 $3 $4 $5
 echo ./bin/$2 -P $1 -n $3 -c ./config_test.yaml -b $4 -p $5
 ./bin/$2 -P $1 -n $3 -c ./config_test.yaml -b $4 -p $5
 done
