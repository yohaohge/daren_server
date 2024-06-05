#!/bin/bash

tag=$1
if [ "$tag" == "" ]; then 
tag="mtv"
fi

# 阿里云镜像仓库
docker login --username=549505007@qq.com registry.cn-hangzhou.aliyuncs.com

docker build -t svr_${tag} . --network=host

docker tag svr_${tag} registry.cn-hangzhou.aliyuncs.com/jasonyuan18/server:${tag}
docker push registry.cn-hangzhou.aliyuncs.com/jasonyuan18/server:${tag}