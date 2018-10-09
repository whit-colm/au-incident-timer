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

/*var configDefaults = topic{
	Name: "Default Topic",
	Triggers: []string{},
	Sine: time.Now(),
	LastUser: "238741157960482816",
	LastSet: time.Now(),
	BreakUser: "238741157960482816",
	BreakTime: time.Now(),
	StreakUser: "238741157960482816",
	StreakSet: time.Now(),
	Message: "Default Message",
}*/

var config struct {
	Guild map[string]struct {
		Topics map[string]topic `json:"topic"`
	} `json:"guild"`
}

var Commands = make(map[string]*f.Command)

func init() {
	Commands["bpconfig"] = &f.Command{
		Name: "The Bad Post Configuration Tool",
		Help: `bpconfig is the modificaation tool for the bad post timer. It is used to define
topics and triggers (as regular expressions). If you have retained vestages of
your sanity and don't want to learn regex, then TL;DR to match "word" not within 
other words, use \b(word)\b
If you're a regex elder god, it uses https://github.com/google/re2/wiki/Syntax (minus \C)
# Arguments
**--new <topic name>** : Create a new topic with the topic name <topic name>
**--del <topic name>** : Delete <topic name>. Dangerous!
**-ls** : List all topics
**<topic name> <--addtrigger|-a> <regex>** : Add a trigger to <topic name> fulfilling <regex> (one regex per --add is supported)
**<topic name> <--rmtrigger|-rm> <item number>** : Remove the <item number> trigger from <topic name>
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
			Topics map[string]topic `json:"topic"`
		})
	}

	f.Session.AddHandler(incidentHandler)
}

func incidentHandler(s *dsg.Session, m *dsg.MessageCreate) {
	if m.Message.Author.Bot {
		return
	}
	var guildID string
	if channel, err := session.Channel(m.Message.ChannelID); err != nil {
		dat.Log.Println(err)
		return
	} else {
		guildID = channel.GuildID
	}
	if config.Guild[guildID].Topics == nil {
		config.Guild[guildID].Topics = make(map[string]topic)
	}
	// Local topics reduces verbosity by grabbing the config.Guild[guildID] value.
	localTopics := config.Guild[guildID].Topics
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
	config.Guild[guildID].Topics = localTopics
}

func bpconfig(session *dsg.Session, message *dsg.Message) {
	var guildID string
	if channel, err := session.Channel(message.ChannelID); err != nil {
		dat.Log.Println(err)
		return
	} else {
		guildID = channel.GuildID
	}
	if config.Guild[guildID].Topics == nil {
		config.Guild[guildID].Topics = make(map[string]topic)
	}
	// Local topics reduces verbosity by grabbing the config.Guild[guildID] value.
	localTopics := config.Guild[guildID].Topics
	defer dat.Save("incident-counter/config.json")

	flgs := flags.Parse(message.Content)
	var topicname string
	for i := range flgs {
		if flgs[i].Name == "--unflagged" {
			topicname = flgs[i].Value
		} else if flgs[i].Name == "--new" {
			localTopics[flgs[i].Value] = topic{}
			localTopics[flgs[i].Value].Name = flgs[i].Value
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Created topic %v.", flgs[i].Value))
		} else if flgs[i].Name == "--del" {
			delete(localTopics, flgs[i].Value)
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Deleted topic %v.", flgs[i].Value))
		}
		if topicname != "" {
			switch flgs[i].Name {
			case "--add", "-a":
				_, err := regexp.MatchString(flgs[i].Value, "Some String")
				if err != nil {
					session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Hey nerd, **%v** doesn't make any sense. Try again idiot. `(error: %v)`", flgs[i].Value, err))
					continue
				}
				localTopics[topicname].Triggers = append(localTopics[topicname].Triggers, flgs[i].Value)
				session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Created uh..... trigger %v.", flgs[i].Value))
			case "--remove", "-rm":
				index, err := strconv.Atoi(flgs[i].Value)
				if err != nil {
					dat.Log.Println(err)
					dat.AlertDiscord(session, message, err)
				}
				go deleteFromSlice(&topicname.Triggers, index, false)
				session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Trigger #%v committed not in memory.", flgs[i].Value))
			case "--list", "-ls", "-l":
				str := fmt.Sprintf("List of triggers for %v:", topicname)
				for i := range localTopics[topicname].Triggers {
					str += fmt.Sprintf("\n`%v`**:**`%v`**.**", i, localTopics[topicname].Triggers)
				}
				session.ChannelMessageSend(message.ChannelID, str)
			}
		}
	}
	config.Guild[guildID].Topics = localTopics
}

var japeFlavourText = []string{
	"You went too far, you fool.",
	"Awe heck, I can't believe you've done this.",
	"OP is a coward.",
	"Think before you speak.",
	"You made a stupid post so now you get no rights.",
	"Do you value your toes?",
}

func deleteFromSlice(a *[]interface{}, i int, preserveOrder bool) {
	if preserveOrder {
		a = append(a[:i], a[i+1:]...)
	} else {
		a[i] = a[len(a)-1]
		a = a[:len(a)-1]
	}
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
