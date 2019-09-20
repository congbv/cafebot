package chat

import (
	"strings"
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yarikbratashchuk/cafebot/order"
)

type orderOp string

const (
	opInit       orderOp = "init"
	opWhere      orderOp = "where"
	opWhen       orderOp = "when"
	opAddMeal    orderOp = "add_meal"
	opRemoveMeal orderOp = "remove_meal"
	opFinish     orderOp = "finish_order"

	opWhereHere     = "here"
	opWhereTakeaway = "takeaway"
)

// processOrder handles user order updates
func (s *service) processOrder(u *api.User, opData string) (*order.Order, bool) {
	var finished bool
	o := s.order.Get(u)

	splited := strings.Split(opData, "=")
	if len(splited) < 2 {
		return o, finished
	}

	op := orderOp(splited[0])
	val := splited[1]

	switch op {
	case opInit:
		o = s.order.InitOrder(u)

	case opWhere:
		o = s.order.SetTakeaway(u, val == opWhereTakeaway)

	case opWhen:
		t, err := time.Parse("15:04", val)
		if err != nil {
			log.Errorf("parsing time %+v: %s", val, err)
			return o, finished
		}
		o = s.order.SetTime(u, t)

	case opAddMeal, opRemoveMeal:
		meal, ok := s.conf.Cafe.Menu.MealByHash(val)
		if !ok {
			log.Errorf("getting meal by hash: %+v: %s", val)
			return o, finished
		}
		switch op {
		case opAddMeal:
			o = s.order.AddMeal(u, meal)
		case opRemoveMeal:
			o = s.order.RemoveMeal(u, meal)
		}

	case opFinish:
		o, finished = s.order.FinishOrder(u)

	default:
	}

	return o, finished
}
