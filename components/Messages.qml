import QtQuick 2.0


ListView {
    objectName: "messages"
    id: messages
    spacing: 1
    property int maxWidth
    signal selected(string jid, string name)

    delegate: Message {
        width: messages.width
        height: 64
        maxWidth: messages.maxWidth
        user: name
        timestamp: time
        msg: msgText
        avatar: avatarSrc
        unread: unreadStat

        onClicked: messages.selected(jid, name)
    }

    model: ListModel{}

    function addMessage(msg) {
        for (var i = 0; i < model.count; i++) {
            if (model.get(i).jid === msg.jid) {
                model.remove(i)
                break
            }
        }
        model.insert(0, {"jid":msg.jid,
                         "name": msg.name,
                         "time": msg.time,
                         "msgText": msg.text,
                         "avatarSrc": msg.avatar.length === 0 ? "contact.svg": msg.avatar,
                         "unreadStat": msg.unread})
    }

    function markRead(jid) {
        for (var i = 0; i < model.count; i++) {
            var m = model.get(i)
            if (m.jid === jid) {
                m.unreadStat = false
                break
            }
        }
    }
}
