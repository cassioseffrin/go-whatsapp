package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"https://github.com/cassioseffrin/go-whatsapp"
	"https://github.com/cassioseffrin/go-whatsapp/binary/proto"
)

func main() {
 
	wac, err := whatsapp.NewConn(5 * time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
		return
	}

	err = login(wac)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
		return
	}

	<-time.After(3 * time.Second)

	previousMessage := "👍"
	quotedMessage := proto.Message{
		Conversation: &previousMessage,
	}

	ContextInfo := whatsapp.ContextInfo{
		QuotedMessage:   &quotedMessage,
		QuotedMessageID: "5496852447@s.whatsapp.net",
		Participant:     "Cassio", //Who sent the original message
	}

	msg := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: "555496852447@s.whatsapp.net",
		},
		ContextInfo: ContextInfo,
		Text:        "Enviada por Cassio!!!",
	}

	msgId, err := wac.Send(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error sending message: %v", err)
		os.Exit(1)
	} else {
		fmt.Println("Message Sent -> ID : " + msgId)
	}
}

// func main() {
// 	//create new WhatsApp connection
// 	wac, err := whatsapp.NewConn(2 * time.Second)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
// 		return
// 	}

// 	err = login(wac)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
// 		return
// 	}

// 	<-time.After(3 * time.Second)

// 	var numbers []string

// 	numbers = append(numbers, "5554996852447")

// 	for index, element := range numbers {
// 		print(index)

// 		msg := whatsapp.TextMessage{
// 			Info: whatsapp.MessageInfo{
// 				RemoteJid: element + "@s.whatsapp.net",
// 			},
// 			Text: fmt.Sprintf("teste %v", index),
// 		}

// 		// err = wac.Send(msg)
// 		msgId, err := wac.Send(msg)
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "error sending message: %v", err)
// 			os.Exit(1)
// 		} else {
// 			fmt.Println("Message Sent -> ID : " + msgId)
// 		}
// 	}

// }

func login(wac *whatsapp.Conn) error {
	//load saved session
	session, err := readSession()
	if err == nil {
		//restore session
		session, err = wac.RestoreWithSession(session)
		if err != nil {
			return fmt.Errorf("restoring failed: %v\n", err)
		}
	} else {
		//no saved session -> regular login
		qr := make(chan string)
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()
		session, err = wac.Login(qr)
		if err != nil {
			return fmt.Errorf("error during login: %v\n", err)
		}
	}

	//save session
	err = writeSession(session)
	if err != nil {
		return fmt.Errorf("error saving session: %v\n", err)
	}
	return nil
}

func readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func writeSession(session whatsapp.Session) error {
	file, err := os.Create(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}
