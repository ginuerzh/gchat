import QtQuick 2.0

Item {
    id: bubble
    property string text
    property string username
    property bool floatRight
    property url avatar
    property string timestamp
    property int maxWidth

    height: Math.max(avatarIcon.height, rect.height)

    Rectangle {
        id: avatarIcon
        height:52
        width: height
        Image {
            anchors.fill: parent
            source: avatar
        }
    }

    Rectangle {
        id: rect
        anchors.top: bubble.top
        Item {
            id: title
            anchors.left: parent.left
            anchors.leftMargin: 5
            width: Math.max(name_field.width + timestamp_field.width + 10, text_field.width)
            height: timestamp_field.height

            Text {
                id: name_field
                text: username
                anchors.left: parent.left
            }

            Text {
                id: timestamp_field
                text: timestamp
                anchors.right: parent.right
            }
        }

        Text {
            id: text_field
            text: bubble.text
            wrapMode: Text.WrapAtWordBoundaryOrAnywhere
            anchors.left: parent.left
            anchors.bottom: parent.bottom
            anchors.margins: 5
        }

        Component.onCompleted: {
            text_field.width = Math.min(text_field.contentWidth, bubble.maxWidth - avatarIcon.width - 20)
            text_field.height = text_field.contentHeight

            width = Math.max(text_field.width, title.width) + 10
            height = Math.max(text_field.height + title.height + 15, avatarIcon.height)
        }
    }

    Component.onCompleted: {
        if (floatRight == true) {
            avatarIcon.anchors.right = bubble.right
            avatarIcon.anchors.bottom = bubble.bottom
            rect.anchors.right = avatarIcon.left
            rect.color = "aquamarine"
        } else {
            avatarIcon.anchors.left = bubble.left
            rect.anchors.left = avatarIcon.right
        }
    }
}
