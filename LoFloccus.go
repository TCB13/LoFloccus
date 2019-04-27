/*
LoFloccus - Sync Floccus to a Local Folder!
 */

package main

import (
	"context"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/uitools"
	"github.com/therecipe/qt/widgets"
	"golang.org/x/net/webdav"
	"gopkg.in/ini.v1"
	"math/rand"
	"time"

	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type ServerConfig struct {
	address string
	port    int
	dir     string
	user    string
	passwd  string
}

type AppSettings struct {
	startMinimized    bool
	startLaunchServer bool
}

var (
	configFile   = "LoFloccus-Settings.ini"
	cfg *ini.File
	windowWidget *widgets.QWidget
	serverHandle *http.Server
	serverConfig = ServerConfig{
		address: "127.0.0.1",
		port:    0,
		dir:     "",
		user:    "floccus",
		passwd:  "",
	}
	appSettings = AppSettings{
		startMinimized:false,
		startLaunchServer:false,
	}
)

func main() {

	// Load AppConfig from config file
	loadAppConfig()

	// Launch Qt Application
	core.QCoreApplication_SetAttribute(core.Qt__AA_ShareOpenGLContexts, true)

	// Create main window from the .ui file
	widgets.NewQApplication(len(os.Args), os.Args)
	loader := uitools.NewQUiLoader(nil)
	file := core.NewQFile2(":/qml/LoFloccus.ui")
	if ok := file.Open(core.QIODevice__ReadOnly); !ok {
		log.Fatal("Could not open the UI file")
	}
	defer file.Close()
	var icon = gui.NewQIcon5(":/qml/icon.ico")
	windowWidget = loader.Load(file, nil)
	windowWidget.SetFixedSize2(windowWidget.Width(), windowWidget.Height())
	windowWidget.SetWindowIcon(icon)
	windowWidget.ConnectDestroyQWidget(func() {
		exitApp()
	})
	setupSystray(icon)

	var (
		uiLogo = widgets.NewQLabelFromPointer(windowWidget.FindChild("logo", core.Qt__FindChildrenRecursively).Pointer())

		uiSrvStatus = widgets.NewQLabelFromPointer(windowWidget.FindChild("srv_status", core.Qt__FindChildrenRecursively).Pointer())
		uiXBELPath  = widgets.NewQLabelFromPointer(windowWidget.FindChild("xbel_path", core.Qt__FindChildrenRecursively).Pointer())

		uiSrvAddr   = widgets.NewQLabelFromPointer(windowWidget.FindChild("srv_addr", core.Qt__FindChildrenRecursively).Pointer())
		uiSrvUser   = widgets.NewQLabelFromPointer(windowWidget.FindChild("srv_user", core.Qt__FindChildrenRecursively).Pointer())
		uiSrvPasswd = widgets.NewQLabelFromPointer(windowWidget.FindChild("srv_passwd", core.Qt__FindChildrenRecursively).Pointer())

		uiBtnXBELLocation = widgets.NewQPushButtonFromPointer(windowWidget.FindChild("btn_xbel_localtion", core.Qt__FindChildrenRecursively).Pointer())
		uiBtnServerControl = widgets.NewQPushButtonFromPointer(windowWidget.FindChild("btn_server_control", core.Qt__FindChildrenRecursively).Pointer())
		uiBtnHideControl   = widgets.NewQPushButtonFromPointer(windowWidget.FindChild("btn_hide_tray", core.Qt__FindChildrenRecursively).Pointer())

		uiSettingsStartMinimized = widgets.NewQCheckBoxFromPointer(windowWidget.FindChild("start_minimized", core.Qt__FindChildrenRecursively).Pointer())
		uiSettingsStartServer = widgets.NewQCheckBoxFromPointer(windowWidget.FindChild("start_open", core.Qt__FindChildrenRecursively).Pointer())
	)

	var pixMap = gui.NewQPixmap5(":/qml/logo.png", "", core.Qt__NoFormatConversion)
	uiLogo.SetPixmap(pixMap)

	uiSrvAddr.SetText("http://" + serverConfig.address + ":" + strconv.Itoa(serverConfig.port))
	uiSrvUser.SetText(serverConfig.user);
	uiSrvPasswd.SetText(serverConfig.passwd);
	uiSrvStatus.SetText("<font color='red'>OFF</color>")

	// Button "xxx Server" Action: Turn on or off the server
	uiBtnServerControl.SetIcon(icon)
	uiBtnServerControl.ConnectClicked(func(checked bool) {
		if checked {
			uiSrvStatus.SetText("<font color='green'>ON</color>")
			uiBtnServerControl.SetText("Stop Server")
			uiBtnXBELLocation.SetDisabled(true)
			serverStart()
			return
		}
		uiSrvStatus.SetText("<font color='red'>OFF</color>")
		uiBtnServerControl.SetText("Start Server")
		serverStop()
		uiBtnXBELLocation.SetDisabled(false)
	})

	// Button "Hide to Tray" Action: Hide Main Window
	uiBtnHideControl.ConnectClicked(func(checked bool) {
		windowWidget.Hide()
	})

	// Button + Label "XBEL Location"
	uiXBELPath.SetText(serverConfig.dir)
	uiBtnXBELLocation.ConnectClicked(func(checked bool) {
		log.Printf("Change XBEL dir.")
		dialog := widgets.NewQFileDialog(nil, core.Qt__Dialog)
		dialog.SetViewMode(widgets.QFileDialog__List)
		dialog.SetFileMode(widgets.QFileDialog__Directory)
		dialog.ConnectFileSelected(func(path string) {
			log.Printf("New XBEL path: " + path);
			serverConfig.dir = path
			uiXBELPath.SetText(serverConfig.dir)
			saveAppConfig()
		})
		dialog.Exec()
	})

	// Checkboxes app settings
	if appSettings.startMinimized {
		uiSettingsStartMinimized.SetCheckState(core.Qt__Checked)
	}
	uiSettingsStartMinimized.ConnectClicked(func(checked bool) {
		appSettings.startMinimized = checked
		saveAppConfig()
	})
	if appSettings.startLaunchServer {
		uiSettingsStartServer.SetCheckState(core.Qt__Checked)
	}
	uiSettingsStartServer.ConnectClicked(func(checked bool) {
		appSettings.startLaunchServer = checked
		saveAppConfig()
	})

	// Apply startup settings
	if !appSettings.startMinimized {
		windowWidget.Show()
	}

	if appSettings.startLaunchServer {
		uiBtnServerControl.SetChecked(true)
		uiSrvStatus.SetText("<font color='green'>ON</color>")
		uiBtnServerControl.SetText("Stop Server")
		uiBtnXBELLocation.SetDisabled(true)
		serverStart()
	}

	widgets.QApplication_Exec()
}

func setupSystray(icon *gui.QIcon) {
	// Build the menu
	var menu = widgets.NewQMenu(nil)
	openItem := menu.AddAction("Open LoFloccus")
	openItem.ConnectTriggered(func(checked bool) {
		log.Printf("Tray menu: Open main window.")
		windowWidget.Show()
	})
	quitItem := menu.AddAction("Quit")
	quitItem.ConnectTriggered(func(checked bool) {
		log.Printf("Tray menu: Quit app.")
		exitApp()
	})
	// Set tray icon
	var tray = widgets.NewQSystemTrayIcon(nil)
	tray.SetIcon(icon)
	tray.SetToolTip("LoFloccus App")
	tray.SetContextMenu(menu)
	tray.ConnectActivated(func(reason widgets.QSystemTrayIcon__ActivationReason) {
		if reason != widgets.QSystemTrayIcon__Context {
			windowWidget.Show()
		}
	})
	tray.Show()
}

func exitApp() {
	serverStop()
	core.QCoreApplication_Exit(0)
}

func loadAppConfig() {
	var err error
	cfg, err = ini.Load(configFile)
	if err != nil {
		log.Printf("No config file %v found. Creating a new one from defaults.", err)
		cfg = ini.Empty()

		// Generate random defaults
		rand.Seed(time.Now().UnixNano())
		serverConfig.passwd = "local-" + strconv.Itoa(rand.Intn(999-100)+100);
		serverConfig.port = rand.Intn(65535-40000) + 40000;

		// Save App Config
		saveAppConfig()
	} else {
		// ServerConfig Section
		serverConfig.address = cfg.Section("ServerConfig").Key("address").String();
		serverConfig.port, err = cfg.Section("ServerConfig").Key("port").Int();
		serverConfig.dir = cfg.Section("ServerConfig").Key("dir").String();
		serverConfig.user = cfg.Section("ServerConfig").Key("user").String();
		serverConfig.passwd = cfg.Section("ServerConfig").Key("passwd").String();

		// AppSettings Section
		appSettings.startMinimized, err = cfg.Section("AppSettings").Key("startMinimized").Bool()
		appSettings.startLaunchServer, err = cfg.Section("AppSettings").Key("startLaunchServer").Bool()

	}
}

func saveAppConfig()  {
	// ServerConfig Section
	cfg.Section("ServerConfig").Key("address").SetValue(serverConfig.address)
	cfg.Section("ServerConfig").Key("port").SetValue(strconv.Itoa(serverConfig.port))
	cfg.Section("ServerConfig").Key("dir").SetValue(serverConfig.dir)
	cfg.Section("ServerConfig").Key("passwd").SetValue(serverConfig.passwd)
	cfg.Section("ServerConfig").Key("user").SetValue(serverConfig.user)

	// AppSettings Section
	cfg.Section("AppSettings").Key("startMinimized").SetValue(strconv.FormatBool(appSettings.startMinimized))
	cfg.Section("AppSettings").Key("startLaunchServer").SetValue(strconv.FormatBool(appSettings.startLaunchServer))

	// Save new settings
	cfg.SaveTo(configFile)
}

func serverStart() {

	log.Printf("WEBDAV: Starting...")

	mux := http.NewServeMux()
	srvWebdav := &webdav.Handler{
		FileSystem: webdav.Dir(serverConfig.dir),
		LockSystem: webdav.NewMemLS(),
		Logger: func(request *http.Request, err error) {
			if err != nil {
				log.Printf("WEBDAV [%s]: %s, ERROR: %s\n", request.Method, request.URL, err)
			} else {
				log.Printf("WEBDAV [%s]: %s \n", request.Method, request.URL)
			}
		},
	}

	serverHandle = &http.Server{
		Addr:    serverConfig.address + ":" + strconv.Itoa(serverConfig.port),
		Handler: mux,
	}

	mux.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {

		// Authentication
		username, password, ok := request.BasicAuth()
		if !ok {
			response.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			response.WriteHeader(http.StatusUnauthorized)
			return
		}
		if username != serverConfig.user || password != serverConfig.passwd {
			http.Error(response, "WebDAV: need authorized!", http.StatusUnauthorized)
			return
		}

		// Restrict WebDav to the current folder & read/writes to .xbel files
		if (!strings.HasSuffix(request.RequestURI, ".xbel") && !strings.HasSuffix(request.RequestURI, ".xbel.lock") && request.RequestURI != "/") || (request.RequestURI == "/" && (request.Method != "HEAD" && request.Method != "PROPFIND")) {
			errorFsAccessMsg := "LoFloccus: unauthorized filesystem access detected. LoFloccus WebDAV server is restricted to '*.xbel' files."
			log.Printf(request.RequestURI)
			log.Printf(request.Method)
			log.Printf(errorFsAccessMsg)
			http.Error(response, errorFsAccessMsg, http.StatusUnauthorized)
			return
		}

		srvWebdav.ServeHTTP(response, request)
	})

	go func() {
		if err := serverHandle.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error with WebDAV server: %s", err)
		}
	}()

}

func serverStop() {
	log.Printf("WEBDAV: Shutting down...")
	if err := serverHandle.Shutdown(context.TODO()); err != nil {
		log.Fatalf("WEBDAV: error shutting down - %s", err)
	}
	serverHandle = nil
	log.Printf("WEBDAV: Server is down.")
}
