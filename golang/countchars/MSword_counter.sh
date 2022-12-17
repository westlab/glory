#!/bin/bash

set -ue

if [[ ! -e $1 ]]; then
    echo $1 ' : does not exist'
    exit 1
fi

SED=""
if [ `uname` = "Darwin" ]; then # for mac
    SED="gsed"
elif [ `uname` = "Linux" ]; then # for linux
    SED="sed"
fi

file_mes=`file $1`

if [[ $file_mes =~ 'Composite Document File V2 Document' ]]; then # doc file
    characters=`file $1 | $SED -r 's/[^"-~\ ]/@/g' | $SED "s/,/\n/g" |grep -E 'Number of Characters' | $SED 's/[^0-9]//g'`
    echo $characters
elif [[ $file_mes =~ 'Microsoft Word 2007+' ]]; then # docx file
    # 空白を含まない文字数
    characters=`unzip -p $1 word/document.xml | xmllint --format - | grep -E '<w:t["-~\ ]*?>'| $SED -r 's/<\/?w:t["-=\?-~\ ]*>//g'| tr -d '\r' | tr -d '\n' | $SED -r 's/\ {2,}//g' | $SED -r 's/&(gt|lt);/@/g'| tr -d ' ' |tr -d ' 　' | wc -m | tr -d ' '`
    echo $characters
else
    echo $1 ': not MS word file'
fi

