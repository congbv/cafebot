// chat holds all stuff related to telegram chat interactions
package chat

import (
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
		func() { err = initIntrHandler(bot, conf.Cafe) },
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

		reqdata := update.CallbackQuery.Data
		intrData, opData, err := splitReqData(reqdata)
		if err != nil {
			log.Errorf("splitting request data %+v: %s", reqdata, err)
			return
		}

		order := s.processOrder(update.CallbackQuery.From, opData)
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

func splitReqData(reqdata string) (string, string, error) {
	splited := strings.Split(reqdata, "?")
	opdata := ""
	if len(splited) > 1 {
		opdata = splited[1]
	}
	return splited[0], opdata, nil
}
