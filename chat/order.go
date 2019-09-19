package chat

import (
	"strings"
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yarikbratashchuk/cafebot/order"
)

type orderOp string

const (
	opAddMeal    orderOp = "add_meal"
	opRemoveMeal orderOp = "remove_meal"
	opWhen       orderOp = "when"
	opWhere      orderOp = "where"

	opWhereHere     = "here"
	opWhereTakeaway = "takeaway"
)

// processOrder handles user order updates
func (s *service) processOrder(u *api.User, opData string) order.Order {
	o := s.order.Get(u)

	splited := strings.Split(opData, "=")
	if len(splited) < 2 {
		return *o
	}

	op := orderOp(splited[0])
	val := splited[1]

	switch op {
	case opAddMeal, opRemoveMeal:
		meal, ok := s.conf.Cafe.Menu.MealByHash(val)
		if !ok {
			log.Errorf("getting meal by hash: %+v: %s", val)
			return *o
		}
		switch op {
		case opAddMeal:
			o = s.order.AddMeal(u, meal)
		case opRemoveMeal:
			o = s.order.RemoveMeal(u, meal)
		}

	case opWhen:
		t, err := time.Parse("15:04", val)
		if err != nil {
			log.Errorf("parsing time %+v: %s", val, err)
			return *o
		}
		o = s.order.SetTime(u, t)

	case opWhere:
		o = s.order.SetTakeaway(u, val == opWhereTakeaway)

	default:
	}

	// no need to pass pointer further
	return *o
}
