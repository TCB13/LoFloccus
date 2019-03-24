# pip install wsgidav cheroot pystray
# Make Windows .exe with: pyinstaller -i icon.ico -F -w .\LoFloccus.py

from tkinter import *
from tkinter import messagebox
from tkinter import filedialog
import logging
import threading
import shelve
import random
import pystray
from pystray import Menu, MenuItem
from PIL import Image, ImageTk
from cheroot import wsgi
from wsgidav.wsgidav_app import WsgiDAVApp

global server
global folder

shelf = shelve.open("LoFloccus-settings")

# Create random password/port if needed
if "port" not in shelf:
    shelf["port"] = random.randint(40000, 65535)
if "passwd" not in shelf:
    shelf["passwd"] = "local-" + str(random.randint(100, 999))
# Display settings
if "startup.tray" not in shelf:
    shelf["startup.tray"] = False
if "startup.server" not in shelf:
    shelf["startup.server"] = False

shelf.sync()


def start_server():
    global server
    global shelf

    config = {
        "host": "127.0.0.1",
        "port": shelf["port"],
        "provider_mapping": {
            "/": shelf["path"],
        },
        "verbose": 1,
        "enable_loggers": [],
        "http_authenticator": {
            "domain_controller": None,
            "accept_basic": False,
            "accept_digest": True,
            "default_to_digest": True

        },
        "simple_dc": {
            "user_mapping":
                {
                    "*":
                        {
                            "floccus": {
                                "password": shelf["passwd"]
                            }
                        }
                }
        }
    }

    print("Starting server...")
    lblStatus.configure(text="ON", foreground="green")
    btnServer.config(state="normal", text="<< Stop Server >>")
    btnChangeDir.config(state="disabled")

    # Start the server
    app = WsgiDAVApp(config)
    server_args = {
        "bind_addr": (config["host"], config["port"]),
        "wsgi_app": app,
    }
    server = wsgi.Server(**server_args)
    logger = logging.getLogger("wsgidav")
    logger.propagate = True
    logger.setLevel(logging.DEBUG)
    server.start()


def stop_server():
    global server

    print("Stopping server...");
    btnServer.config(state="disabled")
    btnChangeDir.config(state="active")

    try:
        server
    except NameError:
        print("Server not running.")
    else:
        print("Stop server.");
        server.stop()
        del server
        print("Server is down!");
        btnServer.config(state="normal", text="<< Start Server >>")
        lblStatus.configure(text="OFF", foreground="red")


def toggle_server():
    global server
    try:
        server
    except NameError:
        if "path" in shelf:
            threading.Thread(target=start_server).start()
        else:
            messagebox.showerror("Error", "You must set a directory where to store your XBEL file.")
            print("No server path defined.")
    else:
        threading.Thread(target=stop_server).start()


def change_dir():
    userInput = filedialog.askdirectory()
    if userInput == "":
        return None

    shelf["path"] = userInput
    shelf.sync()
    lblStorePath.config(text=shelf["path"])


def on_closing():
    global window
    if messagebox.askokcancel("Quit", "Do you want to quit? You won't be able to sync your favourites."):
        stop_server()
        window.destroy()


def setup_tray(icon):
    icon.visible = True


def hide_tray():
    global window
    print("tray hide")
    # hide window
    window.withdraw()
    window.overrideredirect(True)

    # Create tray icon
    icon = pystray.Icon("Floccus Local XBEL")
    icon.icon = Image.open("logo.png")
    icon.visible = True
    icon.menu = Menu(
        MenuItem("Open", lambda: restore_tray(icon, False)),
        MenuItem("Exit", lambda: restore_tray(icon, True)),
    )
    icon.title = "tooltip"
    icon.run(setup_tray)


def restore_tray(icon, close):
    global window
    print("tray restore")
    icon.visible = False
    icon.stop()

    # show app
    window.overrideredirect(False)
    window.update()
    window.deiconify()
    window.attributes("-topmost", True)

    if close:
        on_closing()


def check_toggle(key, value):
    global shelf
    shelf[key] = value.get()
    shelf.sync()
    print(shelf[key])


