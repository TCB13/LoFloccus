#include "lofloccus.h"

#include <QApplication>

int main(int argc, char *argv[])
{
    QApplication a(argc, argv);
    LoFloccus w;
    //w.show();
    //QApplication::setQuitOnLastWindowClosed(false);
    return a.exec();
}
