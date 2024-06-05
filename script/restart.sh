#!/bin/bash

getScriptDir() 
{
    oldwd=`pwd`
    rw=`dirname $0`
    cd $rw
    sw=`pwd`
    cd $oldwd
    echo $sw 
}

curDir=`getScriptDir`

pid=`pidof $curDir'/LittleVideo'`
cp -f nohup.out nohup.out.1
cat /dev/null > nohup.out
if [ "$pid" == "" ]; then
    nohup ./LittleVideo &
else
    kill -SIGUSR2 ${pid}
fi
