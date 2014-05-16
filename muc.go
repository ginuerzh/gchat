// muc
package main

import (
	xmpp "github.com/ginuerzh/goxmpp"
	"github.com/ginuerzh/goxmpp/xep"
	"sync"
)

type MUC struct {
	groups map[string]*GroupChat
	locker *sync.RWMutex
}

func NewMUC() *MUC {
	return &MUC{
		groups: make(map[string]*GroupChat),
		locker: new(sync.RWMutex),
	}
}

func (muc *MUC) Add(group *GroupChat) {
	muc.locker.Lock()
	defer muc.locker.Unlock()

	muc.groups[group.Room.Jid] = group
}

func (muc *MUC) Remove(jid string) *GroupChat {
	muc.locker.Lock()
	defer muc.locker.Unlock()

	jid = xmpp.ToJID(jid).Bare()
	room := muc.groups[jid]
	delete(muc.groups, jid)

	return room
}

func (muc *MUC) Group(jid string) *GroupChat {
	muc.locker.RLock()
	defer muc.locker.RUnlock()

	return muc.groups[xmpp.ToJID(jid).Bare()]
}

type GroupChat struct {
	Room   *xep.ChatRoom
	Dialog *Dialog
}

func NewGroupChat(room *xep.ChatRoom, dialog *Dialog) *GroupChat {
	return &GroupChat{
		Room:   room,
		Dialog: dialog,
	}
}
