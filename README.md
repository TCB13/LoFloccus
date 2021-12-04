# LoFloccus
Sync Floccus to a Local Folder!

**LoFloccus** is a small companion app for Floccus that empowers you to sync your bookmarks with whatever service or tool you would like to!

![LoFloccus](https://cdn.iklive.eu/tcb13/2021/lofloccus-1-2-0.png.png) 

## Download LoFloccus

- **Windows**: https://github.com/TCB13/LoFloccus/releases/download/1.2.0/LoFloccus-Win.zip
- **macOS**: https://github.com/TCB13/LoFloccus/releases/download/1.1.2/LoFloccus-macOS.zip

## Why and How

**Floccus** (https://github.com/marcelklehr/floccus) is a great browser extensions that allows you to sync your browser bookmarks with your selfhosted server (e.g. Nextcloud or a WebDAV server).

Due to browser restrictions, you can't store your browser bookmarks in a local file and then sync it with Dropbox, Syncthing, rsync etc. **LoFloccus** aims to make this possible by introducing a small companion app that is a secure, self-contained WebDAV server.

This tool was designed to:
1) Accept WebDAV connections from the Floccus;
2) Restrict Floccus access to a single directory and read/write access limited to `*.xbel` bookmarks files;
3) Generate a random port, username and password for each setup;
4) Store your XBEL bookmarks location and other settings across sessions;
5) Minimize to Windows tray / macOS top menu bar.

Enjoy the best of Floccus and combine it with favourite sync tool!

## Compile from Source
LoFloccus was developed with Qt Creator 4.11.1, Qt 5.14.1 and Golang 1.14.2.

1. Install Qt
2. Install Golang
3. (Optionally) compile libs/libLoFloccusDav.go as described on the file
4. Open the project in Qt Creator and run it.

### Deploy for Windows
1. `cd LoFloccus`
2. `windres.exe .\Windows.rc -o Windows.syso`
3. (Optionally) `cd libs; go build -buildmode c-shared -o libLoFloccusDavWin64.a libLoFloccusDav.go; cd ..`
4. `qtdeploy build desktop .`
5. Run `deploy/windows/LoFloccus.exe`!

### Deploy for macOS
1. `cd LoFloccus`
2. (Optionally) Install imagemagick `brew update && brew install imagemagick`
3. (Optionally) `generate-icns-from-svg.sh`
4. (Optionally) `cd libs; go build -buildmode c-shared -o libLoFloccusDavDarwin.a libLoFloccusDav.go; cd ..`
5. `qtdeploy build desktop .`
6. Run `deploy/darwin/LoFloccus.app`
