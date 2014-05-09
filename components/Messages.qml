import QtQuick 2.0


ListView {
    objectName: "messageView"
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

    function addMessage(user, msg) {
        for (var i = 0; i < model.count; i++) {
            if (model.get(i).jid === user.jid) {
                model.remove(i)
                break
            }
        }
        model.insert(0, {"jid":user.jid,
                         "name": user.name,
                         "time": msg.time.format("15:04"),
                         "msgText": msg.text,
                         "avatarSrc": user.avatar.length === 0 ? "contact.svg": user.avatar,
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
