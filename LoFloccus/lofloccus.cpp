#include "lofloccus.h"
#include "ui_lofloccus.h"
#include <QDebug>
#include <QSettings>
#include <QRandomGenerator>
#include <QMenu>
#include <QAction>
#include <QMessageBox>
#include <QCloseEvent>
#include <QHideEvent>
#include <QEvent>
#include <QFileDialog>
#include <QNetworkInterface>
#include <QStandardPaths>

#ifdef Q_OS_WIN
#include "libLoFloccusDavWin64.h"
#endif
#ifdef Q_OS_DARWIN
#include "libLoFloccusDavDarwin.h"
#endif

LoFloccus::LoFloccus(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::LoFloccus)
{
    ui->setupUi(this);
    appIcon = QIcon(":/assets/icon.ico");
    running = false;

    #ifdef Q_OS_DARWIN
    darwinBridge = new PlatformDarwin();
    #endif

    // Make sure nobody can resize the app under any platform
    this->setFixedSize(this->width(), this->height());
    this->setWindowTitle(this->windowTitle() + " - v" + APP_VERSION);

    // Fetch settings from storage and/or write defaults
    this->initSettings(false, false);

    // Populate UI with the loaded settings
    this->reloadUiState();

    // Build the systray menu
    this->initSystray();

    // Deal with previously saved settings
    if (!settings->value("startminimized").toBool()) {
        show();
    }
    if (settings->value("startminimized").toBool() && !settings->value("hidetosystray").toBool()) {
        showMinimized();
    }
    if (settings->value("startminimized").toBool() && settings->value("hidetosystray").toBool()) {
       #ifdef Q_OS_DARWIN
       darwinBridge->makeAppAccessory();
       #endif
       sysTray->showMessage("LoFloccus", "LoFloccus is running in the background. Click the menu for more options.", appIcon);
    }

    // Start the server with the app
    this->startServer();
}

void LoFloccus::closeEvent(QCloseEvent *event)
{
    if (minimizeAndCloseToTray) {
        hide();
        event->ignore();
    } else {
        event->accept();
    }
}

void LoFloccus::hideEvent(QHideEvent *event)
{
    event->accept();
    if (minimizeAndCloseToTray) {
        hide();
        #ifdef Q_OS_DARWIN
        darwinBridge->makeAppAccessory();
        #endif
        sysTray->showMessage("LoFloccus", "LoFloccus is running in the background. Click the menu for more options.", appIcon);
    }
}

LoFloccus::~LoFloccus()
{
    // Stop the server on close
    if (running) {
        this->stopServer();
    }

    delete ui;
}

void LoFloccus::initSettings(bool makeExistingSettingsPortable = false, bool makeExistingSettingsLocal = false)
{   
    QString localSettingsFile = QStandardPaths::writableLocation(QStandardPaths::AppConfigLocation) + "/settings.ini";
    QString portableSettingsFile = "lofloccus-settings.ini";

    // Deal with settings location move
    if (makeExistingSettingsPortable) {
        if (QFile::exists(portableSettingsFile)) {
            QFile::remove(portableSettingsFile);
        }
        QFile::copy(localSettingsFile, portableSettingsFile);
        QFile::remove(localSettingsFile);
    }
    if (makeExistingSettingsLocal) {
        if (QFile::exists(localSettingsFile)) {
            QFile::remove(localSettingsFile);
        }
        QFile::copy(portableSettingsFile, localSettingsFile);
        QFile::remove(portableSettingsFile);
    }

    // Initialize a shared settings object with the appropriate path
    settings = new QSettings(QFile::exists(portableSettingsFile) ? portableSettingsFile : localSettingsFile, QSettings::IniFormat);

    // Generate random defaults for port and password
    QString defaultPort = QString::number(QRandomGenerator::global()->bounded(40000, 65535));
    QString defaultPasswd = "fpass-" + QString::number(QRandomGenerator::global()->bounded(100000, 999999));

    // Set default values for each setting. This will write default values or
    // keep existing ones making sure our settings on disk are always complete
    settings->setValue("serveraddr", settings->value("serveraddr", "127.0.0.1"));
    settings->setValue("serverport", settings->value("serverport", defaultPort));
    settings->setValue("serverpath", settings->value("serverpath", QCoreApplication::applicationDirPath()));
    settings->setValue("serveruser", settings->value("serveruser", "floccus"));
    settings->setValue("serverpasswd", settings->value("serverpasswd", defaultPasswd));

    settings->setValue("startminimized", settings->value("startminimized", false));
    settings->setValue("hidetosystray", settings->value("hidetosystray", false));
    settings->setValue("sharednetwork", settings->value("sharednetwork", false));
    settings->setValue("portablemode", settings->value("portablemode", false));
}

