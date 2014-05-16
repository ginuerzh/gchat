import QtQuick 2.0

ListView {
    objectName: "groupChatView"
    id: groupchat
    property int maxWidth
    property string jid
    property bool show

    signal loaded(string jid)

    spacing: 10

    delegate: MUCBubble {
        width: groupchat.width
        text: str
        username: name
        timestamp: time
        avatar: avatarSrc
        floatRight: !isReceive
        maxWidth: groupchat.maxWidth
    }
    model:ListModel {}

    function addBubble(user, bubble) {
        model.append({"str": bubble.text,
                      "name": user.name,
                       "isReceive": user.jid === jid,
                      "time": bubble.time.format("01-02 15:04"),
                      "avatarSrc": user.avatar.length === 0 ? "contact.svg": user.avatar})
        positionViewAtEnd()
    }
}

