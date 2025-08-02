#! /bin/bash

if [ $# -lt 4 ]; then
   echo "usage: /bin/bash build.sh {APP_PATH} {GOARCH} {GOOS} {PACKAGE}"
   exit 1;
fi

echo "开始构建..."

APP_PATH=$1
GOARCH=$2
GOOS=$3
PACKAGE=$4

if [ x"${GIT_COMMIT_ID}" == "x" ];then
	GIT_COMMIT_ID=`git log -1 --abbrev=8 --pretty=format:%h`
fi
if [ x"${GIT_COMMIT_DATE}" == "x" ];then
	GIT_COMMIT_DATE=`date "+%Y%m%d%H%M%S"`
fi

echo "PATH=${APP_PATH} OS=${GOOS} ARCH=${GOARCH} PACKAGE=${PACKAGE} GIT_COMMIT_ID=${GIT_COMMIT_ID} GIT_COMMIT_DATE=${GIT_COMMIT_DATE}"

export GOOS=${GOOS}
export GOARCH=${GOARCH}
export BINARY=${PACKAGE}

APP_VERSION="$(cat VERSION)"
cd ${APP_PATH}
make
archiveFile="${PACKAGE%.exe}-${GOOS}-${GOARCH}-${APP_VERSION}.tar"
tar zcf ${archiveFile}.gz -C bin ${BINARY}
cd -

PACKAGE_FILE_NAME="${APP_PATH}/${archiveFile}.gz"
TAG_NAME="${APP_VERSION}-${GIT_COMMIT_DATE}-${GIT_COMMIT_ID}"
echo PACKAGE_FILE_NAME=${PACKAGE_FILE_NAME}
echo TAG_NAME=${TAG_NAME}
echo -n "$PACKAGE_FILE_NAME" > PACKAGE_NAME
echo -n "$TAG_NAME" > TAG_NAME
