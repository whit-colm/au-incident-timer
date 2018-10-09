package main

import (
	"github.com/aurumbot/flags"
	"github.com/aurumbot/lib/dat"
	f "github.com/aurumbot/lib/foundation"
	dsg "github.com/bwmarrin/discordgo"
)

var config struct {
	Guild map[string]struct {
		Topic []struct {
			Name     string        `json:"name"`
			Triggers []string      `json:"triggers"`
			Since    time.Time     `json:"since"`
			User     string        `json:"user"`
			Message  string        `json:"message"`
			Streak   time.Duration `json:"streak"`
		} `json:"topic"`
	} `json:"guild"`
}

var Commands = make(map[string]*f.Command)

func init() {
	Commands["bpconfig"] = &f.Command{
		Name: "The Bad Post Configuration Tool",
		Help: `bpconfig is the modificaation tool for the bad post timer. It is used to define
topics and triggers (as regular expressions). If you have retained vestages of
your sanity but want to learn regex, try https://regexr.com.
# Arguments
**-new <topic name>** : Create a new topic with the topic name <topic name>
**-rm <topic name>**  : Remove the <topic name>. Dangerous!
**-ls** : List all topics
**<topic name> <--addtrigger|-a> <//regex//>** : Add a trigger to <topic name> fulfilling //regex// (must be in //double slashes//)
**<topic name> <--rmtrigger|-rm> <item number>** : remove the <item number> trigger from <topic name>
**<topic name> <-ls>** : List all topic triggers
# Usage:
` + f.Config.Prefix + `bpconfig -new foo
` + f.Config.Prefix + `bpconfig foo -a //^[A-z].+$//`,
		Perms:   dsg.PermissionManageWebhooks,
		Version: "v1.0",
		Action:  bpconfig,
	}

	dat.Load("incident-counter/config.json", &config)
	if config.Guild == nil {
		config.Guild = make(map[string]struct {
			Topic []struct {
				Name     string        `json:"name"`
				Triggers []string      `json:"triggers"`
				Since    time.Time     `json:"since"`
				User     string        `json:"user"`
				Message  string        `json:"message"`
				Streak   time.Duration `json:"streak"`
			} `json:"topic"`
		})
	}

	f.Session.AddHandler(incidentHandler)
}

func incidentHandler(s *dsg.Session, m *dsg.MessageCreate) {
	if message.Author.Bot {
		return
	}
	var guildID string
	if channel, err := session.Channel(m.Message.ChannelID); err != nil {
		dat.Log.Println(err)
		return
	} else {
		guildID = channel.GuildID
	}
	// Local topics reduces verbosity by grabbing the config.Guild[guildID] value.
	localTopics := config.Guild[guildID].Topics
	for li := range localTopics {
		for i := range localTopics[li].Triggers {
			match, err := regexp.Compile(localTopics[li].Triggers[i])
			if err != nil {
				dat.Log.Println(err)
				return
			}

		}
	}
}

func bpconfig(session *dsg.Session, message *dsg.Message) {
	var guildID string
	if channel, err := session.Channel(message.ChannelID); err != nil {
		dat.Log.Println(err)
		return
	} else {
		guildID = channel.GuildID
	}
	// Local topics reduces verbosity by grabbing the config.Guild[guildID] value.
	localTopics := config.Guild[guildID].Topics

	flgs := flags.Parse(message.Content)
	for i := range flgs {

	}
}

func removeTopicConfirm(s *dsg.Session, m *dsg.MessageCreate) {

}
