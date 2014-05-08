import QtQuick 2.0

Rectangle {
    id: message
    property int maxWidth

    property string user
    property string timestamp
    property string msg
    property string avatar
    property bool unread
    signal clicked()

    Row {
        spacing: 10
        anchors.verticalCenter: parent.verticalCenter
        Rectangle {
            id: avatarIcon
            height: message.height - 10
            width: height + 5
            Image {
                anchors.fill: parent
                anchors.leftMargin: 5
                source: avatar
            }
        }

        Column {
            spacing: 5
            width: maxWidth - avatarIcon.width - 25

            Item {
                width: parent.width
                height: childrenRect.height

                Text {
                    anchors.left: parent.left
                    elide: Text.ElideRight
                    text: user
                    font.bold: message.unread
                    font.pointSize: 14
                }

                Text {
                    text: timestamp
                    anchors.right: parent.right
                    anchors.leftMargin: 5
                }
            }
            Text {
                width: parent.width
                elide: Text.ElideRight
                text: msgText
                font.bold: message.unread
            }
        }
    }

    MouseArea {
        id: clickable
        anchors.fill: parent
        onClicked: {
            message.unread = false
            parent.clicked()
        }
    }
}
