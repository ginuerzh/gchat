import QtQuick 2.0

Item {
    id: bubble
    property string text
    property bool floatRight
    property url avatar
    property string timestamp
    property int maxWidth

    height: Math.max(avatarIcon.height, rect.height)

    Rectangle {
        id: avatarIcon
        height:52
        width: height
        //anchors.top: bubble.top
        Image {
            anchors.fill: parent
            source: avatar
        }
    }

    Rectangle {
        id: rect
        anchors.top: avatarIcon.top
        Text {
            id: text_field
            text: bubble.text
            wrapMode: Text.WrapAtWordBoundaryOrAnywhere
            anchors.left: parent.left
            anchors.top: parent.top
            anchors.margins: 5
        }

        Text {
            id: timestamp_field
            text: timestamp
            anchors.left: parent.left
            anchors.bottom: parent.bottom
            anchors.margins: 5
        }

        Component.onCompleted: {
            text_field.width = Math.min(text_field.contentWidth, bubble.maxWidth - avatarIcon.width - 20)
            text_field.height = text_field.contentHeight

            width = Math.max(text_field.width, timestamp_field.width) + 10
            height = Math.max(text_field.height + timestamp_field.height + 15, avatarIcon.height)
        }
    }

    Component.onCompleted: {
        if (floatRight == true) {
            avatarIcon.anchors.right = bubble.right
            rect.anchors.right = avatarIcon.left
        } else {
            avatarIcon.anchors.left = bubble.left
            rect.anchors.left = avatarIcon.right
            rect.color = "palegreen"
        }
    }
}
