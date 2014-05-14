import QtQuick 2.0

Rectangle {
    id: buddy

    property int maxWidth
    property string name
    property string avatar
    property int show
    property string status
    signal clicked()

    property int showChat: 1
    property int showDnd: 2
    property int showAway: 3
    property int showXa: 4
    property int showUnavail: 5

    onShowChanged: {
        switch(show) {
        case showChat:
            statusIcon.color = "lightgreen"
            break
        case showDnd:
            statusIcon.color = "red"
            break
        case showAway:
            statusIcon.color = "gold"
            break
        case showXa:
            statusIcon.color = "gray"
            break
        default:
            statusIcon.color = "white"
            break
        }
    }

    Row {
        spacing: 10
        anchors.verticalCenter: parent.verticalCenter
        Rectangle {
            id: avatarWrapper
            height: buddy.height - 10
            width: height + 5
            Image {
                anchors.fill: parent
                anchors.leftMargin: 5
                source: avatar
            }
        }

        Column {
            Row {
                spacing: 5

                Rectangle {
                    id: statusIcon
                    width: 10
                    height: 10
                    radius: 5
                    //antialiasing: true
                    anchors.verticalCenter: parent.verticalCenter
                }

                Text {
                    text: name
                    elide: Text.ElideRight
                    font.pointSize: 14
                }
            }
            Text {
                width :maxWidth - avatarWrapper.width - 15
                maximumLineCount: 1
                elide: Text.ElideRight
                text: status
            }
        }
    }

    MouseArea {
        id: clickable
        anchors.fill: parent
        onClicked: parent.clicked()
    }
}
