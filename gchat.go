// gchat
package main

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
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
	"strings"
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

type Buddy struct {
	Jid    string
	Name   string
	Avatar string
	Group  string
	groups []string
	Show   int
	shows  map[string]int
	Status string
}

func NewBuddy(jid, name string, groups []string) *Buddy {
	if name == "" {
		name = jid
	}
	return &Buddy{
		Jid:    jid,
		Name:   name,
		groups: groups,
		Show:   ShowUnavail,
		shows:  make(map[string]int),
	}
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

type Chat struct {
	client *client.Client
	config *Config
	dir    string

	window *qml.Window

	buddies *BuddyList
	user    *Buddy
	dialogs map[string][]*Message

	pLock *sync.Mutex // mutex lock for presence
	cLock *sync.Mutex // mutex lock for card
}

func NewChat(dataDir string, config *Config) *Chat {
	if len(dataDir) == 0 {
		dataDir = os.Getenv("HOME") + "/.gchat"
	}

	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		panic(err)
		os.Exit(1)
	}

	return &Chat{
		config: config,
		dir:    dataDir,
		pLock:  &sync.Mutex{},
		cLock:  &sync.Mutex{},
	}
}

func (chat *Chat) Init() {
	chat.buddies = NewBuddyList()
	chat.dialogs = make(map[string][]*Message)
}

func (chat *Chat) ObjectByName(objectName string) qml.Object {
	return chat.window.Root().ObjectByName(objectName)
}

func (chat *Chat) LoadConfig() {
	addr := "talk.google.com"
	port := "443"

	serverAddr := chat.ObjectByName("serverAddr").String("text")
	serverPort := chat.ObjectByName("serverPort").String("text")

	if len(serverAddr) > 0 {
		addr = serverAddr
	}
	if len(serverPort) > 0 {
		port = serverPort
	}
	chat.config.Server = addr + ":" + port

	chat.config.Resource = chat.ObjectByName("resource").String("text")
	chat.config.UseOldTLS = chat.ObjectByName("sslSwitch").Bool("checked")
	chat.config.EnableProxy = chat.ObjectByName("proxySwitch").Bool("checked")
	if chat.config.EnableProxy {
		chat.config.UseSysProxy = chat.ObjectByName("sysProxySwitch").Bool("checked")
		if chat.config.UseSysProxy {
			proxy := os.Getenv("HTTP_PROXY")
			if proxy == "" {
				proxy = os.Getenv("http_proxy")
			}
			if proxy != "" {
				url, err := url.Parse(proxy)
				if err == nil {
					chat.config.Proxy = url.Host
				}
			}
		} else {
			proxyServerAddr := chat.ObjectByName("proxyServerAddr").String("text")
			proxyServerPort := chat.ObjectByName("proxyServerPort").String("text")

			if len(proxyServerAddr) > 0 && len(proxyServerPort) > 0 {
				chat.config.Proxy = proxyServerAddr + ":" + proxyServerPort
			}
		}
	}
	chat.config.Username = chat.ObjectByName("usernameInput").String("text")
	chat.config.Password = chat.ObjectByName("passwordInput").String("text")
	ciphertext, err := Encrypt(chat.config.Password)
	if err == nil {
		chat.config.Password = fmt.Sprintf("%x", ciphertext)
	}
	chat.config.AutoLogin = chat.ObjectByName("autoLogin").Bool("checked")

	if err := chat.config.Save(chat.dir + "/chat.conf"); err != nil {
		log.Println(err)
	}

}

