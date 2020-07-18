package main

import (
	"fmt"
	"log"

	irc "github.com/thoj/go-ircevent"
)

// THIS IS STILL WIP and I HAD PROBLEMS WORKING WITH THIS!!

//IRCSERVERNAMEPORT is the server used to report to
var IRCSERVERNAMEPORT = "hostname.example.com:6667"

//IRCBOTNICKNAME is nickname used by the account
var IRCBOTNICKNAME = "<NICK>"

//IRCSERVERPASSWORD is password used to log into the IRC server
var IRCSERVERPASSWORD = "<PASSWORD>"

// IRCUSERNAME is  the username of the account for the IRC server
var IRCUSERNAME = "<USER1>"

// IRCNAME is the name of the bot
var IRCNAME = "JARVIS"

// JOBCHANNEL is the irc channel name the bot will use
var JOBCHANNEL = fmt.Sprintf("#%sBOTNOTIFICATION", IRCNAME)

type ircclient struct {
	irc *irc.Connection
}

var ircMessageQueue = []string{}

func (cli *ircclient) prepareClient() {
	irccon := irc.IRC(IRCBOTNICKNAME, IRCUSERNAME)
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = false
	//irccon.Password = IRCSERVERPASSWORD
	irccon.AddCallback("001", func(e *irc.Event) {
		irccon.Join(JOBCHANNEL)
	})
	irccon.AddCallback("366", func(e *irc.Event) {
		//	irccon.Privmsg(JOBCHANNEL, "testworking")
	})
	cli.irc = irccon

}

func (cli *ircclient) connect() error {

	var err error = nil

	go func() {
		if err = cli.irc.Connect(IRCSERVERNAMEPORT); err != nil {
			return
		}
		cli.irc.Loop()
	}()

	return err
}

func (cli *ircclient) notifyMessage(message string) error {
	log.Printf("\nNotifyMessage called with %s\n", message)

	//ircMessageQueue = append(ircMessageQueue, message)
	// /target := fmt.Sprintf("%s", JOBCHANNEL)
	//cli.irc.Privmsg("#JARVISBOTNOTIFICATION", message)
	//cli.irc.Notice("#JARVISBOTNOTIFICATION", message) // sends a message to either a certain nick or a channel
	//cli.irc.Action("#JARVISBOTNOTIFICATION", message)

	return nil

}

func (cli *ircclient) disConnect() error {

	cli.irc.Disconnect()
	return nil

}
