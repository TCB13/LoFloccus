#!/bin/bash

mkdir icon.iconset
convert -background none -size 16x16 -gravity center -extent 16x16     icon.svg icon.iconset/icon_16x16.png
convert -background none -size 32x32 -gravity center -extent 32x32     icon.svg icon.iconset/icon_16x16@2x.png
convert -background none -size 32x32 -gravity center -extent 32x32     icon.svg icon.iconset/icon_32x32.png
convert -background none -size 64x64 -gravity center -extent 64x64     icon.svg icon.iconset/icon_32x32@2x.png
convert -background none -size 64x64 -gravity center -extent 64x64     icon.svg icon.iconset/icon_64x64.png
convert -background none -size 128x128 -gravity center -extent 128x128   icon.svg icon.iconset/icon_64x64@2x.png
convert -background none -size 128x128 -gravity center -extent 128x128   icon.svg icon.iconset/icon_128x128.png
convert -background none -size 256x256 -gravity center -extent 256x256   icon.svg icon.iconset/icon_128x128@2x.png
convert -background none -size 256x256 -gravity center -extent 256x256   icon.svg icon.iconset/icon_256x256.png
convert -background none -size 512x512 -gravity center -extent 512x512   icon.svg icon.iconset/icon_256x256@2x.png
convert -background none -size 512x512 -gravity center -extent 512x512   icon.svg icon.iconset/icon_512x512.png
convert -background none -size 1024x1024 -gravity center -extent 1024x1024 icon.svg icon.iconset/icon_512x512@2x.png
iconutil --convert icns icon.iconset
rm -R icon.iconset

