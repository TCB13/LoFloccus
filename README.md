# LoFloccus
Sync Floccus with a local XBEL file.

**LoFloccus** is a small companion app for Floccus that empowers you to sync your bookmarks with whatever service or tool you would like to!

![LoFloccus](https://cdn.iklive.eu/tcb13/2019/lofloccus.png)

Floccus (https://github.com/marcelklehr/floccus) is a great browser extensions that allows you to sync your browser bookmarks with your selfhosted server (e.g. Nextcloud or a WebDAV server).

Unfortunately (and due to browser restrictions) you can't simply store your browser bookmarks in a local file and then sync them with Dropbox, Syncthing, rsync or other methods. **LoFloccus** aims to make this possible by introducing a companion app that is a secure, portable and self-contained local WebDAV server with a very small footprint.

This tool was designed to:
1) Accept WebDAV connections from the localhost;
2) Generate a random port, username and password for each setup;
3) Store your XBEL location and other settings across sessions;
4) Minimize to Windows tray / macOS top bar.

Enjoy the best of Floccus and combine it with favourite sync tool!
