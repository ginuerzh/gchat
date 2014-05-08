import QtQuick 2.0
import Ubuntu.Components 0.1
import "components"
import "ui"

/*!
    \brief MainView with a Label and Button elements.
*/

MainView {
    // objectName for functional testing purposes (autopilot-qt5)
    id: mainView
    objectName: "mainView"

    // Note! applicationName needs to match the "name" field of the click manifest
    applicationName: "com.ubuntu.developer..dialog"

    /*
     This property enables the application to change orientation
     when the device is rotated. The default is false.
    */
    //automaticOrientation: true
    width: units.gu(45)
    height: units.gu(75)

    Component.onCompleted: {
        pageStack.push(loginTabs)

        // avoid display bug of userTabs title
        pageStack.push(userTabs)
        pageStack.pop()
    }

    PageStack {
        id: pageStack

        Tabs {
            visible: false
            objectName: "loginTabs"
            id: loginTabs
            Tab {
                title: i18n.tr("Login")
                Page {
                    objectName: "loginPage"
                    id: loginPage

                    Login {
                        id: login
                        anchors.horizontalCenter: parent.horizontalCenter
                    }

                    function logined(ok, username, result) {
                        login.logining = false
                        if (ok) {
                            //buddyTab.title = i18n.tr(username)
                            pageStack.pop()
                            pageStack.push(userTabs)
                        } else {
                            console.log(result)
                        }
                    }
                }
            }
            Tab {
                title: i18n.tr("Preferences")
                Page {
                    objectName: "preferencePage"
                    Preference {}
                }
            }
        }

        Tabs{
            objectName: "userTabs"
            id: userTabs
            signal logout()
            visible: false

            Tab {
                id:buddyTab
                title: i18n.tr("Buddies")
                Page {
                    BuddyList {
                        id: buddyList
                        anchors.fill: parent

                        onSelected: pageStack.showDialog(jid, name)
                    }

                    tools: ToolbarItems {
                        ToolbarButton {
                            text: "Logout"
                            iconName: "system-log-out"
                            onTriggered: {
                                userTabs.logout()
                                buddyList.model.clear()
                                pageStack.clear()
                                pageStack.push(loginTabs)
                            }
                        }
                    }
                }
            }

            Tab {
                id: messageTab
                title: i18n.tr("Messages")
                Page {
                    Messages{
                        id: messages
                        anchors.fill: parent
                        maxWidth: mainView.width

                        onSelected: pageStack.showDialog(jid, name)
                    }
                }
            }
        }

        Page {
            id: dialogPage
            title: i18n.tr("Dialog")
            visible: false
            Dialog {
                id: dialog
                anchors.fill: parent
                anchors.margins: 5
                maxWidth: mainView.width
                show: parent.visible
            }
            tools: ToolbarItems {
                locked: true
                opened: true

                TextField {
                    objectName: "msgInput"
                    id: msgInput
                    anchors.verticalCenter: parent.verticalCenter
                    placeholderText: "Input chat message"

                    onAccepted: sendConfirm.triggered(null)
                }
                ToolbarButton {
                    objectName: "sendConfirm"
                    id: sendConfirm
                    text: "Send"
                    signal sended(string jid, string text)

                    iconName: "media-playback-start"
                    onTriggered: {
                        if (msgInput.text.length > 0)
                            sended(dialog.jid, msgInput.text)
                        msgInput.text = ""
                    }
                }
            }
        }

        function showDialog(jid, name) {
            dialogPage.title = i18n.tr(name)
            pageStack.push(dialogPage)
            if (dialog.jid !== jid) {
                dialog.model.clear()
                dialog.jid = jid
                dialog.loaded(jid)
            }
            messages.markRead(jid)
        }
    }

}
