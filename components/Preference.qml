import QtQuick 2.0
import Ubuntu.Components 0.1

Column {
    anchors.left: parent.left
    anchors.right: parent.right
    anchors.top: parent.top
    anchors.leftMargin: 10
    anchors.rightMargin: 10
    anchors.topMargin: 20
    spacing: 10

    Item {
        width: parent.width
        height: 40
        Text {
            text: "Server"
            anchors.left: parent.left
            anchors.verticalCenter: parent.verticalCenter
        }
        TextField {
            objectName: "serverAddr"
            id: serverAddr
            placeholderText: "talk.google.com"
            anchors.right: parent.right
            anchors.verticalCenter: parent.verticalCenter
            KeyNavigation.tab: serverPort
        }
    }

    Item {
        width: parent.width
        height: 40
        Text {
            text: "Server Port"
            anchors.left: parent.left
            anchors.verticalCenter: parent.verticalCenter
        }
        TextField {
            objectName: "serverPort"
            id: serverPort
            placeholderText: "443"
            anchors.right: parent.right
            anchors.verticalCenter: parent.verticalCenter
            inputMethodHints: Qt.ImhDigitsOnly
            KeyNavigation.tab:resource
        }
    }

    Item {
        width: parent.width
        height: 40
        Text {
            text: "Resource"
            anchors.left: parent.left
            anchors.verticalCenter: parent.verticalCenter
        }
        TextField {
            objectName: "resource"
            id: resource
            anchors.right: parent.right
            anchors.verticalCenter: parent.verticalCenter
            KeyNavigation.tab:proxyServerAddr
        }
    }

    Item {
        width: parent.width
        height: 40

        Text {
            text: "Use the old SSL method"
            anchors.left: parent.left
            anchors.verticalCenter: parent.verticalCenter
        }
        Switch {
            objectName: "sslSwitch"
            id: sslSwitch
            anchors.right: parent.right
            anchors.verticalCenter: parent.verticalCenter
        }
    }

    Item {
        width: parent.width
        height: 40

        Text {
            text: "Use Proxy"
            anchors.left: parent.left
            anchors.verticalCenter: parent.verticalCenter
        }
        Switch {
            objectName: "proxySwitch"
            id: proxySwitch
            anchors.right: parent.right
            anchors.verticalCenter: parent.verticalCenter
        }
    }
    Item {
        width: parent.width
        height: 40
        visible: proxySwitch.checked

        Text {
            text: "Use System Proxy"
            anchors.left: parent.left
            anchors.verticalCenter: parent.verticalCenter
        }
        Switch {
            objectName: "sysProxySwitch"
            id: sysProxySwitch
            anchors.right: parent.right
            anchors.verticalCenter: parent.verticalCenter
        }
    }

    Item {
        width: parent.width
        height: 40
        visible: proxySwitch.checked && !sysProxySwitch.checked

        Text {
            text: "Proxy Server"
            anchors.left: parent.left
            anchors.verticalCenter: parent.verticalCenter
        }
        TextField {
            objectName: "proxyServerAddr"
            id: proxyServerAddr
            anchors.right: parent.right
            anchors.verticalCenter: parent.verticalCenter
            KeyNavigation.tab: proxyServerPort
        }
    }
    Item {
        width: parent.width
        height: 40
        visible: proxySwitch.checked && !sysProxySwitch.checked

        Text {
            text: "Proxy Server Port"
            anchors.left: parent.left
            anchors.verticalCenter: parent.verticalCenter
        }
        TextField {
            objectName: "proxyServerPort"
            id: proxyServerPort
            anchors.right: parent.right
            anchors.verticalCenter: parent.verticalCenter
            inputMethodHints: Qt.ImhDigitsOnly
        }
    }

    Item {
        width: parent.width
        height: 40

        Text {
            text: "Auto Login"
            anchors.left: parent.left
            anchors.verticalCenter: parent.verticalCenter
        }
        Switch {
            objectName: "autoLogin"
            id: autoLogin
            anchors.right: parent.right
            anchors.verticalCenter: parent.verticalCenter
        }
    }
}
