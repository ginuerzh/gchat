import QtQuick 2.0

ListView {
    objectName: "roomView"
    id: roomlist
    spacing: 2
    signal selected(string jid, string name)
    property int maxWidth

    delegate: Room {
        width: roomlist.width
        height: 50
        name: roomName
        maxWidth: roomlist.maxWidth
        onClicked: {
            roomlist.selected(jid, name)
        }
    }

    model: ListModel{}

    function appendRoom(room) {
        model.append({"jid":room.jid,
                      "roomName": room.jid.split("@", 1)[0]})
    }

    function clearRooms() {
        model.clear()
    }
}
