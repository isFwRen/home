#!/bin/sh
cd /home/stronger/app
echo $1
datetime=`date +%w`
echo $datetime
cp $1 $1_$datetime


