import QtQuick 2.0

ListView {
    objectName: "buddies"
    id: buddylist
    spacing: 1
    signal selected(string jid, string name)

    delegate: Buddy {
        width: buddylist.width
        height: 64
        name: buddyName
        show: showText
        status: statusText
        avatar: avatarSrc
        onClicked: {
            buddylist.selected(jid, name)
        }
    }

    header: Component {
        Buddy {
            width: buddylist.width
            height: 64
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
        headerItem.name = u.name
        headerItem.show = u.show
        headerItem.status = u.status
        headerItem.avatar = u.avatar.length === 0 ? "contact.svg": u.avatar
    }

    function addBuddy(buddy) {
        model.append({"jid":buddy.jid,
                      "buddyName": buddy.name,
                      "group": buddy.group,
                      "avatarSrc": buddy.avatar.length === 0 ? "contact.svg": buddy.avatar,
                      "showText": buddy.show,
                      "statusText": buddy.status})
    }

    function updateBuddy(buddy) {
        var index = 0
        var o = {"jid":buddy.jid,
            "buddyName": buddy.name,
            "group": buddy.group,
            "avatarSrc": buddy.avatar.length === 0 ? "contact.svg": buddy.avatar,
            "showText": buddy.show,
            "statusText": buddy.status}

        for (var i = 0; i < model.count; i++) {
            if (model.get(i).jid === buddy.jid) {
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
