#!/bin/bash
version=$(cat cmd/frps-panel/cmd.go | grep 'const version' | egrep -o '[0-9.]+')
cd ./release || exit
rm -rf *.zip
list=$(ls frps-panel-*)
echo "$list"
for binFile in $list
  do
    tmpFile=frps-panel
    newBinFile=$binFile
    if echo "$binFile" | grep  -q -E "\.exe";then
      tmpFile=frps-panel.exe
      newBinFile=${newBinFile%%.exe*}
    fi
    cp "$binFile" "$tmpFile"
    zip -r "$newBinFile-$version".zip "$tmpFile" frps-panel.ini assets -x "*.git*" "*.idea*" "*.DS_Store" "*.contentFlavour"
    rm -rf "$binFile" "$tmpFile"
  done
  rm -rf frps-panel.ini assets
