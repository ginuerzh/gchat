// dialog
package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"
)

type Message struct {
	Jid    string
	Text   string
	Time   time.Time
	Unread bool
}

type Dialog struct {
	Jid      string
	Messages []*Message
}

func NewDialog(jid string, messages ...*Message) *Dialog {
	dialog := new(Dialog)
	dialog.Jid = jid
	dialog.Messages = append(dialog.Messages, messages...)
	return dialog
}

func (dialog *Dialog) Append(filename string, messages ...*Message) error {
	dialog.Messages = append(dialog.Messages, messages...)

	/*
		file, err := os.OpenFile(filename,
			os.O_WRONLY|os.O_APPEND|os.O_CREATE,
			os.ModePerm)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := file.Seek(0, os.SEEK_END); err != nil {
			return err
		}

		encoder := json.NewEncoder(file)
		for _, msg := range messages {
			if err := encoder.Encode(msg); err != nil {
				return err
			}
		}
	*/
	return dialog.Save(filename)
}

func (dialog *Dialog) Load(filename string) error {
	if len(filename) == 0 {
		return nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		msg := new(Message)
		if err := decoder.Decode(msg); err != nil {
			if err == io.EOF {
				break
			}
			log.Println(err)
			return err
		}
		dialog.Append("", msg)
	}
	return nil
}

func (dialog *Dialog) Save(filename string) error {
	if len(filename) == 0 {
		return nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for _, msg := range dialog.Messages {
		if err := encoder.Encode(msg); err != nil {
			return err
		}
	}

	return nil
}
