#! /bin/bash

if [ $# -lt 1 ]; then
   echo "usage: /bin/bash build.sh {APP_PATH} "
   exit 1;
fi

echo "开始构建..."

APP_PATH=$1

if [ x"${GIT_COMMIT_ID}" == "x" ];then
	GIT_COMMIT_ID=`git log -1 --abbrev=8 --pretty=format:%h`
fi
if [ x"${GIT_COMMIT_DATE}" == "x" ];then
	GIT_COMMIT_DATE=`date "+%Y%m%d%H%M%S"`
fi

echo "PATH=${APP_PATH} GIT_COMMIT_ID=${GIT_COMMIT_ID} GIT_COMMIT_DATE=${GIT_COMMIT_DATE}"

APP_VERSION="$(cat VERSION)"
cd ${APP_PATH}
npm install
npm run build
cd -

TAG_NAME="${APP_VERSION}-${GIT_COMMIT_DATE}-${GIT_COMMIT_ID}"
echo TAG_NAME=${TAG_NAME}
echo -n "$TAG_NAME" > TAG_NAME
