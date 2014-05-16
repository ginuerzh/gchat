import QtQuick 2.0

Rectangle {
    id: occupant
    property int maxWidth
    property string nickname
    signal clicked()

    Row {
        anchors.verticalCenter: parent.verticalCenter
        anchors.left: parent.left
        anchors.right: parent.right
        anchors.leftMargin: 5
        anchors.rightMargin: 5


        Text {
            anchors.verticalCenter: parent.verticalCenter
            elide: Text.ElideRight
            text: nickname
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
