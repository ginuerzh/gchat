import QtQuick 2.0

Rectangle {
    id: button
    property alias text: textField.text
    signal clicked()
    radius: 5

    Text {
        id: textField
        anchors.centerIn: parent
    }

    MouseArea {
        anchors.fill: parent
        onClicked: button.clicked()
    }
}
