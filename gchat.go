// gchat
package main

import (
	//"crypto/sha1"
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

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Chat struct {
	client *client.Client
	config *Config
	dir    string

	features []string // server features

	engine *qml.Engine
	window *qml.Window

	buddies *BuddyList
	user    *Buddy
	//dialogs map[string][]*Message

	lock  *sync.Mutex
	bLock *sync.Mutex
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
		lock:   &sync.Mutex{},
		bLock:  &sync.Mutex{},
	}
}

func (chat *Chat) Init(user *Buddy) {
	chat.user = user
	chat.buddies = NewBuddyList()

	if err := os.MkdirAll(chat.MessagePath(), os.ModePerm); err != nil {
		fmt.Println(err)
	}
	if err := os.MkdirAll(chat.AvatarPath(), os.ModePerm); err != nil {
		fmt.Println(err)
	}
}

func (chat *Chat) MessagePath() string {
	return chat.dir + "/" + chat.user.Jid + "/messages"
}

func (chat *Chat) AvatarPath() string {
	return chat.dir + "/" + chat.user.Jid + "/avatars"
}

func (chat *Chat) AvatarFile(jid string) (name string, hash string) {
	avatars, err := filepath.Glob(chat.AvatarPath() + "/" + jid + "*")
	if err != nil {
		log.Println(err)
	}
	if len(avatars) == 0 {
		return
	}
	avatar := avatars[0]
	a := strings.SplitN(filepath.Base(avatar), " ", 2)
	if len(a) != 2 {
		return
	}
	a = strings.SplitN(a[1], ".", 2)
	if len(a) != 2 {
		return
	}

	return avatar, a[0]
}

func (chat *Chat) MessageFile(jid string) string {
	return chat.MessagePath() + "/" + xmpp.ToJID(jid).Bare() + ".json"
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
	chat.config.NoTLS = !chat.ObjectByName("sslSwitch").Bool("checked")
	chat.config.Proxy = ""
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
	log.Println("server:", chat.config.Server,
		"proxy:", chat.config.Proxy,
		"username:", chat.config.Username,
		"notls:", chat.config.NoTLS)
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
	chat.ObjectByName("sslSwitch").Set("checked", !chat.config.NoTLS)

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
		"notls:", chat.config.NoTLS)
}

func (chat *Chat) addBubble(jid string, bubble *Message, logToFile bool) {
	buddy := chat.buddies.Buddy(jid)

	if logToFile {
		filename := chat.MessageFile(xmpp.ToJID(jid).Bare())
		if err := buddy.Dialog.Append(filename, bubble); err != nil {
			log.Println(err)
		}
	}

	if bubble.Jid == chat.user.Jid {
		buddy = chat.user
	}

	dialogView := chat.ObjectByName("dialogView")
	if dialogView.String("jid") == jid {
		dialogView.Call("addBubble", buddy, bubble)
	}
}

func (chat *Chat) addMessage(buddy *Buddy, msg *Message) {
	if buddy == nil || msg == nil {
		return
	}
	chat.ObjectByName("messageView").Call("addMessage", buddy, msg)
}

