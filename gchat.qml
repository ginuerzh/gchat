import QtQuick 2.0
import Ubuntu.Components 0.1
import "components"

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
                id: messageTab
                title: i18n.tr("Messages")
                Page {
                    Messages{
                        id: messages
                        anchors.fill: parent
                        maxWidth: mainView.width

                        onSelected: pageStack.showChat(jid, name)
                    }

                     tools: ToolbarItems {
                         locked: true
                         opened: false
                     }
                }
            }
            Tab {
                id:buddyTab
                title: i18n.tr("Buddies")
                Page {
                    BuddyList {
                        id: buddyList
                        anchors.fill: parent
                        maxWidth: mainView.width

                        onSelected: pageStack.showChat(jid, name)
                    }

                    tools: ToolbarItems {
                        ToolbarButton {
                            text: "Logout"
                            iconName: "system-log-out"
                            onTriggered: {
                                mainView.logout()
                            }
                        }
                    }
                }
            }
            Tab {
                id: roomTab
                title:i18n.tr("Chat Rooms")
                Page {
                    RoomList {
                        id: roomList
                        anchors.fill: parent
                        maxWidth: mainView.width
                        onSelected: pageStack.showRoomInfo(jid, name)
                    }
                }
            }
        }

        Tabs {
            objectName: "chatRoomTabs"
            id: chatRoomTabs
            visible: false
            Tab {
                id: roomInfoTab
                title: i18n.tr("Chat Room Info")
                Page {
                    RoomInfo {
                        id: roomInfo
                        anchors.fill: parent
                        maxWidth: mainView.width
                    }
                }
            }
            Tab {
                title: i18n.tr("Chat")
                Page {
                    MUC {
                        objectName: "groupChatView"
                        id: groupChat
                        anchors.fill: parent
                        anchors.margins: 5
                        maxWidth: mainView.width
                    }
                    tools: ToolbarItems {
                        locked: true
                        opened: true

                        TextField {
                            id: mucMsgInput
                            anchors.verticalCenter: parent.verticalCenter
                            placeholderText: "Input chat message"

                            onAccepted: mucSend.triggered(null)
                        }
                        ToolbarButton {
                            objectName: "mucSend"
                            id: mucSend
                            text: "Send"
                            signal sended(string jid, string text)

                            iconName: "media-playback-start"
                            onTriggered: {
                                if (mucMsgInput.text.length > 0)
                                    sended(groupChat.jid, mucMsgInput.text)
                                mucMsgInput.text = ""
                            }
                        }
                    }
                }
            }
        }

        Page {
            id: chatPage
            title: i18n.tr("Chat")
            visible: false
            Chat {
                id: chat
                anchors.fill: parent
                anchors.margins: 5
                maxWidth: mainView.width
                show: parent.visible
            }
            tools: ToolbarItems {
                locked: true
                opened: true

                TextField {
                    id: msgInput
                    anchors.verticalCenter: parent.verticalCenter
                    placeholderText: "Input chat message"

                    onAccepted: chatSend.triggered(null)
                }
                ToolbarButton {
                    objectName: "chatSend"
                    id: chatSend
                    text: "Send"
                    signal sended(string jid, string text)

                    iconName: "media-playback-start"
                    onTriggered: {
                        if (msgInput.text.length > 0)
                            sended(chat.jid, msgInput.text)
                        msgInput.text = ""
                    }
                }
            }
        }

        function showChat(jid, name) {
            chatPage.title = i18n.tr(name)
            pageStack.push(chatPage)
            if (chat.jid !== jid) {
                chat.model.clear()
                chat.jid = jid
                chat.loaded(jid)
            }
            messages.markRead(jid)
        }

        function showRoomInfo(jid, name) {
            roomInfoTab.title = i18n.tr(jid.split("@", 1)[0])
            pageStack.push(chatRoomTabs)
            if (roomInfo.jid !== jid) {
                roomInfo.model.clear()
                groupChat.jid = jid
                roomInfo.jid = jid
                roomInfo.name = name
                roomInfo.loaded(jid)
            }
        }
    }

    function logout() {
        userTabs.logout()
        buddyList.clearBuddies()
        messages.clearMessages()
        roomList.clearRooms()

        pageStack.clear()
        pageStack.push(loginTabs)
    }
}
