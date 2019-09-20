// chat holds all stuff related to telegram chat interactions
package chat

import (
	"fmt"
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yarikbratashchuk/cafebot/config"
	"github.com/yarikbratashchuk/cafebot/order"
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
		if finished {
			err = s.sendOrderToChannel(s.conf.Cafe.Chan, order)
			if err != nil {
				log.Errorf(
					"sending order to cafe: %s, order: %+v",
					err,
					order,
				)
				s.sendError(q.Message.Chat.ID)
				return
			}
		}
		intr.handle(intrData, update, order)

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

func (s *service) sendOrderToChannel(channel string, o order.Order) error {
	log.Debugf("sending order to cafe: u: %s, o: %+v", o.User.UserName, o)

	text := fmt.Sprintf(
		"<b>%s %s\n(@%s)</b>\n\n%s",
		o.User.FirstName,
		o.User.LastName,
		o.User.UserName,
		generatePreviewText(o),
	)

	msg := api.NewMessageToChannel(channel, text)
	msg.ParseMode = api.ModeHTML
	_, err := s.bot.Send(msg)

	return err
}

func (s *service) sendError(chatID int64) {
	s.bot.Send(api.NewMessage(chatID, text["error"]))
}
