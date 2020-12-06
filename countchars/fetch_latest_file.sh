#!/bin/bash

# 最新のdocxファイルの絶対パス、最終更新日時(RFC3339Nano)、文字数を表示する
# statの表示形式がmacとlinux(centos)で異なるため,linuxのみで使用可能

set -u

SED=""
if [ `uname` = "Darwin" ]; then # for mac
    SED="gsed"
    echo "this script does not support for macOS"
    exit 1
elif [ `uname` = "Linux" ]; then # for linux
    SED="sed"
fi

docxArray=($(ls -t $1/*.docx 2> /dev/null))

if [ "$?" -ne 0 ]; then
  echo "no docx file"
  exit 0
fi

filePath=`echo $docxArray`
updateTime=`stat $filePath | grep Modify | $SED -r 's/^Modify: ([0-9]{4}-[0-9]{2}-[0-9]{2}).*([0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{9}).*$/\1T\2+09:00/'`
count=`./MSword_counter.sh $filePath`


echo -e "$filePath\n$updateTime\n$count"
