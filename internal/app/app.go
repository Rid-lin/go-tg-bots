package app

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rid-lin/go-tg-bots/for_Vasiliy/internal/config"
	log "github.com/sirupsen/logrus"

	tdlib "github.com/Arman92/go-tdlib"
)

type App struct {
	cfg    *config.Config
	client *tdlib.Client
	Log    *log.Logger
}

func New(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
		// Create new instance of client
		client: tdlib.NewClient(tdlib.Config{
			APIID:               cfg.APIID,
			APIHash:             cfg.APIHash,
			SystemLanguageCode:  "en",
			DeviceModel:         "Server",
			SystemVersion:       "1.0.0",
			ApplicationVersion:  "1.0.0",
			UseMessageDatabase:  true,
			UseFileDatabase:     true,
			UseChatInfoDatabase: true,
			UseTestDataCenter:   false,
			DatabaseDirectory:   "./tdlib-db",
			FileDirectory:       "./tdlib-files",
			IgnoreFileNames:     false,
		}),
		Log: log.New(),
	}

}

func (a *App) Configure() {
	a.configureLogger()
	a.configureClient()
	a.authorizationClient()
}

func (a *app) Start() {

	// https://github.com/KaoriEl/go-tdlib/blob/master/examples/customEvents/getCustomEvents.go
	// Handle Ctrl+C , Gracefully exit and shutdown tdlib
	var ch = make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		a.client.Close()
		a.client.DestroyInstance()
		os.Exit(1)
	}()

	go func() {
		var ChatIDSearch int64
		// Create an filter function which will be used to filter out unwanted tdlib messages
		eventFilter := func(msg *tdlib.TdMessage) bool {
			updateMsg, ok := (*msg).(*tdlib.UpdateNewMessage)
			if !ok {
				return false
			}
			sender, ok := updateMsg.Message.Sender.(*tdlib.MessageSenderUser)
			if ok {
				if sender.GetMessageSenderEnum() == tdlib.MessageSenderUserType {
					a.log.Debugf("UserID:%v,", sender.UserID)
				}
			}
			a.log.Debugf("ChatID:%v\n", updateMsg.Message.ChatID)
			flag := false
			for _, ChatIDs := range a.cfg.ChatsIDSearch {
				if updateMsg.Message.ChatID == ChatIDs {
					ChatIDSearch = ChatIDs
					flag = true
					break
				}
			}
			// ChatIDSearch = a.cfg.ChatIDSearch
			// flag := updateMsg.Message.ChatID == a.cfg.ChatIDSearch
			return flag
			// if updateMsg.Message.Sender.GetMessageSenderEnum() == tdlib.MessageSenderUserType {
			// 	sender := updateMsg.Message.Sender.(*tdlib.MessageSenderUser)
			// 	return sender.UserID == 1055350095
			// }
			// return false
		}

		// Here we can add a receiver to retreive any message type we want
		// We like to get UpdateNewMessage events and with a specific FilterFunc
		receiver := a.client.AddEventReceiver(&tdlib.UpdateNewMessage{}, eventFilter, 5)
		for newMsg := range receiver.Chan {
			// fmt.Println(newMsg)
			updateMsg, ok := (newMsg).(*tdlib.UpdateNewMessage)
			if !ok {
				continue
			}
			// We assume the message content is simple text: (should be more sophisticated for general use)
			msgText := updateMsg.Message.Content.(*tdlib.MessageText)
			if !ok {
				continue
			}
			flag := false
			for _, word := range a.cfg.Words {
				if strings.Contains(fmt.Sprint(msgText.Text), word) {
					flag = true
					break
				}
			}
			// https://github.com/KaoriEl/go-tdlib/blob/master/examples/sendText/sendText.go

			// Should get chatID somehow, check out "getChats" example
			if flag {
				fmt.Println("Search word in MsgText:  ", msgText.Text)
				option := tdlib.MessageSendOptions{
					DisableNotification: false, // Pass true to disable notification for the message
					FromBackground:      false, // Pass true if the message is sent from the background
				}
				_, err := a.client.ForwardMessages(a.cfg.ChatID, ChatIDSearch, []int64{updateMsg.Message.ID}, &option, false, false)
				if err != nil {
					a.log.Error(err)
				}
			}
			// fmt.Println("MsgText:  ", msgText.Text)
			// fmt.Print("\n")

		}

	}()

	// rawUpdates gets all updates comming from tdlib
	rawUpdates := a.client.GetRawUpdatesChannel(100)
	for update := range rawUpdates {
		// Show all updates
		log.Trace(update.Data)
		// fmt.Println(update.Data)
		// fmt.Print("\n\n")
	}
}

