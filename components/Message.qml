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
        spacing: 5
        anchors.verticalCenter: parent.verticalCenter

        Rectangle {
            id: avatarIcon
            height: message.height
            width: height
            Image {
                anchors.fill: parent
                anchors.margins: 5
                source: avatar
            }
        }

        Column {
            spacing: 5
            width: maxWidth - avatarIcon.width - 5

            Rectangle {
                width: parent.width
                height: 30
                Text {
                    anchors.left: parent.left
                    anchors.verticalCenter: parent.verticalCenter
                    width: parent.width - timeField.width - 5
                    elide: Text.ElideRight
                    text: user
                    font.bold: message.unread
                }

                Text {
                    id: timeField
                    text: timestamp
                    anchors.right: parent.right
                    anchors.verticalCenter: parent.verticalCenter
                    anchors.rightMargin: 5
                    font.bold: message.unread
                }
            }

            Text {
                width: parent.width
                maximumLineCount: 1
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
