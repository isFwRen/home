#!/bin/sh
BUILD_ID=DONTKILLME

cd /home/stronger/service/
echo "端口  二进制文件  项目编码  判断关进程:1还是开进程:2  screen名字  下载后缀为这些的单 下载路径"
echo $1 $2 $3 $4 $5 $6 $7

find -name "$2_*" -mtime +1 -delete
# datetime=`date +%Y%m%d\-%H%M%S`
datetime=``

if [ "$4" = "1" ]; then
    echo "开始关闭进程"
    info=`ps aux | grep $5 | grep SCREEN  |awk '/'$1'/'`
    info1=`netstat -lnpt | grep $1 | awk '{print $7}' | awk -F '/' '{print $1}'`
    echo $info  
    echo $info1
    if [ -n "$info" ]; then
        screen -S $5 -X quit
    fi
    if [ -n "$info1" ]; then
        kill -9 $info1
    fi
    echo "完成关闭进程"

    echo "开始拷贝二进制文件"
    # sshpass -p $1 scp $2 ./bin/$3_$datetime
    cp ./bin/$2 ./bin/$2_$datetime
    echo "完成拷贝二进制文件"
fi
if [ "$4" = "2" ]; then
	echo "开始开启进程"
    # chmod 777 $2
    echo screen -dmS $5 ./deploy/start_screen.sh $1 $2 $3 $6 $7
    screen -dmS $5 ./deploy/start_screen.sh $1 $2 $3 $6 $7
    # screen -dmS $2 ./deploy/$2.sh
    echo "完成开启进程"
fi


