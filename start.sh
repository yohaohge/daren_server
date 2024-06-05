#! /bin/bash

set -e

cur_dir=`pwd`
dst_dir='/usr/local/mtv_svr/'

pid=`ps -ef | grep "/usr/local/mtv_svr/mtv_svr" | grep -v grep | awk '{print $2}'`
if ! [ -z "$pid" ]; then
    echo "close old. pid: "${pid} 
    kill $pid
fi

cp ./mtv_svr ${dst_dir}
cp ./config/env.online.toml ${dst_dir}/config/env.toml
cp ./config/config.json ${dst_dir}/config/config.json

cd $dst_dir
nohup ./mtv_svr > nohup.out 2>&1 &

