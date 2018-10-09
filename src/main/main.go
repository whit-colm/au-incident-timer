package main

import (
	"github.com/aurumbot/flags"
	"github.com/aurumbot/lib/dat"
	f "github.com/aurumbot/lib/foundation"
	dsg "github.com/bwmarrin/discordgo"
	"time"
)

type topic struct {
	Name       string    `json:"name"`
	Triggers   []string  `json:"triggers"`
	Since      time.Time `json:"since"`
	LastUser   string    `json:"lastuser"`
	LastSet    time.Time `json:"lastset"`
	BreakUser  string    `json:"breakuser"`
	BreakTime  time.Time `json:"breaktime"`
	StreakUser string    `json:"streakuser"`
	StreakSet  time.Time `json:"streakset"`
	Message    string    `json:"message"`
}

var config struct {
	Guild map[string]struct {
		Topic []topic `json:"topic"`
	} `json:"guild"`
}

var Commands = make(map[string]*f.Command)

func init() {
	Commands["bpconfig"] = &f.Command{
		Name: "The Bad Post Configuration Tool",
		Help: `bpconfig is the modificaation tool for the bad post timer. It is used to define
topics and triggers (as regular expressions). If you have retained vestages of
your sanity and don't want to learn regex, then TL;DR to match "word" not within 
other words, use //\b(word)\b//.
If you're a regex elder god, it uses https://github.com/google/re2/wiki/Syntax (minus \C)
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
			Topic []topic `json:"topic"`
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
	defer config.Guild[guildID].Topics = localTopics
	defer dat.Save("incident-counter/config.json")
	for li := range localTopics {
		for i := range localTopics[li].Triggers {
			match, err := regexp.MatchString(localTopics[li].Triggers[i], m.Message.Content)
			if err != nil {
				dat.Log.Println(err)
				return
			}
			if match {
				matchedTopic(localTopics[li], s, m.Message)
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
	defer config.Guild[guildID].Topics = localTopics
	defer dat.Save("incident-counter/config.json")

	flgs := flags.Parse(message.Content)
	var topicname string
	for i := range flgs {
		if flgs[i].Name == "--unflagged" {
			topicname = flgs[i].Value
		}
	}
}

const japeFlavourText = []string{
	"You went too far, you fool.",
	"Awe heck, I can't believe you've done this.",
	"OP is a coward.",
	"Think before you speak.",
	"You made a stupid post so now you get no rights.",
}

// If matchedTopic is called, it is already assumed the regex has been met.
// Regex will not be checked.
func matchedTopic(t topic, s *dsg.Session, m *dsg.Message) {
	rand.Seed(time.Now().Unix())
	if time.Since(t.Since) > t.Streak {
		s.ChannelMessageSendEmbed(m.ChannelID, &dsg.MessageEmbed{
			Author:      &dsg.MessageEmbedAuthor{},
			Color:       0x073642,
			Title:       fmt.Sprintf("**Streak Broken. New Record!** %v", japeFlavourText[rand.Intn(len(japeFlavourText))]),
			Description: fmt.Sprintf("%v has broken the %v streak.", m.Author.Mention(), t.Name),
			Thumbnail: &dsg.MessageEmbedThumbnail{
				URL:    fmt.Sprintf("https://cdn.whitmans.io/streakreactions/newstreak.png"),
				Width:  64,
				Height: 64,
			},
			Fields: []*dsg.MessageEmbedField{
				&dsg.MessageEmbedField{
					Name: "Previous Streak",
					Value: fmt.Sprintf("Set by <@%v> at %v, and broken by <@%v> at %v (%v) with the messsage\n```%v```",
						t.StreakUser,
						t.StreakSet.Format("Mon, Jan 2 2006 at 15:04 (MST)"),
						t.BreakUser,
						t.BreakTime.Format("Mon, Jan 2 2006 at 15:04 (MST)"),
						t.BreakTime.Sub(t.StreakSet),
						t.Message),
					Inline: true,
				},
				&dsg.MessageEmbedField{
					Name: "Highest Streak",
					Value: fmt.Sprintf("Set by <@%v> at %v, and broken by %v at %v (%v) with the message\n```%v```",
						t.LastUser,
						t.LastSet.Format("Mon, Jan 2 2006 at 15:04 (MST)"),
						m.Author.Mention(),
						time.Now().Format("Mon, Jan 2 2006 at 15:04 (MST)"),
						time.Since(t.LastSet),
						m.Message.Content),
					Inline: true,
				},
			},
		})
	} else {
		s.ChannelMessageSendEmbed(m.ChannelID, &dsg.MessageEmbed{
			Author:      &dsg.MessageEmbedAuthor{},
			Color:       0x073642,
			Title:       fmt.Sprintf("**Streak Broken.** %v", japeFlavourText[rand.Intn(len(japeFlavourText))]),
			Description: fmt.Sprintf("%v has broken the %v streak.", m.Author.Mention(), t.Name),
			Thumbnail: &dsg.MessageEmbedThumbnail{
				URL:    fmt.Sprintf("https://cdn.whitmans.io/streakreactions/streakbroken.png"),
				Width:  64,
				Height: 64,
			},
			Fields: []*dsg.MessageEmbedField{
				&dsg.MessageEmbedField{
					Name: "Highest Streak",
					Value: fmt.Sprintf("Set by <@%v> at %v, and broken by <@%v> at %v (%v) with the messsage\n```%v```",
						t.StreakUser,
						t.StreakSet.Format("Mon, Jan 2 2006 at 15:04 (MST)"),
						t.BreakUser,
						t.BreakTime.Format("Mon, Jan 2 2006 at 15:04 (MST)"),
						t.BreakTime.Sub(t.StreakSet),
						t.Message),
					Inline: true,
				},
				&dsg.MessageEmbedField{
					Name: "Current Streak",
					Value: fmt.Sprintf("Set by <@%v> at %v, and broken by %v at %v (%v) with the message\n```%v```",
						t.LastUser,
						t.LastSet.Format("Mon, Jan 2 2006 at 15:04 (MST)"),
						m.Author.Mention(),
						time.Now().Format("Mon, Jan 2 2006 at 15:04 (MST)"),
						time.Since(t.LastSet),
						m.Message.Content),
					Inline: true,
				},
			},
		})
	}
}
