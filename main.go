// main
package main

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	xmpp "github.com/ginuerzh/goxmpp"
	"github.com/ginuerzh/goxmpp/client"
	"github.com/ginuerzh/goxmpp/core"
	"github.com/ginuerzh/goxmpp/xep"
	"gopkg.in/qml.v0"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	_ = iota
	ShowChat
	ShowDnd
	ShowAway
	ShowXa
	ShowUnavail
)

func show(s string) int {
	i := ShowUnavail
	switch s {
	case "chat":
		i = ShowChat
	case "dnd":
		i = ShowDnd
	case "away":
		i = ShowAway
	case "xa":
		i = ShowXa
	}
	return i
}

type Buddy struct {
	Jid    string
	Name   string
	Avatar string
	Group  string
	groups []string
	Show   int
	Status string
}

type BuddyList struct {
	buddies map[string]*Buddy
	groups  map[string][]*Buddy
}

func NewBuddyList() *BuddyList {
	return &BuddyList{
		buddies: make(map[string]*Buddy),
		groups:  make(map[string][]*Buddy),
	}
}

func (l *BuddyList) Add(buddy *Buddy) {
	l.buddies[xmpp.ToJID(buddy.Jid).Bare()] = buddy
	if len(buddy.groups) == 0 {
		l.groups["Buddies"] = append(l.groups["Buddies"], buddy)
		return
	}
	for _, group := range buddy.groups {
		l.groups[group] = append(l.groups[group], buddy)
	}
}

func (l *BuddyList) Buddy(jid string) *Buddy {
	return l.buddies[xmpp.ToJID(jid).Bare()]
}

type Message struct {
	Jid    string
	Name   string
	Text   string
	Time   string
	Avatar string
	Unread bool
}

var (
	dataPath = os.Getenv("HOME") + "/.gchat"

	xmppClient *client.Client

	window      *qml.Window
	messageView qml.Object
	dialogView  qml.Object

	buddyList *BuddyList
	dialogs   map[string][]*Message

	user         *Buddy
	presenceLock = &sync.Mutex{}
	vCardLock    = &sync.Mutex{}
)

func addBubble(jid string, bubble *Message) {
	dialogs[jid] = append(dialogs[jid], bubble)

	if dialogView.String("jid") == jid {
		dialogView.Call("addBubble", bubble)
	}
}

func addMessage(msg *Message) {
	messageView.Call("addMessage", msg)
}

