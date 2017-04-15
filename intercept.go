package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
	"github.com/legolord208/stdutil"
)

func messageCreate(session *discordgo.Session, e *discordgo.MessageCreate) {
	channel, err := session.Channel(e.ChannelID)
	if err != nil {
		stdutil.PrintErr(tl("failed.channel"), err)
		return
	}

	var guild *discordgo.Guild
	if !channel.IsPrivate {
		guild, err = session.Guild(channel.GuildID)
		if err != nil {
			stdutil.PrintErr(tl("failed.guild"), err)
			return
		}
	}

	if messageCommand(session, e.Message, guild, channel) {
		return
	}

	lastMsg = location{
		guild:   guild,
		channel: channel,
	}

	hasOutput := false

	print := false
Switch:
	switch messages {
	case MessagesAll:
		print = true
	case MessagesPrivate:
		if channel.IsPrivate {
			print = true
		}
	case MessagesMentions:
		if channel.IsPrivate || e.MentionEveryone {
			print = true
			break
		}

		for _, u := range e.Mentions {
			if u.ID == UserId {
				print = true
				break Switch
			}
		}

		user, err := session.GuildMember(guild.ID, UserId)
		if err != nil {
			stdutil.PrintErr(tl("failed.user"), err)
			break
		}

		for _, role := range user.Roles {
			for _, role2 := range e.MentionRoles {
				if role == role2 {
					print = true
					break Switch
				}
			}
		}
	case MessagesCurrent:
		if (guild == nil || loc.guild == nil) && loc.channel != nil && channel.ID != loc.channel.ID {
			break
		}
		if guild != nil && loc.guild != nil && guild.ID != loc.guild.ID {
			break
		}

		print = true
	}
	if print {
		printMessage(session, e.Message, true, guild, channel, color.Output)
		hasOutput = true
	}

	if len(luaMessageEvents) > 0 {
		hasOutput = true

		color.Unset()
		ColorAutomated.Set()

		fmt.Print("\r" + strings.Repeat(" ", 20) + "\r")
		luaMessageEvent(session, e.Message)

		color.Unset()
		ColorDefault.Set()
	}
	if hasOutput {
		printPointer(session)
	}
}

func messageCommand(session *discordgo.Session, e *discordgo.Message, guild *discordgo.Guild, channel *discordgo.Channel) (isCmd bool) {
	if e.Author.ID != UserId {
		return
	} else if !intercept {
		return
	}

	contents := strings.TrimSpace(e.Content)
	if !strings.HasPrefix(contents, "console.") {
		return
	}
	cmd := contents[len("console."):]

	isCmd = true

	if strings.EqualFold(cmd, "ping") {
		first := time.Now().UTC()

		timestamp, err := e.Timestamp.Parse()
		if err != nil {
			stdutil.PrintErr(tl("failed.timestamp"), err)
			return
		}

		in := first.Sub(timestamp)

		// Discord 'bug' makes us receive the message before the timestamp, sometimes.
		text := ""
		if in.Nanoseconds() >= 0 {
			text += "Incoming: `" + in.String() + "`"
		} else {
			text += "Message was received earlier than timestamp. Discord bug."
		}

		middle := time.Now().UTC()

		msg := &discordgo.MessageEdit{}
		msg.SetContent("Pong! 1/2")
		msg.Embed = &discordgo.MessageEmbed{
			Description: text + "\nCalculating outgoing..",
		}
		_, err = session.ChannelMessageEditComplex(e.ChannelID, e.ID, msg)

		last := time.Now().UTC()

		text += "\nOutgoing: `" + last.Sub(middle).String() + "`"
		text += "\n\n\nIncoming is the time it takes for the message to reach DiscordConsole."
		text += "\nOutgoing is the time it takes for DiscordConsole to reach discord."

		msg = &discordgo.MessageEdit{}
		msg.SetContent("Pong! 2/2")
		msg.Embed = &discordgo.MessageEmbed{
			Description: text,
		}
		_, err = session.ChannelMessageEditComplex(e.ChannelID, e.ID, msg)
		if err != nil {
			stdutil.PrintErr(tl("failed.msg.edit"), err)
		}
		return
	}

	loc.push(guild, channel)

	var w io.Writer
	var str *bytes.Buffer
	if output {
		str = bytes.NewBuffer(nil)
		w = str
	} else {
		go func() {
			err := session.ChannelMessageDelete(e.ChannelID, e.ID)
			if err != nil {
				stdutil.PrintErr(tl("failed.msg.delete"), err)
			}
		}()
		color.Unset()
		ColorAutomated.Set()

		fmt.Println(cmd)
		w = color.Output
	}
	command(session, cmd, w)

	if !output {
		color.Unset()
		ColorDefault.Set()
		printPointer(session)
	} else {
		first := true
		send := func(buf string) {
			if buf == "" {
				return
			}

			buf = "```\n" + buf + "\n```"
			if first {
				first = false
				_, err := session.ChannelMessageEdit(e.ChannelID, e.ID, buf)
				if err != nil {
					stdutil.PrintErr(tl("failed.msg.edit"), err)
					return
				}
			} else {
				_, err := session.ChannelMessageSend(e.ChannelID, buf)
				if err != nil {
					stdutil.PrintErr(tl("failed.msg.send"), err)
					return
				}
			}
		}

		buf := ""
		for {
			line, err := str.ReadString('\n')
			if err != nil {
				break
			}

			if len(line)+len(buf)+8 < MsgLimit {
				buf += line
			} else {
				send(buf)
				buf = ""
			}
		}
		send(buf)
	}

	color.Unset()
	ColorDefault.Set()
	return
}