#!/bin/sh -x

set -e

SIZES="
16,16x16
32,16x16@2x
32,32x32
64,32x32@2x
128,128x128
256,128x128@2x
256,256x256
512,256x256@2x
512,512x512
1024,512x512@2x
"

FILES=$PWD/icon.svg
for SVG in $FILES
do
  echo "Processing $f file..."
  # take action on each file. $f store current file name
  BASE=$(basename "$SVG" | sed 's/\.[^\.]*$//')
    ICONSET="$BASE.iconset"
    mkdir -p "$PWD/icons/$ICONSET"
    for PARAMS in $SIZES; do
        SIZE=$(echo $PARAMS | cut -d, -f1)
        LABEL=$(echo $PARAMS | cut -d, -f2)
        svg2png -w $SIZE -h $SIZE "$SVG" "$PWD/icons/$ICONSET"/icon_$LABEL.png || true
    done

    iconutil -c icns "$PWD/icons/$ICONSET" || true
    rm -rf "$ICONSET"
done