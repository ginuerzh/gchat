import QtQuick 2.0

ListView {
    objectName: "dialogView"
    id: dialog
    property int maxWidth
    property string jid
    property bool show

    signal loaded(string jid)

    spacing: 10

    delegate: Bubble {
        width: dialog.width
        text: str
        timestamp: time
        avatar: avatarSrc
        floatRight: !isReceive
        maxWidth: dialog.maxWidth
    }
    model:ListModel {}

    function addBubble(bubble) {
        model.append({"str": bubble.text,
                      "isReceive": bubble.jid === jid,
                      "time": bubble.time,
                      "avatarSrc": bubble.avatar.length === 0 ? "contact.svg": bubble.avatar})
    }
}