window = Tk()
window.title("LoFloccus - v1.0")
window.iconbitmap("icon.ico")
window.resizable(False, False)
# Window Size and Positioning
windowWidth = 480
windowHeight = 470
ws = window.winfo_screenwidth()
hs = window.winfo_screenheight()
x = (ws / 2) - (windowWidth / 2)
y = (hs / 2) - (windowHeight / 2)
window.geometry('%dx%d+%d+%d' % (windowWidth, windowHeight, x, y))

btnPadding = (10, 5)
labelPad = (20, 0)

# Title

img = ImageTk.PhotoImage(Image.open("logo.png"))
panel = Label(window, image=img)
panel.grid(column=0, row=0, columnspan=2,
           pady=(20, 0))

Label(window, text="LoFloccus", font=("Arial Bold", 15)).grid(column=0, row=1, columnspan=2,
                                                              pady=(0, 0))
Label(window, text="Sync Floccus with a local XBEL file.", font=("Arial Bold", 10)).grid(column=0, row=2, columnspan=2,
                                                                                         pady=(0, 20))
Label(window, text="Server Status: ").grid(column=0, row=3, padx=labelPad, sticky=W)
lblStatus = Label(window, text="OFF", foreground="red")
lblStatus.grid(column=1, row=3, sticky=W)

Label(window, text="Store XBEL in: ").grid(column=0, row=4, padx=labelPad, sticky=W)

if "path" in shelf:
    lblStorePath = Label(window, text=shelf["path"])
else:
    lblStorePath = Label(window, text="(null)")

lblStorePath.grid(column=1, row=4, sticky=W)

btnChangeDir = Button(window, text="Change Directory", command=change_dir)
btnChangeDir.grid(column=0, row=5, pady=btnPadding, padx=labelPad, sticky=W)
btnServer = Button(window, text="<< Start Server >>", command=toggle_server)
btnServer.grid(column=1, row=5, pady=btnPadding, sticky=W)
btnTray = Button(window, text="Hide to Tray Now", command=hide_tray)
btnTray.grid(column=1, row=5, pady=btnPadding, sticky=E, padx=(0, 15))

Label(window, text="Pick an empty folder. Never store other files in the same folder where your XBEL is.").grid(
    column=0, row=6, columnspan=2, pady=(5, 0), padx=labelPad, sticky=W)

# Notes
Label(window, text="Connect your Floccus to (Add \"XBEL in WebDAV Account\"):", font=("Arial Bold", 10)).grid(column=0,
                                                                                                              row=7,
                                                                                                              columnspan=2,
                                                                                                              pady=(
                                                                                                                  15,
                                                                                                                  5),
                                                                                                              padx=labelPad,
                                                                                                              sticky=W)
Label(window, text="Server: ").grid(column=0, row=8, padx=labelPad, sticky=W)
Label(window, text="http://127.0.0.1:" + str(shelf["port"])).grid(column=1, row=8, sticky=W)
Label(window, text="Username: ").grid(column=0, row=9, padx=labelPad, sticky=W)
Label(window, text="floccus").grid(column=1, row=9, sticky=W)
Label(window, text="Password: ").grid(column=0, row=10, padx=labelPad, sticky=W)
Label(window, text=shelf["passwd"]).grid(column=1, row=10, sticky=W)

# More settings
Label(window, text="App Settings:", font=("Arial Bold", 10)).grid(column=0, row=11, columnspan=2, pady=(15, 5),
                                                                  padx=labelPad, sticky=W)

checkStartTrayVal = BooleanVar(value=shelf["startup.tray"])
checkStartServerVal = BooleanVar(value=shelf["startup.server"])

Checkbutton(window, text="Start minimized to tray/top menu bar", variable=checkStartTrayVal,
            command=lambda: check_toggle("startup.tray", checkStartTrayVal),
            height=1).grid(column=0, row=12, columnspan=2, padx=labelPad, sticky=W)

Checkbutton(window, text="Start server when open", variable=checkStartServerVal,
            command=lambda: check_toggle("startup.server", checkStartServerVal),
            height=1).grid(column=0, row=13, columnspan=2, padx=labelPad, sticky=W)

window.protocol("WM_DELETE_WINDOW", on_closing)

# Autostart settings
if shelf["startup.tray"]:
    hide_tray()

if shelf["startup.server"]:
    toggle_server()

window.mainloop()
