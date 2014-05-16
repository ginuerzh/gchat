import QtQuick 2.0


ListView {
    objectName: "roomInfoView"
    id: roomInfo
    spacing: 1
    property int maxWidth
    property string jid
    property string name
    signal loaded(string jid)

    header: Column {
        width: roomInfo.maxWidth
        spacing: 5

        Row {
            anchors.left: parent.left
            anchors.right: parent.right
            anchors.leftMargin: 5
            anchors.rightMargin: 5
            spacing: 5

            Text {
                id: nameTitle
                text: "Name:"
            }
            Text {
                text: roomInfo.name
                width: roomInfo.maxWidth - nameTitle.width - 15
                wrapMode: Text.WrapAtWordBoundaryOrAnywhere
            }
        }
        Row {
            anchors.left: parent.left
            anchors.right: parent.right
            anchors.leftMargin: 5
            anchors.rightMargin: 5
            spacing: 5
            Text {
                id: descTitle
                text: "Description:"
            }
            Text {
                id: roomDesc
                width: roomInfo.maxWidth - descTitle.width - 15
                wrapMode: Text.WrapAtWordBoundaryOrAnywhere
            }
        }
        Row {
            anchors.left: parent.left
            anchors.right: parent.right
            anchors.leftMargin: 5
            anchors.rightMargin: 5
            spacing: 5
            Text {
                id: subjectTitle
                text: "Subject:"
            }
            Text {
                id: roomSubject
                width: roomInfo.maxWidth - subjectTitle.width - 15
                wrapMode: Text.WrapAtWordBoundaryOrAnywhere
            }
        }

        Row {
            anchors.left: parent.left
            anchors.right: parent.right
            anchors.leftMargin: 5
            anchors.rightMargin: 5
            spacing: 5
            Text {
                text:"Occupants:"
            }
            Text {
                id: roomOccupants
            }
        }

        function updateRoomInfo(info) {
            roomDesc.text = info.description
            roomSubject.text = info.subject
            roomOccupants.text = info.occupants
        }
    }

    delegate: Occupant {
        width: messages.width
        height: 40
        maxWidth: messages.maxWidth
        nickname: name
    }

    model: ListModel{}

    function setRoomInfo(info) {
        headerItem.updateRoomInfo(info)
    }

    function appendOccupant(occupant) {
        model.append({"name": occupant})
    }
}