func (chat *Chat) restoreConfig() {
	if len(chat.config.Server) > 0 {
		a := strings.SplitN(chat.config.Server, ":", 2)
		chat.ObjectByName("serverAddr").Set("text", a[0])
		if len(a) != 2 {
			chat.ObjectByName("serverPort").Set("text", "5222")
		} else {
			chat.ObjectByName("serverPort").Set("text", a[1])
		}
	}
	chat.ObjectByName("resource").Set("text", chat.config.Resource)
	chat.ObjectByName("sslSwitch").Set("checked", chat.config.UseOldTLS)

	chat.ObjectByName("proxySwitch").Set("checked", chat.config.EnableProxy)
	chat.ObjectByName("sysProxySwitch").Set("checked", chat.config.UseSysProxy)

	if len(chat.config.Proxy) > 0 {
		a := strings.SplitN(chat.config.Proxy, ":", 2)
		chat.ObjectByName("proxyServerAddr").Set("text", a[0])
		if len(a) != 2 {
			chat.ObjectByName("proxyServerPort").Set("text", "80")
		} else {
			chat.ObjectByName("proxyServerPort").Set("text", a[1])
		}
	}

	chat.ObjectByName("usernameInput").Set("text", chat.config.Username)

	password := chat.config.Password
	if plaintext, err := Decrypt(chat.config.Password); err == nil {
		password = string(plaintext)
	}
	chat.ObjectByName("passwordInput").Set("text", password)

	chat.ObjectByName("autoLogin").Set("checked", chat.config.AutoLogin)

	log.Println("server:", chat.config.Server,
		"proxy:", chat.config.Proxy,
		"username:", chat.config.Username,
		"password:", password,
		"oldTLS:", chat.config.UseOldTLS)
}

func (chat *Chat) addBubble(jid string, bubble *Message) {
	chat.dialogs[jid] = append(chat.dialogs[jid], bubble)

	dialogView := chat.ObjectByName("dialogView")
	if dialogView.String("jid") == jid {
		dialogView.Call("addBubble", bubble)
	}
}

func (chat *Chat) addMessage(msg *Message) {
	chat.ObjectByName("messageView").Call("addMessage", msg)
}

func (chat *Chat) Run() error {
	qml.Init(nil)
	engine := qml.NewEngine()
	component, err := engine.LoadFile("gchat.qml")
	if err != nil {
		return err
	}
	window := component.CreateWindow(nil)
	chat.window = window

	chat.restoreConfig()

	chat.ObjectByName("loginView").On("login", func(username, password string, status string) {
		if len(username) == 0 || len(password) == 0 {
			return
		}

		chat.login(username, password, status)
	})

	chat.ObjectByName("userTabs").On("logout", func() {
		chat.client.Close()
	})

	dialogView := chat.ObjectByName("dialogView")
	dialogView.On("loaded", func(jid string) {
		for _, bubble := range chat.dialogs[xmpp.ToJID(jid).Bare()] {
			dialogView.Call("addBubble", bubble)
		}
	})

	chat.ObjectByName("sendConfirm").On("sended", func(jid, text string) {
		chat.client.Send(xmpp.NewMessage("chat", jid, text, ""))
		chat.addBubble(jid, &Message{
			Jid:    chat.user.Jid,
			Text:   text,
			Time:   time.Now().Format("15:04"),
			Avatar: chat.user.Avatar,
		})
	})

	// handle Auto login
	if chat.config.AutoLogin {
		password := chat.config.Password
		if plaintext, err := Decrypt(chat.config.Password); err == nil {
			password = string(plaintext)
		}
		chat.login(chat.config.Username, password, "chat")
	}

	window.Show()
	window.Wait()

	return nil
}

