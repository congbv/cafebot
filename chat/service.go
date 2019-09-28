// chat holds all stuff related to telegram chat interactions
package chat

import (
	"fmt"
	"time"

	"cafebot/config"
	"cafebot/order"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
)

// service is responsible for:
// - receiving commands from telegram chat
// - making calls to the order service
// - calling command and interaction handlers
type service struct {
	conf  config.Config
	bot   *api.BotAPI
	order order.Service

	done chan struct{}
}

func NewService(
	conf config.Config,
	orderService order.Service,
) (*service, error) {
	bot, err := initBotAPI(conf.TgAPIToken)
	if err != nil {
		return nil, err
	}

	for _, f := range []func(){
		func() { err = initCmdHandler(bot) },
		func() { err = initIntrHandler(bot, conf.Cafe) },
	} {
		f()
		if err != nil {
			return nil, err
		}
	}

	s := &service{
		conf:  conf,
		bot:   bot,
		order: orderService,

		done: make(chan struct{}),
	}

	return s, nil
}

func (s *service) Run() error {
	log.Info("starting chat interface")

	u := api.NewUpdate(0)
	u.Timeout = 60

	updates, err := s.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	updates.Clear()

	go func() {
		for {
			select {
			case update := <-updates:
				s.handleUpdate(update)
			case <-s.done:
				return
			case <-time.After(1 * time.Second):
			}
		}
	}()

	return nil
}

func (s *service) Stop() { close(s.done) }

func (s *service) handleUpdate(update api.Update) {
	if update.CallbackQuery != nil {
		log.Debugf("handling callback query: %+v", update.CallbackQuery)

		q := update.CallbackQuery
		reqdata := q.Data
		intrData, opData, err := splitIntrOpData(reqdata)
		if err != nil {
			log.Errorf("splitting request data %+v: %s", reqdata, err)
			return
		}

		order, finished := s.processOrder(q.From, opData)
		if order == nil {
			return
		}
		if finished {
			errText := s.sendOrderToChannel(
				s.conf.Cafe.OrderChan,
				*order,
			)
			if errText != "" {
				log.Errorf(
					"sending order to cafe: %s, order: %s",
					errText,
					order,
				)
				if errText != text["err_no_username"] {
					errText = text["err_internal"]
				}
				s.sendError(q.Message.Chat.ID, errText)
				return
			}
		}
		intr.handle(intrData, update, *order)

		return
	}

	if update.Message != nil {
		log.Debugf("handling command: %+v", update.Message)

		command := update.Message.Command()
		err := cmd.handle(command, update)
		if err != nil {
			log.Errorf("running command %+v: %s", command, err)
			return
		}
	}
}

func (s *service) sendOrderToChannel(channel string, o order.Order) string {
	log.Debugf("sending order to cafe: u: %s, o: %+v", o.User.UserName, o)

	if o.User.UserName == "" {
		return text["err_no_username"]
	}

	userNameText := generateUserNameText(o)
	previewText := generatePreviewText(o)

	msg := api.NewMessageToChannel(
		channel,
		fmt.Sprintf("%s\n\n%s", userNameText, previewText),
	)
	msg.ParseMode = api.ModeHTML
	_, err := s.bot.Send(msg)
	if err != nil {
		return err.Error()
	}

	return ""
}

func (s *service) sendError(chatID int64, errText string) {
	s.bot.Send(api.NewMessage(chatID, errText))
}
