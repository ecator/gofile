#!/bin/bash

# 编译打包一条龙

tag=`git tag -l | tail -n1`
if [ -z "$tag" ]
then
  echo "no tags" >&2
  exit 1
fi

platform=`uname -p`
if [ "$platform" == "unknown" ]
then
  platform=`uname -m`
fi
kernel=`uname -s`

target_dir="GoFile-${tag}-${kernel}-${platform}"

echo "making dist/$target_dir"
[ ! -d dist ] && mkdir dist
[ -e "dist/$target_dir" ] && rm -rf "dist/$target_dir"

go build -o "dist/$target_dir/gofile" -ldflags "-w" main.go
cd web && npm run build && cd ..
mkdir "dist/$target_dir/web" && cp -a web/dist "dist/$target_dir/web/"

cd dist
tar -cJf "${target_dir}.tar.xz" $target_dir

if [ $? -eq 0 ]
then
  echo "dist/$target_dir.tar.xz ok"
fi