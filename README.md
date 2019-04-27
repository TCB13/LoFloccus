# LoFloccus
Sync Floccus to a Local Folder!

**LoFloccus** is a small companion app for Floccus that empowers you to sync your bookmarks with whatever service or tool you would like to!

![LoFloccus](https://cdn.iklive.eu/tcb13/2019/lofloccus-1-1-0.png)

## Download LoFloccus

- **Windows**: https://github.com/TCB13/LoFloccus/releases/download/1.1.1/LoFloccus-1-1-1.zip
- **macOS**: Available Soon!

## Why and How

**Floccus** (https://github.com/marcelklehr/floccus) is a great browser extensions that allows you to sync your browser bookmarks with your selfhosted server (e.g. Nextcloud or a WebDAV server).

Unfortunately (and due to browser restrictions) you can't store your browser bookmarks in a local file and then sync it with Dropbox, Syncthing, rsync etc. **LoFloccus** aims to make this possible by introducing a small companion app that is a secure, self-contained WebDAV server.

This tool was designed to:
1) Accept WebDAV connections from the localhost;
2) Restrict Floccus access to a single directory and read/write access limited to `*.xbel` bookmarks files;
3) Generate a random port, username and password for each setup;
4) Store your XBEL bookmarks location and other settings across sessions;
5) Minimize to Windows tray / macOS top menu bar.

Enjoy the best of Floccus and combine it with favourite sync tool!

## Compile from Source

LoFloccus was built using Golang and Qt.

#### Windows
**Develop**:
1. Install Qt
2. Install https://github.com/therecipe/qt
3. Add user variable `CGO_LDFLAGS`: `C:\Users\Public\env_windows_amd64_Tools\mingw730_64\x86_64-w64-mingw32\lib\libmsvcrt.a`
4. Add to user `Path` variable: `C:\Users\Public\env_windows_amd64\5.12.0\mingw73_64\bin` and `C:\Users\Public\env_windows_amd64_Tools\mingw730_64\bin`
5. `go run LoFloccus.go`!

**Deploy**:
1. `cd LoFloccus`
2. `windres.exe .\Windows.rc -o Windows.syso`
3. `qtdeploy build desktop .`
4. Run `deploy/windows/LoFloccus.exe`!
