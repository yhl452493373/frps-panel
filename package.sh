#!/bin/bash
cd ./release || exit
rm -rf *.zip
list=$(ls frps-multiuser-*)
echo "$list"
for binFile in $list
  do
    cp "$binFile" frps-multiuser
    zip -r "$binFile".zip frps-multiuser frps-multiuser.ini assets -x "*.git*" "*.idea*" "*.DS_Store" "*.contentFlavour"
    rm -rf "$binFile" frps-multiuser
  done
  rm -rf frps-multiuser.ini assets