var (
	flagServer  = flag.String("server", "talk.google.com:443", "xmpp server")
	flagProxy   = flag.String("proxy", "", "proxy server")
	useSysProxy = flag.Bool("sproxy", false, "Use system proxy")
	oldTLS      = flag.Bool("tls", false, "use old tls")
	enableDebug = flag.Bool("debug", false, "enable debug")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func initSettings() {
	*flagProxy = ""
	addr := "talk.google.com"
	port := "443"

	serverAddr := window.Root().ObjectByName("serverAddr").String("text")
	serverPort := window.Root().ObjectByName("serverPort").String("text")

	if len(serverAddr) > 0 {
		addr = serverAddr
	}
	if len(serverPort) > 0 {
		port = serverPort
	}
	*flagServer = addr + ":" + port

	*oldTLS = window.Root().ObjectByName("sslSwitch").Bool("checked")
	if window.Root().ObjectByName("proxySwitch").Bool("checked") {

		*useSysProxy = window.Root().ObjectByName("sysProxySwitch").Bool("checked")
		if *useSysProxy {
			proxy := os.Getenv("HTTP_PROXY")
			if proxy == "" {
				proxy = os.Getenv("http_proxy")
			}
			if proxy != "" {
				url, err := url.Parse(proxy)
				if err == nil {
					*flagProxy = url.Host
				}
			}
		} else {
			proxyServerAddr := window.Root().ObjectByName("proxyServerAddr").String("text")
			proxyServerPort := window.Root().ObjectByName("proxyServerPort").String("text")

			if len(proxyServerAddr) > 0 && len(proxyServerPort) > 0 {
				*flagProxy = proxyServerAddr + ":" + proxyServerPort
			}
		}
	}
	log.Println("server:", *flagServer, ", proxy:", *flagProxy, ", oldTLS:", *oldTLS)
}

func run() error {
	qml.Init(nil)
	engine := qml.NewEngine()

	component, err := engine.LoadFile("gchat.qml")
	if err != nil {
		return err
	}
	window = component.CreateWindow(nil)
	window.Show()

	window.Root().ObjectByName("loginView").On("login", func(username, password string, remember bool) {
		if len(username) == 0 || len(password) == 0 {
			return
		}
		initSettings()
		xmppClient = createClient(*flagServer, username, password, *flagProxy, *oldTLS, *enableDebug)
		go xmppClient.Run()
	})
	window.Root().ObjectByName("userTabs").On("logout", func() {
		xmppClient.Close()
	})

	messageView = window.Root().ObjectByName("messages")

	dialogView = window.Root().ObjectByName("dialog")
	dialogView.On("loaded", func(jid string) {
		for _, bubble := range dialogs[xmpp.ToJID(jid).Bare()] {
			dialogView.Call("addBubble", bubble)
		}
	})

	window.Root().ObjectByName("sendConfirm").On("sended", func(jid, text string) {
		xmppClient.Send(xmpp.NewMessage("chat", jid, text, ""))
		addBubble(jid, &Message{
			Jid:    user.Jid,
			Text:   text,
			Time:   time.Now().Format("15:04"),
			Avatar: user.Avatar,
		})
	})

	window.Wait()
	return nil
}

func createClient(server, username, password, proxy string, oldTLS, debug bool) *client.Client {
	cli := client.NewClient(server, username, password,
		&client.Options{
			Debug:     debug,
			NoTLS:     !oldTLS,
			Proxy:     proxy,
			TlsConfig: &tls.Config{InsecureSkipVerify: true}})

	cli.OnLogined(func(err error) {
		if err != nil {
			fmt.Println("login:", err)
			window.Root().ObjectByName("loginPage").Call("logined", false, "", err.Error())
			return
		}

		buddyList = NewBuddyList()
		dialogs = make(map[string][]*Message)

		user = &Buddy{Jid: cli.Jid.String(), Name: cli.Jid.Bare(), Show: show("chat")}
		window.Root().ObjectByName("buddies").Call("setUser", user)
		window.Root().ObjectByName("loginPage").Call("logined", true, user.Name, "")

		//cli.Send(xmpp.NewIQ("get", client.GenId(), "", &xep.DiscoItemsQuery{}))
		//cli.Send(xmpp.NewIQ("get", client.GenId(), "", &xep.DiscoInfoQuery{}))
		cli.Send(xmpp.NewIQ("get", client.GenId(), "", &core.RosterQuery{}))
		cli.Send(xmpp.NewIQ("get", client.GenId(), "", &xep.VCard{}))
	})

	cli.HandleFunc(xmpp.NSRoster+" query", func(header *core.StanzaHeader, e xmpp.Element) {
		if header.Types != "result" {
			return
		}
		//fmt.Println(e)
		for _, item := range e.(*core.RosterQuery).Items {
			if xmpp.ToJID(item.Jid).Bare() == xmppClient.Jid.Bare() {
				continue
			}
			buddy := &Buddy{Jid: item.Jid, Show: ShowUnavail}
			if len(item.Name) > 0 {
				buddy.Name = item.Name
			} else {
				buddy.Name = item.Jid
			}
			buddy.groups = item.Group
			buddyList.Add(buddy)
		}

		buddyView := window.Root().ObjectByName("buddies")
		for group, buddies := range buddyList.groups {
			for _, buddy := range buddies {
				buddy.Group = group
				buddyView.Call("addBuddy", buddy)
			}
		}

		cli.Send(xmpp.NewStanza("presence"))
	})

	cli.HandleFunc(xmpp.NSDiscoItems+" query", func(header *core.StanzaHeader, e xmpp.Element) {
		//fmt.Println(e)
	})
	cli.HandleFunc(xmpp.NSDiscoInfo+" query", func(header *core.StanzaHeader, e xmpp.Element) {
		//fmt.Println(e)
	})

	cli.HandleFunc(xmpp.NSClient+" message", func(header *core.StanzaHeader, e xmpp.Element) {
		msg := e.(*xmpp.Stanza)
		body := ""
		for _, e := range msg.E() {
			if e.Name() == "body" {
				body = e.(*core.MsgBody).Body
				break
			}
		}
		if len(body) > 0 {
			msg := &Message{
				Jid:    xmpp.ToJID(header.From).Bare(),
				Name:   buddyList.Buddy(header.From).Name,
				Text:   body,
				Time:   time.Now().Format("15:04"),
				Avatar: buddyList.Buddy(header.From).Avatar,
				Unread: true,
			}
			if dialogView.Bool("show") {
				msg.Unread = false
			}
			addBubble(msg.Jid, msg)
			addMessage(msg)
		}
	})

	cli.HandleFunc(xmpp.NSClient+" presence", func(header *core.StanzaHeader, e xmpp.Element) {
		buddy := buddyList.Buddy(header.From)
		if buddy == nil {
			return
		}

		for _, e := range e.(*xmpp.Stanza).Elements {
			switch e.FullName() {
			case xmpp.NSClient + " show":
				buddy.Show = show(e.(*core.PresenceShow).Show)
			case xmpp.NSClient + " status":
				buddy.Status = e.(*core.PresenceStatus).Status
			case xmpp.NSVcardUpdate + " x":
				avatar := e.(*xep.VCardUpdate).Photo
				if len(avatar) == 0 {
					continue
				}
				path := dataPath + "/" + xmppClient.Jid.Bare() + "/avatar"
				if matchs, _ := filepath.Glob(path + "/" + avatar + ".*"); len(matchs) > 0 {
					//fmt.Println("avatar exists", buddy.Avatar)
					buddy.Avatar = matchs[0]
					continue
				}
				cli.Send(xmpp.NewIQ("get", client.GenId(), xmpp.ToJID(header.From).Bare(), &xep.VCard{}))
			}
		}
		if buddy.Show == 0 {
			buddy.Show = ShowChat // default is chat
			if header.Types == "unavailable" {
				buddy.Show = ShowUnavail
			}
		}
		presenceLock.Lock()
		window.Root().ObjectByName("buddies").Call("updateBuddy", buddy)
		presenceLock.Unlock()
	})

	cli.HandleFunc(xmpp.NSVcardTemp+" vCard", func(header *core.StanzaHeader, e xmpp.Element) {
		card := e.(*xep.VCard)
		if card.Photo == nil {
			return
		}
		data, err := base64.StdEncoding.DecodeString(card.Photo.BinVal)
		if err != nil {
			fmt.Println(err)
			return
		}

		h := sha1.New()
		h.Write(data)
		hex := fmt.Sprintf("%x", h.Sum(nil))

		path := dataPath + "/" + xmppClient.Jid.Bare() + "/avatar"
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			fmt.Println(err)
			return
		}
		filename := path + "/" + hex + ".jpg"
		//fmt.Println(filename)
		if err := ioutil.WriteFile(filename, data, os.ModePerm); err != nil {
			fmt.Println(err)
			return
		}

		if buddy := buddyList.Buddy(header.From); buddy != nil {
			buddy.Avatar = filename

			vCardLock.Lock()
			window.Root().ObjectByName("buddies").Call("updateBuddy", buddy)
			vCardLock.Unlock()
		}
		if len(header.From) == 0 {
			user.Avatar = filename
			if len(card.FName) > 0 {
				user.Name = card.FName
			}
			window.Root().ObjectByName("buddies").Call("setUser", user)
		}
	})

	return cli
}
