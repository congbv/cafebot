// chat holds all stuff related to telegram chat interactions
package chat

import (
	"errors"
	"net/url"
	"strings"
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
		func() { err = initIntrHandler(conf.Cafe, bot) },
	} {
		if err != nil {
			continue
		}
		f()
	}
	if err != nil {
		return nil, err
	}

	s := &service{
		conf: conf,
		bot:  bot,

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
		log.Debugf("received callback: %#v", update.CallbackQuery)

		reqdata, params, err := splitDataParams(update.CallbackQuery.Data)
		if err != nil {
			log.Errorf("cutOffParams: %s", err)
			return
		}

		user := update.CallbackQuery.From
		s.processOrder(user, params)
		order := *s.order.Get(user)

		intr.handle(reqdata, update, order)
		return
	}

	if update.Message != nil {
		log.Debugf("received message: %#v", update.Message)

		cmd.handle(update.Message.Command(), update)
	}
}

func splitDataParams(reqdata string) (string, url.Values, error) {
	splited := strings.Split(reqdata, "?")
	data := splited[0]
	if data == "" {
		return "", nil, errors.New("empty data value")
	}
	if len(splited) < 2 {
		return data, nil, nil
	}
	params, err := url.ParseQuery(splited[1])
	if err != nil {
		return "", nil, err
	}
	return data, params, nil
}
