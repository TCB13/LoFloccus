#ifndef LOFLOCCUS_H
#define LOFLOCCUS_H

#include <QMainWindow>
#include <QSettings>
#include <QSystemTrayIcon>
#include <QMenu>

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
    void on_btn_server_control_clicked();
    void on_btn_xbel_localtion_clicked();

    void on_startopen_clicked();

    void on_startminimized_clicked();

    void on_hidetosystray_clicked();

    void on_sharednetwork_clicked();

protected:
    #ifdef Q_OS_WIN
    void closeEvent(QCloseEvent *event) override;
    void hideEvent(QHideEvent *event) override;
    #endif

private:
    bool minimizeAndCloseToTray;
    bool running;

    QSettings *settings;
    QSystemTrayIcon *sysTray;

    Ui::LoFloccus *ui;
    void initSettings();
    void reloadUiState();

    void startServer();
    void stopServer();
    void restartServer();
    void initSystray();
    QList<QString> getSystemIPAddresses(bool locals, bool v4, bool v6);
};
#endif // LOFLOCCUS_H
