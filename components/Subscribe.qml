import QtQuick 2.0

Rectangle {
    id: subscribe
    property int maxWidth

    property string jid
    property string msg
    property string avatar
    property string timestamp

    signal accepted(string jid)
    signal declined(string jid)

    Row {
        spacing: 5
        anchors.verticalCenter: parent.verticalCenter

        Rectangle {
            id: avatarIcon
            height: subscribe.height
            width: height
            Image {
                anchors.fill: parent
                anchors.margins: 5
                source: avatar
            }
        }

        Column {
            width: maxWidth - avatarIcon.width - 5

            Rectangle {
                width: parent.width
                height: 30
                Text {
                    anchors.left: parent.left
                    anchors.verticalCenter: parent.verticalCenter
                    width: parent.width - timeField.width - 5
                    elide: Text.ElideRight
                    text: jid
                }

                Text {
                    id: timeField
                    text: timestamp
                    anchors.right: parent.right
                    anchors.verticalCenter: parent.verticalCenter
                    anchors.rightMargin: 5
                    font.bold: true
                }
            }

            Row {
                spacing: 5
                Text {
                    width: 140
                    maximumLineCount: 1
                    elide: Text.ElideRight
                    text: msg
                }

                Button {
                    text: "Accept"
                    width: 70
                    height: 25
                    color: "palegreen"
                    onClicked: subscribe.accepted(jid)
                }
                Button {
                    text: "Decline"
                    width: 70
                    height: 25
                    color: "salmon"
                    onClicked: subscribe.declined(jid)
                }
            }
        }
    }
}
