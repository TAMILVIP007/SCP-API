package src

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

var bot *gotgbot.Bot

func init() {
	var err error
	bot, err = gotgbot.NewBot(Envars.Token, &gotgbot.BotOpts{
		Client: http.Client{},
		DefaultRequestOpts: &gotgbot.RequestOpts{
			Timeout: gotgbot.DefaultTimeout,
			APIURL:  gotgbot.DefaultAPIURL,
		},
	})
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}
	updater := ext.NewUpdater(&ext.UpdaterOpts{
		Dispatcher: ext.NewDispatcher(&ext.DispatcherOpts{
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				log.Println("an error occurred while handling update:", err.Error())
				return ext.DispatcherActionNoop
			},
			MaxRoutines: ext.DefaultMaxRoutines,
		}),
	})
	dispatcher := updater.Dispatcher
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	err = updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	log.Printf("%s has been started...\n", bot.User.Username)
}

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, "hi", &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}

func SendRequest(request interface{}) bool {
	var msgtosend string
	switch req := request.(type) {
	case *BanRequest:
		msgtosend += "<b>New Ban Request</b>\n"
		msgtosend += fmt.Sprintf("<b>User ID:</b> %s\n", req.UserId)
		msgtosend += fmt.Sprintf("<b>Reason:</b> %s\n", req.Reason)
		msgtosend += fmt.Sprintf("<b>Requested By:</b> %s\n", req.From)
		msgtosend += fmt.Sprintf("<b>Ban Class:</b> %s\n", req.BanClass)
		msgtosend += fmt.Sprintf("<b>Note:</b> %s\n", req.Notes)
		msgtosend += fmt.Sprintf("<b>Proof:</b> %s\n", req.EvidenceLink)
	case *UnbanRequest:
		msgtosend += "<b>New Unban Request</b>\n"
		msgtosend += fmt.Sprintf("<b>User ID:</b> %s\n", req.UserId)
		msgtosend += fmt.Sprintf("<b>Reason:</b> %s\n", req.Reason)
		msgtosend += fmt.Sprintf("<b>Requested By:</b> %s\n", req.From)
	default:
		return false
	}

	_, err := bot.SendMessage(Envars.LogChat, msgtosend, &gotgbot.SendMessageOpts{
		ParseMode: "html",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "Approve",
						CallbackData: "approve",
					},
					{
						Text:         "Deny",
						CallbackData: "deny",
					},
				},
			},
		},
	})
	return err == nil
}
