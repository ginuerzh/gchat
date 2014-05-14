import QtQuick 2.0

ListView {
    objectName: "buddyView"
    id: buddylist
    spacing: 1
    signal selected(string jid, string name)
    property int maxWidth

    delegate: Buddy {
        width: buddylist.width
        height: 60
        name: buddyName
        show: showText
        status: statusText
        avatar: avatarSrc
        maxWidth: buddyList.maxWidth
        onClicked: {
            buddylist.selected(jid, name)
        }
    }

    header: Component {
        Buddy {
            width: buddylist.width
            maxWidth: buddyList.maxWidth
            height: 60
        }
    }

    model: ListModel{}

    section.property: "group"
    section.criteria: ViewSection.FullString
    section.labelPositioning: ViewSection.InlineLabels | ViewSection.CurrentLabelAtStart
    section.delegate: Rectangle {
        width: buddylist.width
        height: 30
        color: "lightgray"
        Text {
            anchors.verticalCenter: parent.verticalCenter
            anchors.left: parent.left
            anchors.leftMargin: 5
            text: section
            font.bold: true
        }
    }

    function setUser(u) {
        headerItem.name = u.name.length === 0 ? u.jid : u.name
        headerItem.show = u.show
        headerItem.status = u.status
        headerItem.avatar = u.avatar.length === 0 ? "contact.svg": u.avatar
    }

    function addBuddy(buddy) {
        var pos = 0
        for (var i = 0; i < model.count; i++) {
            if (model.get(i).group === buddy.group) {
                pos = i
                break
            }
        }

        model.insert(pos, {"jid": buddy.jid,
                         "buddyName": buddy.name.length === 0 ? buddy.jid : buddy.name,
                         "group": buddy.group,
                         "avatarSrc": buddy.avatar.length === 0 ? "contact.svg": buddy.avatar,
                         "showText": buddy.show,
                         "statusText":buddy.status})

    }

    function appendBuddy(buddy) {
       // console.log(buddy.groups)
        model.append({"jid":buddy.jid,
                      "buddyName": buddy.name.length === 0 ? buddy.jid : buddy.name,
                      "group": buddy.group,
                      "avatarSrc": buddy.avatar.length === 0 ? "contact.svg": buddy.avatar,
                      "showText": buddy.show,
                      "statusText": buddy.status})
    }

    function removeBuddy(buddy) {
        for (var i = 0; i < model.count; i++) {
            if (model.get(i).jid === buddy.jid) {
                model.remove(i)
            }
        }
    }

    function clearBuddies() {
        buddylist.model.clear()
    }

    function updateBuddy(buddy) {
        var index = 0
        var o = {"jid":buddy.jid,
            "buddyName": buddy.name.length === 0 ? buddy.jid : buddy.name,
            "group": buddy.group,
            "avatarSrc": buddy.avatar.length === 0 ? "contact.svg": buddy.avatar,
            "showText": buddy.show,
            "statusText": buddy.status}

        for (var i = 0; i < model.count; i++) {
            var m = model.get(i)
            if (m.jid === buddy.jid && m.group === buddy.group) {
                model.remove(i)
                break
            }
        }
        for (var j = 0; j < model.count; j++) {
            var m = model.get(j)
            if (buddy.group === m.group) { // same group
                index = j
                if (buddy.show <= m.showText) {
                    break
                }
            }

        }

        model.insert(index, o)
    }
}
