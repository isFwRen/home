#!/bin/sh
BUILD_ID=DONTKILLME

echo "name:$1 system:$2"
echo "开始替换package.json"
sed "s/\[Name\]/$1/g" package.json.tmp > package.json
echo "替换完成package.json"
echo "开始替换config.js"
sed "s/\[Name\]/$1/g" config.js.tmp > config.js
echo "替换完成config.js"
echo "开始打包"
if [ "$2" = "32" ]; then
	echo "32位系统"
 	#npm run package -- --arch="ia32"
	npm run package -- --platform=win32 --arch="ia32"
fi

if [ "$2" = "64" ]; then
	echo "64位系统"
	npm run package -- --platform=win32
fi