func (chat *Chat) Run() error {
	qml.Init(nil)
	chat.engine = qml.NewEngine()
	component, err := chat.engine.LoadFile("gchat.qml")
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

	msgView := chat.ObjectByName("messageView")
	msgView.On("accepted", func(jid string) {
		chat.client.Send(xmpp.NewPresence("subscribed", "", jid))
		chat.client.Send(xmpp.NewPresence("subscribe", "", jid))
	})
	msgView.On("declined", func(jid string) {
		chat.client.Send(xmpp.NewPresence("unsubscribed", "", jid))
	})

	dialogView := chat.ObjectByName("dialogView")
	dialogView.On("loaded", func(jid string) {
		buddy := chat.buddies.Buddy(jid)
		for _, bubble := range buddy.Dialog.Messages {
			chat.addBubble(jid, bubble, false)
		}
	})

	chat.ObjectByName("sendConfirm").On("sended", func(jid, text string) {
		chat.client.Send(xmpp.NewMessage("chat", jid, text, ""))
		msg := &Message{
			Jid:  chat.user.Jid,
			Text: text,
			Time: time.Now(),
		}
		chat.addBubble(jid, msg, true)
		chat.addMessage(chat.buddies.Buddy(jid), msg)
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
			NoTLS:     chat.config.NoTLS,
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

		user := NewBuddy(cli.Jid.Bare(), "", nil, "")
		user.Show = showPriv(status)
		chat.Init(user)
		chat.engine.Context().SetVar("loginUser", user)

		chat.ObjectByName("buddyView").Call("setUser", chat.user)
		chat.ObjectByName("loginPage").Call("logined", true, chat.user.Name, "")

		//cli.Send(xmpp.NewIQ("get", client.GenId(), chat.client.Jid.Domain(), &xep.DiscoInfoQuery{}))
		//cli.Send(xmpp.NewIQ("get", client.GenId(), chat.client.Jid.Domain(), &xep.DiscoItemsQuery{}))

		cli.Send(xmpp.NewIQ("get", client.GenId(), "", &core.RosterQuery{}))
		cli.Send(xmpp.NewIQ("get", client.GenId(), "", &xep.VCard{}))
	})

	cli.OnError(func(err error) {
		log.Println(err) // TODO error handling
	})

	// ping response
	cli.HandleFunc(xmpp.NSPing+" ping", func(header *core.StanzaHeader, e xmpp.Element) {
		cli.Send(xmpp.NewIQ("result", header.Ids, header.From, nil))
	})

	// roster
	cli.HandleFunc(xmpp.NSRoster+" query", func(header *core.StanzaHeader, e xmpp.Element) {
		//fmt.Println(e)
		if header.Types == "error" {
			return
		}

		removeBuddy := func(buddy *Buddy) {
			buddy = chat.buddies.Remove(buddy.Jid)
			if buddy == nil {
				return
			}

			chat.lock.Lock()
			chat.ObjectByName("buddyView").Call("removeBuddy", buddy)
			chat.lock.Unlock()
			chat.ObjectByName("messageView").Call("removeMessage", buddy.Jid)
		}

		initBuddy := func(buddy *Buddy) {
			buddy.Avatar, buddy.avatarHash = chat.AvatarFile(buddy.Jid)
			buddy.Dialog.Load(chat.MessageFile(buddy.Jid))

			if len(buddy.Dialog.Messages) > 0 {
				chat.addMessage(buddy, buddy.Dialog.Messages[len(buddy.Dialog.Messages)-1]) // show the last message
			}
		}

		if header.Types == "set" {
			for _, item := range e.(*core.RosterQuery).Items {
				switch item.Subscription {
				case "remove":
					removeBuddy(NewBuddy(item.Jid, item.Name, item.Group, item.Subscription))
				case "none":
					break
				default:
					chat.bLock.Lock()
					buddy := chat.buddies.Buddy(item.Jid)
					if buddy != nil {
						buddy.Name = item.Name
						buddy.Groups = item.Group
						buddy.Subscription = item.Subscription

						chat.lock.Lock()
						chat.ObjectByName("buddyView").Call("updateBuddy", buddy)
						chat.lock.Unlock()

						chat.bLock.Unlock()
						break
					}

					buddy = NewBuddy(item.Jid, item.Name, item.Group, item.Subscription)

					initBuddy(buddy)
					chat.buddies.Add(buddy)
					buddy.Group = "Buddies"
					chat.lock.Lock()
					chat.ObjectByName("buddyView").Call("addBuddy", buddy)
					chat.lock.Unlock()

					chat.bLock.Unlock()
				}
			}
			return
		}

		for _, item := range e.(*core.RosterQuery).Items {
			if item.Jid == cli.Jid.Bare() || chat.buddies.Buddy(item.Jid) != nil {
				continue
			}
			buddy := NewBuddy(item.Jid, item.Name, item.Group, item.Subscription)
			initBuddy(buddy)
			chat.buddies.Add(buddy)
		}

		buddyView := chat.ObjectByName("buddyView")
		for group, buddies := range chat.buddies.Groups {
			for _, buddy := range buddies {
				buddy.Group = group
				buddyView.Call("appendBuddy", buddy)
			}
		}

		cli.Send(xmpp.NewStanza("presence"))
	})

	cli.HandleFunc(xmpp.NSDiscoItems+" query", func(header *core.StanzaHeader, e xmpp.Element) {
		//fmt.Println(e)
		if header.Types == "error" {
			return
		}

		for _, item := range e.(*xep.DiscoItemsQuery).Items {
			cli.Send(xmpp.NewIQ("get", client.GenId(), item.Jid, new(xep.DiscoInfoQuery)))
		}
	})
	cli.HandleFunc(xmpp.NSDiscoInfo+" query", func(header *core.StanzaHeader, e xmpp.Element) {
		//fmt.Println(e)

		if header.Types == "error" {
			return
		}
		query := e.(*xep.DiscoInfoQuery)

		for _, id := range query.Identities {
			// See http://xmpp.org/registrar/disco-categories.html
			switch id.Category + " " + id.Type {
			case "server im":
				for _, feature := range query.Features {
					chat.features = append(chat.features, feature.Var)
				}
			case "conference text":
				log.Println("find conference", id.Name, header.From)
				iq, err := cli.SendIQ(xmpp.NewIQ("get", client.GenId(), header.From, new(xep.DiscoItemsQuery)))
				if err != nil {
					log.Println(err)
					break
				}
				if err = iq.Error(); err != nil {
					log.Println(err)
					break
				}

				query := iq.E()[0].(*xep.DiscoItemsQuery)

				for _, item := range query.Items {
					log.Println("find chat room", item.Jid, item.Name)
				}
			case "directory chatroom":

			case "pubsub service":
			case "proxy bytestreams":

			}
		}
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
				Text:   body,
				Time:   time.Now(),
				Unread: true,
			}
			if chat.ObjectByName("dialogView").Bool("show") {
				msg.Unread = false
			}
			chat.addBubble(msg.Jid, msg, true)
			chat.addMessage(chat.buddies.Buddy(msg.Jid), msg)
		}
	})

	cli.HandleFunc(xmpp.NSClient+" presence", func(header *core.StanzaHeader, e xmpp.Element) {
		switch header.Types {
		case "subscribe":
			msg := &Message{
				Jid:  xmpp.ToJID(header.From).Bare(),
				Text: "Subscribe",
				Time: time.Now(),
			}
			chat.ObjectByName("messageView").Call("addSubscribe", msg)
		case "subscribed", "unsubscribe":
			return
		default:
			break
		}

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
				hash := e.(*xep.VCardUpdate).Photo
				if len(hash) == 0 || buddy.avatarHash == hash {
					continue
				}

				buddy.avatarHash = hash
				cli.Send(xmpp.NewIQ("get", client.GenId(), buddy.Jid, &xep.VCard{}))
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

		chat.lock.Lock()
		chat.ObjectByName("buddyView").Call("updateBuddy", buddy)
		chat.lock.Unlock()
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

		buddy := chat.buddies.Buddy(header.From)
		if len(header.From) == 0 {
			buddy = chat.user
		}

		filename := chat.AvatarPath() + "/" + buddy.Jid + " " + buddy.avatarHash + ".jpg"
		if err := ioutil.WriteFile(filename, data, os.ModePerm); err != nil {
			buddy.avatarHash = ""
			fmt.Println(err)
			return
		}

		buddy.Avatar = filename
		if len(card.FName) > 0 {
			buddy.Name = card.FName
		}

		if len(header.From) == 0 {
			chat.ObjectByName("buddyView").Call("setUser", buddy)
		} else {
			chat.lock.Lock()
			chat.ObjectByName("buddyView").Call("updateBuddy", buddy)
			chat.lock.Unlock()
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
