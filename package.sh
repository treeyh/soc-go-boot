#!/bin/bash

## TODO 待完善


set -x

echo "package start ..."

UNAME=`uname`

# 项目分支
PROJECT_BRANCH=$1
# build版本号
BUILD_NUMBER=$2

#取当前目录
BASE_PATH=`cd "$(dirname "$0")"; pwd`

PROJECT_NAME="soc-go-boot"



# 获取源码最近一次 git commit log，包含 commit sha 值，以及 commit message
GitCommitLog=`git log --pretty=oneline -n 1`
# 将 log 原始字符串中的单引号替换成双引号
GitCommitLog=${GitCommitLog//\'/\"}
# 检查源码在 git commit 基础上，是否有本地修改，且未提交的内容
GitStatus=`git status -s`
# 获取当前时间
BuildTime=`date +'%Y-%m-%dT%H%M%S'`
# 获取 Go 的版本
BuildGoVersion=`go version`

# 将以上变量序列化至 LDFlags 变量中
LDFlags=" \
    -X 'github.com/treeyh/soc-go-boot/app/common/buildinfo.GitCommitLog=${GitCommitLog}' \
    -X 'github.com/treeyh/soc-go-boot/app/common/buildinfo.GitStatus=${GitStatus}' \
    -X 'github.com/treeyh/soc-go-boot/app/common/buildinfo.BuildTime=${BuildTime}' \
    -X 'github.com/treeyh/soc-go-boot/app/common/buildinfo.BuildGoVersion=${BuildGoVersion}' \
"

#发布目录
TARGET_PATH="$BASE_PATH/target"

#备份目录
BAK_PATH="$BASE_PATH/bak"
CUR_TIME=`date +%Y%m%d%H%m%s`
BAK_PATH=$BAK_PATH/$CUR_TIME

# 创建编译目标目录和备份目录
mkdir -p $TARGET_PATH
mkdir -p $BAK_PATH

# 清理发布目录
mv $TARGET_PATH/* $BAK_PATH

# 复制docker相关文件
cp -rf $BASE_PATH/build/docker/*  $TARGET_PATH

# 设置环境变量
export GO111MODULE=on
export GOPROXY=https://goproxy.io

function build_package(){
    local package_path=$BASE_PATH/$1



    # 复制配置文件
    cp -rf $BASE_PATH/config  $package_target_path/
    cp -rf $package_path/config/*  $package_target_path/config/

    if [ -d "$package_path/resources" ]; then
      cp -rf $package_path/resources  $package_target_path/;
    fi

    cp -rf $BASE_PATH/init  $package_target_path/
    cp $BASE_PATH/README.md $package_target_path/
    cp $BASE_PATH/CHANGELOG.md $package_target_path/

    # 复制执行命令
    cp -rf $BASE_PATH/bin  $package_target_path/

    # 修改启动脚本中的应用名
    sed -i "s/#APP_NAME#/${project_name}/g" $package_target_path/bin/single.sh
#    if [[ "$UNAME" == "Darwin" ]]; then
#      sed -ie "s/#APP_NAME#/${project_name}/g" $package_target_path/bin/single.sh
#    else
#      sed -i "s/#APP_NAME#/${project_name}/g" $package_target_path/bin/single.sh
#    fi


    cd $package_path

    # 编译
    if [[ "$POL_ENV" == "local" || "$POL_ENV" == "" ]]; then
        go build -ldflags "$LDFlags" -o $package_target_path/${project_name}-app
    else
        GOOS=linux GOARCH=amd64 go build -ldflags "$LDFlags" -o $package_target_path/${project_name}-app
    fi

    echo '------'

    #清理文件
#    rm $package_path/$1
}


build_package "cmd/main.go"

echo "package end ..."