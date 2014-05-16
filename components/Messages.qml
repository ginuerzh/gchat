import QtQuick 2.0


ListView {
    objectName: "messageView"
    id: messages
    spacing: 2
    property int maxWidth
    signal selected(string jid, string name)
    signal accepted(string jid)
    signal declined(string jid)

    delegate: Message {
        width: messages.width
        height: 60
        maxWidth: messages.maxWidth
        user: name
        timestamp: time
        msg: msgText
        avatar: avatarSrc
        unread: unreadStat

        onClicked: messages.selected(jid, name)
    }

    model: ListModel{}

    header: Column {
        //spacing: 1

        Repeater {
            id: repeater

            model: ListModel{}
            Subscribe {
                width: messages.width
                height: 60
                maxWidth: messages.maxWidth
                jid: name
                avatar: avatarSrc
                msg: msgText
                Rectangle {
                    width: parent.width
                    height: 2
                    anchors.bottom: parent.bottom
                    color: "#f0f0f0"
                }
                onAccepted: messages.accepted(jid)
                onDeclined: messages.declined(jid)
            }
        }
        function insert(index, obj) {
            repeater.model.insert(index, obj)
        }
        function remove(jid) {
            console.log("remove", jid)
            for (var i = 0; i < repeater.model.count; i++) {
                var m = repeater.model.get(i)
                if (m.name === jid) {
                    repeater.model.remove(i)
                }
            }
        }
        function clear() {
            repeater.model.clear()
        }
    }

    onAccepted: headerItem.remove(jid)
    onDeclined: headerItem.remove(jid)

    function addMessage(user, msg) {
        for (var i = 0; i < model.count; i++) {
            if (model.get(i).jid === user.jid) {
                model.remove(i)
                break
            }
        }
        model.insert(0, {"jid":user.jid,
                         "name": user.name.length === 0 ? user.jid : user.name,
                         "time": msg.time.format("01-02 15:04"),
                         "msgText": msg.text.trim(),
                         "avatarSrc": user.avatar.length === 0 ? "contact.svg": user.avatar,
                         "unreadStat": msg.unread})
    }

    function removeMessage(jid) {
        for (var i = 0; i < model.count; i++) {
            if (model.get(i).jid === jid) {
                model.remove(i)
                break
            }
        }
    }

    function clearMessages() {
        messages.model.clear()
        messages.headerItem.clear()
    }

    function addSubscribe(msg) {
        headerItem.insert(0, {"jid": msg.jid,
                              "name": msg.jid,
                              "time": msg.time.format("01-02 15:04"),
                              "msgText": msg.text,
                              "avatarSrc": "subscribe.svg"})
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
