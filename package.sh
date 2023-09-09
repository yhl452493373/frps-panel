#!/bin/bash
cd ./release || exit
rm -rf *.zip
list=$(ls frps-panel-*)
echo "$list"
for binFile in $list
  do
    cp "$binFile" frps-panel
    zip -r "$binFile".zip frps-panel frps-panel.ini assets -x "*.git*" "*.idea*" "*.DS_Store" "*.contentFlavour"
    rm -rf "$binFile" frps-panel
  done
  rm -rf frps-panel.ini assets
