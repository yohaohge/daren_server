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

cp -f nohup.out nohup.out.1
cat /dev/null > nohup.out

app=$curDir'/LittleVideo'
echo $app
nohup ${app} &