void LoFloccus::reloadUiState()
{
    ui->xbel_path->setText(settings->value("serverpath").toString());

    // Take care of the server addresses, might be just local or all IPs of the machine
    QString serverPort = settings->value("serverport").toString();
    QList<QString> ips = { "127.0.0.1" };
    if (settings->value("sharednetwork").toBool()) {
        ips << getSystemIPAddresses(false, true, false);
    }
    ui->srv_addr->setText("http://" + ips.mid(0,2).join(":" + serverPort + " or http://") + ":" + serverPort + (ips.length() > 2 ? " (...)" : ""));
    ui->srv_addr->setToolTip("LoFloccus server is running at:\r\nhttp://" + ips.join(":" + serverPort + "\r\nhttp://") + ":" + serverPort);

    ui->srv_user->setText(settings->value("serveruser").toString());
    ui->srv_passwd->setText(settings->value("serverpasswd").toString());

    ui->startminimized->setChecked(settings->value("startminimized").toBool());
    ui->hidetosystray->setChecked(settings->value("hidetosystray").toBool());
    ui->sharednetwork->setChecked(settings->value("sharednetwork").toBool());
    ui->portablemode->setChecked(settings->value("portablemode").toBool());
}

void LoFloccus::initSystray()
{
    minimizeAndCloseToTray = settings->value("hidetosystray").toBool();
    QAction *exitAction = new QAction("Quit", this);
    connect(exitAction, &QAction::triggered, [this]() {
        minimizeAndCloseToTray = false;
        close();
    });

    QAction *openAction = new QAction("Open LoFloccus", this);
    connect(openAction, &QAction::triggered, [this]() {
        #ifdef Q_OS_DARWIN
        darwinBridge->makeAppRegular();
        #endif
        showNormal();
        activateWindow();
    });

    QMenu *trayIconMenu = new QMenu(this);
    trayIconMenu->addAction(openAction);
    trayIconMenu->addAction(exitAction);

    sysTray = new QSystemTrayIcon(this);
    sysTray->setContextMenu(trayIconMenu);
    sysTray->setIcon(appIcon);
    sysTray->show();

    // Force opening the context menu on icon click under Windows
    #ifdef Q_OS_WIN
    connect(sysTray, &QSystemTrayIcon::activated, [this]() {
        /* QSystemTrayIcon::ActivationReason reason
         * if (reason == QSystemTrayIcon::DoubleClick) {
            showNormal();
            activateWindow();
            return;
        }*/
        sysTray->contextMenu()->popup(QCursor::pos());
    });
    #endif
}

void LoFloccus::restartServer()
{
    if (running) {
        this->stopServer();
        this->startServer();
    }
}

void LoFloccus::startServer()
{
    serverStart(settings->value("serveraddr").toString().toUtf8().data(),
                settings->value("serverport").toString().toUtf8().data(),
                settings->value("serverpath").toString().toUtf8().data(),
                settings->value("serveruser").toString().toUtf8().data(),
                settings->value("serverpasswd").toString().toUtf8().data()
                );
    running = true;
    this->reloadUiState();
}

void LoFloccus::stopServer()
{
    serverStop();
    running = false;
    this->reloadUiState();
}

QList<QString> LoFloccus::getSystemIPAddresses(bool locals = true, bool v4 = true, bool v6 = true)
{
    QList<QString> returnList;
    QList<QHostAddress> addresses = QNetworkInterface::allAddresses();
    for (int i=0; i < addresses.count(); i++) {
        if(!locals && addresses[i].isLoopback()) {
            continue;
        }
        if (v4 && addresses[i].protocol() == QAbstractSocket::IPv4Protocol) {
            returnList.append(addresses[i].toString());
        }
        if (v6 && addresses[i].protocol() == QAbstractSocket::IPv6Protocol) {
            returnList.append(addresses[i].toString());
        }
    }
    return returnList;
}

void LoFloccus::on_btn_xbel_localtion_clicked()
{
    QString dir = QFileDialog::getExistingDirectory(this,
                                                    "Pick a Directory",
                                                    settings->value("serverpath").toString(),
                                                    QFileDialog::ShowDirsOnly);
    if (dir.isEmpty()) {
        return;
    }
    settings->setValue("serverpath", dir);
    this->restartServer();
}



void LoFloccus::on_portablemode_clicked()
{
    if (ui->portablemode->isChecked()) {
        this->initSettings(true, false);
    } else {
        this->initSettings(false, true);
    }
    settings->setValue("portablemode", ui->portablemode->isChecked());
}

void LoFloccus::on_startminimized_clicked()
{
    settings->setValue("startminimized", ui->startminimized->isChecked());
}

void LoFloccus::on_hidetosystray_clicked()
{
    settings->setValue("hidetosystray", ui->hidetosystray->isChecked());
    minimizeAndCloseToTray = ui->hidetosystray->isChecked();
}

void LoFloccus::on_sharednetwork_clicked()
{
    settings->setValue("sharednetwork", ui->sharednetwork->isChecked());
    settings->setValue("serveraddr", ui->sharednetwork->isChecked() ? "0.0.0.0" : "127.0.0.1");
    this->restartServer();
}

