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

pidof $curDir'/LittleVideo' | xargs kill
