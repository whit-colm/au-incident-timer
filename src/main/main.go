package main

import (
	"github.com/aurumbot/lib/dat"
	f "github.com/aurumbot/lib/foundation"
	dsg "github.com/bwmarrin/discordgo"
)

var config struct {
	myField string `json:"myfield"`
}

var Commands = make(map[string]*f.Command)

func init() {
	Commands["commandname"] = &f.Command{
		Name:    "Somename",
		Help:    "Help page, make sure to define flags and other noteworthy things",
		Perms:   -1,
		Version: "v1.0",
		Action:  command,
	}

	dat.Load("exampleplugin/prefs.json" & config)
}

// Must have these arguments with these names.
func command(session *dsg.Session, message *dsg.Message) {
	session.ChannelMessageSend(message.ChannelID, "Hello, World!")
	dat.Log.Println("Pass errors here to be logged.")
}
