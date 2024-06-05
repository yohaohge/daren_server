#!/bin/bash
dir=${PWD##*/}
grep "wespy-http-go/" -rl ${PWD} | grep -v setup.sh | xargs sed -i "" "s/wespy-http-go/${dir}/g"
cp ./config/env.local.toml ./config/env.toml