func (chat *Chat) login(username, password string, status string) {
	chat.LoadConfig()

	cli := client.NewClient(chat.config.Server, username, password,
		&client.Options{
			Debug:     chat.config.EnableDebug,
			NoTLS:     !chat.config.UseOldTLS,
			Proxy:     chat.config.Proxy,
			Resource:  chat.config.Resource,
			TlsConfig: &tls.Config{InsecureSkipVerify: true}})

	chat.client = cli

	cli.OnLogined(func(err error) {
		if err != nil {
			fmt.Println("login:", err)
			chat.ObjectByName("loginPage").Call("logined", false, "", err.Error())
			return
		}

		chat.Init()
		chat.user = NewBuddy(cli.Jid.Bare(), "", nil)
		chat.user.Show = showPriv(status)

		chat.ObjectByName("buddyView").Call("setUser", chat.user)
		chat.ObjectByName("loginPage").Call("logined", true, chat.user.Name, "")

		cli.Send(xmpp.NewIQ("get", client.GenId(), chat.client.Jid.Domain(), &xep.DiscoInfoQuery{}))
		cli.Send(xmpp.NewIQ("get", client.GenId(), chat.client.Jid.Domain(), &xep.DiscoItemsQuery{}))

		cli.Send(xmpp.NewIQ("get", client.GenId(), "", &core.RosterQuery{}))
		cli.Send(xmpp.NewIQ("get", client.GenId(), "", &xep.VCard{}))
	})

	cli.HandleFunc(xmpp.NSRoster+" query", func(header *core.StanzaHeader, e xmpp.Element) {
		if header.Types != "result" {
			return
		}
		//fmt.Println(e)
		for _, item := range e.(*core.RosterQuery).Items {
			if item.Jid == cli.Jid.Bare() {
				continue
			}

			buddy := NewBuddy(item.Jid, item.Name, item.Group)
			chat.buddies.Add(buddy)
		}

		buddyView := chat.ObjectByName("buddyView")
		for group, buddies := range chat.buddies.groups {
			for _, buddy := range buddies {
				buddy.Group = group
				buddyView.Call("addBuddy", buddy)
			}
		}

		cli.Send(xmpp.NewStanza("presence"))
	})

	cli.HandleFunc(xmpp.NSDiscoItems+" query", func(header *core.StanzaHeader, e xmpp.Element) {
		fmt.Println(e)
	})
	cli.HandleFunc(xmpp.NSDiscoInfo+" query", func(header *core.StanzaHeader, e xmpp.Element) {
		fmt.Println(e)
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
				Name:   chat.buddies.Buddy(header.From).Name,
				Text:   body,
				Time:   time.Now().Format("15:04"),
				Avatar: chat.buddies.Buddy(header.From).Avatar,
				Unread: true,
			}
			if chat.ObjectByName("dialogView").Bool("show") {
				msg.Unread = false
			}
			chat.addBubble(msg.Jid, msg)
			chat.addMessage(msg)
		}
	})

	cli.HandleFunc(xmpp.NSClient+" presence", func(header *core.StanzaHeader, e xmpp.Element) {
		buddy := chat.buddies.Buddy(header.From)
		if buddy == nil {
			return
		}
		show := 0

		for _, e := range e.(*xmpp.Stanza).Elements {
			switch e.FullName() {
			case xmpp.NSClient + " show":
				show = showPriv(e.(*core.PresenceShow).Show)
			case xmpp.NSClient + " status":
				buddy.Status = e.(*core.PresenceStatus).Status
			case xmpp.NSVcardUpdate + " x":
				avatar := e.(*xep.VCardUpdate).Photo
				if len(avatar) == 0 {
					continue
				}
				path := chat.dir + "/" + cli.Jid.Bare() + "/avatar"
				if matchs, _ := filepath.Glob(path + "/" + avatar + ".*"); len(matchs) > 0 {
					//fmt.Println("avatar exists", buddy.Avatar)
					buddy.Avatar = matchs[0]
					continue
				}
				cli.Send(xmpp.NewIQ("get", client.GenId(), xmpp.ToJID(header.From).Bare(), &xep.VCard{}))
			}
		}
		if show == 0 {
			show = ShowChat // default is chat
			if header.Types == "unavailable" {
				show = ShowUnavail
			}
		}

		buddy.shows[header.From] = show
		buddy.Show = ShowUnavail
		for _, v := range buddy.shows {
			if v < buddy.Show {
				buddy.Show = v
			}
		}

		chat.pLock.Lock()
		chat.ObjectByName("buddyView").Call("updateBuddy", buddy)
		chat.pLock.Unlock()
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

		path := chat.dir + "/" + cli.Jid.Bare() + "/avatar"
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

		if buddy := chat.buddies.Buddy(header.From); buddy != nil {
			buddy.Avatar = filename

			chat.cLock.Lock()
			chat.ObjectByName("buddyView").Call("updateBuddy", buddy)
			chat.cLock.Unlock()
		}
		if len(header.From) == 0 {
			chat.user.Avatar = filename
			if len(card.FName) > 0 {
				chat.user.Name = card.FName
			}
			chat.ObjectByName("buddyView").Call("setUser", chat.user)
		}
	})
	go cli.Run()
}

func showPriv(s string) int {
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
