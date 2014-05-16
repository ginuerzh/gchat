import QtQuick 2.0

Rectangle {
    id: room
    property int maxWidth

    property string name

    signal clicked()

    Row {
        anchors.verticalCenter: parent.verticalCenter
        anchors.left: parent.left
        anchors.right: parent.right
        anchors.leftMargin: 5
        anchors.rightMargin: 5

        Text {
            elide: Text.ElideRight
            text: name
        }
    }

    MouseArea {
        id: clickable
        anchors.fill: parent
        onClicked: {
            parent.clicked()
        }
    }
}
