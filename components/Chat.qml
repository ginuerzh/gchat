import QtQuick 2.0

ListView {
    objectName: "chatView"
    id: chatView
    property int maxWidth
    property string jid
    property bool show

    signal loaded(string jid)

    spacing: 10

    delegate: Bubble {
        width: chatView.width
        text: str
        timestamp: time
        avatar: avatarSrc
        floatRight: !isReceive
        maxWidth: chatView.maxWidth
    }
    model:ListModel {}

    function addBubble(user, bubble) {
        model.append({"str": bubble.text,
                      "isReceive": bubble.jid === jid,
                      "time": bubble.time.format("01-02 15:04"),
                      "avatarSrc": user.avatar.length === 0 ? "contact.svg": user.avatar})
        positionViewAtEnd()
    }
}

