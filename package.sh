#!/bin/bash
cd ./release || exit
rm -rf *.zip
list=$(ls frps-panel-*)
echo "$list"
for binFile in $list
  do
    tmpFile=frps-panel
    if echo "$binFile" | grep  -q -E "\.exe";then
      tmpFile=frps-panel.exe
    fi
    cp "$binFile" "$tmpFile"
    zip -r "$binFile".zip "$tmpFile" frps-panel.ini assets -x "*.git*" "*.idea*" "*.DS_Store" "*.contentFlavour"
    rm -rf "$binFile" "$tmpFile"
  done
  rm -rf frps-panel.ini assets