func (a *app) configureLogger() {
	lvl, err := log.ParseLevel(a.cfg.LogLevel)
	if err != nil {
		log.Errorf("Error parse the level of logs (%v). Installed by default = Info", a.cfg.LogLevel)
		lvl, _ = log.ParseLevel("info")
	}
	a.Log.SetLevel(lvl)
}

func (a *app) configureClient() {
	_, _ = a.client.SetLogVerbosityLevel(1)
	// a.client.SetFilePath("./errors.txt")
}

func (a *App) authorizationClient() {
	// https://github.com/KaoriEl/go-tdlib/blob/master/examples/authorization/basicAuthorization.go
	for {
		currentState, _ := a.client.Authorize()
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			// fmt.Print("Enter phone: ")
			// var number string
			// fmt.Scanln(&number)
			number := a.cfg.PhoneNumber
			_, err := a.client.SendPhoneNumber(number)
			if err != nil {
				fmt.Printf("Error sending phone number: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitCodeType {
			fmt.Print("Enter code: ")
			var code string
			fmt.Scanln(&code)
			_, err := a.client.SendAuthCode(code)
			if err != nil {
				fmt.Printf("Error sending auth code : %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPasswordType {
			// fmt.Print("Enter Password: ")
			// var password string
			// fmt.Scanln(&password)
			password := a.cfg.Password
			_, err := a.client.SendAuthPassword(password)
			if err != nil {
				fmt.Printf("Error sending auth password: %v", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType {
			fmt.Println("Authorization Ready! Let's rock")
			break
		}
	}
}

func (a *App) Start() {

	// https://github.com/KaoriEl/go-tdlib/blob/master/examples/customEvents/getCustomEvents.go
	// Handle Ctrl+C , Gracefully exit and shutdown tdlib
	var ch = make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		a.client.DestroyInstance()
		os.Exit(1)
	}()

	go func() {
		var ChatIDSearchTrue int64
		// Create an filter function which will be used to filter out unwanted tdlib messages
		eventFilter := func(msg *tdlib.TdMessage) bool {
			updateMsg := (*msg).(*tdlib.UpdateNewMessage)
			// if updateMsg.Message.Sender.GetMessageSenderEnum() == tdlib.MessageSenderUserType {
			// 	sender := updateMsg.Message.Sender.(*tdlib.MessageSenderUser)
			// 	a.log.Debugf("UserID:%v,", sender.UserID)
			// }
			a.Log.Debugf("ChatID:%v\n", updateMsg.Message.ChatID)
			flag := false
			for _, ChatIDSearch := range a.cfg.ChatsIDSearch {
				if updateMsg.Message.ChatID == ChatIDSearch {
					ChatIDSearchTrue = updateMsg.Message.ChatID
					flag = true
					break
				}
			}
			return flag
			// if updateMsg.Message.Sender.GetMessageSenderEnum() == tdlib.MessageSenderUserType {
			// 	sender := updateMsg.Message.Sender.(*tdlib.MessageSenderUser)
			// 	return sender.UserID == 1055350095
			// }
			// return false
		}

		// Here we can add a receiver to retreive any message type we want
		// We like to get UpdateNewMessage events and with a specific FilterFunc
		receiver := a.client.AddEventReceiver(&tdlib.UpdateNewMessage{}, eventFilter, 5)
		for newMsg := range receiver.Chan {
			// fmt.Println(newMsg)

			updateMsg, ok := (newMsg).(*tdlib.UpdateNewMessage)
			if !ok {
				continue
			}
			// We assume the message content is simple text: (should be more sophisticated for general use)
			msgText, ok := updateMsg.Message.Content.(*tdlib.MessageText)
			if !ok {
				continue
			}
			flag := false
			for _, word := range a.cfg.Words {
				word = strings.ToLower(word)
				text := strings.ToLower(fmt.Sprint(msgText.Text))
				if strings.Contains(text, word) {
					flag = true
					break
				}
			}
			if flag {
				fmt.Println("Search word in MsgText:  ", msgText.Text)
				fmt.Print("\n")

				// https://github.com/KaoriEl/go-tdlib/blob/master/examples/sendText/sendText.go

				// Should get chatID somehow, check out "getChats" example
				option := tdlib.MessageSendOptions{
					DisableNotification: false, // Pass true to disable notification for the message
					FromBackground:      false, // Pass true if the message is sent from the background
				}
				_, err := a.client.ForwardMessages(a.cfg.ChatID, ChatIDSearchTrue, []int64{updateMsg.Message.ID}, &option, false, false)
				if err != nil {
					a.Log.Error(err)
				}
			}
		}

	}()

	// rawUpdates gets all updates comming from tdlib
	rawUpdates := a.client.GetRawUpdatesChannel(100)
	for update := range rawUpdates {
		// Show all updates
		log.Trace(update.Data)
		// fmt.Println(update.Data)
		// fmt.Print("\n\n")
	}
}
