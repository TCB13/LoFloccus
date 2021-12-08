#ifndef LOFLOCCUS_H
#define LOFLOCCUS_H

#include <QMainWindow>
#include <QSettings>
#include <QSystemTrayIcon>
#include <QMenu>

#ifdef Q_OS_DARWIN
#include "platformdarwin.h"
#endif

QT_BEGIN_NAMESPACE
namespace Ui { class LoFloccus; }
QT_END_NAMESPACE

class LoFloccus : public QMainWindow
{
    Q_OBJECT

public:
    LoFloccus(QWidget *parent = nullptr);
    ~LoFloccus();
    QIcon appIcon;

private slots:
    void on_btn_xbel_localtion_clicked();
    void on_startminimized_clicked();
    void on_hidetosystray_clicked();
    void on_sharednetwork_clicked();
    void on_portablemode_clicked();

protected:
    void closeEvent(QCloseEvent *event) override;
    void hideEvent(QHideEvent *event) override;

private:
    bool minimizeAndCloseToTray;
    bool running;

    QSettings *settings;
    QSystemTrayIcon *sysTray;

    #ifdef Q_OS_DARWIN
    PlatformDarwin *darwinBridge;
    #endif

    Ui::LoFloccus *ui;

    void initSettings(bool makeExistingSettingsPortable, bool makeExistingSettingsLocal);
    void reloadUiState();
    void startServer();
    void stopServer();
    void restartServer();
    void initSystray();
    QList<QString> getSystemIPAddresses(bool locals, bool v4, bool v6);
};
#endif // LOFLOCCUS_H
