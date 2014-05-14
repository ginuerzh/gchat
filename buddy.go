// buddy
package main

import (
	xmpp "github.com/ginuerzh/goxmpp"
	"sync"
)

type Buddy struct {
	Jid          string
	Name         string
	Avatar       string
	avatarHash   string
	Group        string // used by qml
	Groups       []string
	Show         int // used by qml
	shows        map[string]int
	Subscription string
	Status       string
	Dialog       *Dialog
}

func NewBuddy(jid, name string, groups []string, subscription string) *Buddy {
	return &Buddy{
		Jid:          jid,
		Name:         name,
		Groups:       groups,
		Subscription: subscription,
		Show:         ShowUnavail,
		shows:        make(map[string]int),
		Dialog:       NewDialog(jid),
	}
}

type BuddyList struct {
	buddies map[string]*Buddy
	Groups  map[string][]*Buddy
	lock    *sync.RWMutex
}

func NewBuddyList() *BuddyList {
	return &BuddyList{
		buddies: make(map[string]*Buddy),
		Groups:  make(map[string][]*Buddy),
		lock:    &sync.RWMutex{},
	}
}

func (l *BuddyList) Add(buddy *Buddy) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.buddies[xmpp.ToJID(buddy.Jid).Bare()] = buddy
	if len(buddy.Groups) == 0 {
		l.Groups["Buddies"] = append(l.Groups["Buddies"], buddy)
		return
	}
	for _, group := range buddy.Groups {
		l.Groups[group] = append(l.Groups[group], buddy)
	}
}

func (l *BuddyList) Remove(jid string) *Buddy {
	l.lock.Lock()
	defer l.lock.Unlock()

	buddy := l.buddies[xmpp.ToJID(jid).Bare()]
	if buddy == nil {
		return nil
	}

	delete(l.buddies, jid)

	delFromGroup := func(group string) {
		buddies := l.Groups[group]
		for i, buddy := range buddies {
			if buddy.Jid == jid {
				l.Groups[group] = append(buddies[:i], buddies[i+1:]...)
			}
		}
	}

	if len(buddy.Groups) == 0 {
		delFromGroup("Buddies")
	} else {
		for _, group := range buddy.Groups {
			delFromGroup(group)
		}
	}

	return buddy
}

func (l *BuddyList) Buddy(jid string) *Buddy {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.buddies[xmpp.ToJID(jid).Bare()]
}
