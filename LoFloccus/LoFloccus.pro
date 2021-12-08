QT       += core gui network

greaterThan(QT_MAJOR_VERSION, 4): QT += widgets

CONFIG += c++11

# The following define makes your compiler emit warnings if you use
# any Qt feature that has been marked deprecated (the exact warnings
# depend on your compiler). Please consult the documentation of the
# deprecated API in order to know how to port your code away from it.
DEFINES += QT_DEPRECATED_WARNINGS

# You can also make your code fail to compile if it uses deprecated APIs.
# In order to do so, uncomment the following line.
# You can also select to disable deprecated APIs only up to a certain version of Qt.
#DEFINES += QT_DISABLE_DEPRECATED_BEFORE=0x060000    # disables all the APIs deprecated before Qt 6.0.0

SOURCES += \
    main.cpp \
    lofloccus.cpp

HEADERS += \
    lofloccus.h

FORMS += \
    lofloccus.ui

# Generic for all builds
VERSION = 1.2.1
DEFINES += APP_VERSION=\\\"$$VERSION\\\"
INCLUDEPATH += $${PWD}/libs

# Windows Build
win32 {
    RC_ICONS = assets/icon.ico
    QMAKE_TARGET_PRODUCT = LoFloccus
    QMAKE_TARGET_DESCRIPTION = LoFloccus
    QMAKE_TARGET_COPYRIGHT = Copyright © 2019-2021 TCB13
    LIBS += -L$${PWD}/libs -lLoFloccusDavWin64
}

# macOS Build
mac {
    ICON = assets/icon.icns
    QMAKE_TARGET_BUNDLE_PREFIX = "com.tcb13"
    LIBS += -L$${PWD}/libs -lLoFloccusDavDarwin

    SOURCES += platformdarwin.mm
    HEADERS += platformdarwin.h
    LIBS += -framework Foundation
    LIBS += -framework AppKit
}


# Default rules for deployment.
qnx: target.path = /tmp/$${TARGET}/bin
else: unix:!android: target.path = /opt/$${TARGET}/bin
!isEmpty(target.path): INSTALLS += target

DISTFILES +=

RESOURCES += \
    resources.qrc

# Static Builds
QTPREFIX=$$[QT_INSTALL_PREFIX]
equals(QTPREFIX, "C:/Qt-Static/Qt-5.14.2") || equals(QTPREFIX, "/Users/tcb13/Qt-Static/Qt-5.12.12") {
    message("--STATIC BUILD--")
    CONFIG += qt static
    win32 {
        QMAKE_LFLAGS += -static-libgcc -static-libstdc++
    }

} else {
    message("--NON-STATIC BUILD--")
}
