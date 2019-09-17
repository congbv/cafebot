package order_test

import (
	"testing"
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yarikbratashchuk/cafebot/config"
	"github.com/yarikbratashchuk/cafebot/order"
)

func TestInMemoryService(t *testing.T) {
	t.Parallel()

	s := order.NewInMemoryService(config.Config{})
	testID := 123
	u := &api.User{ID: testID}

	testInitOrder(t, s, u)
	testAddMeal(t, s, u)
	testRemoveMeal(t, s, u)
	testSetTime(t, s, u)
	testSetTakeaway(t, s, u)
	testFinishOrder(t, s, u)
}

func testInitOrder(t *testing.T, s order.Service, u *api.User) {
	order := s.InitOrder(u)
	if order.User.ID != u.ID {
		t.Error("userID is not set")
	}
	if cap(order.Meal) != 1 {
		t.Error("meal slice should have cap=1")
	}
}

func testAddMeal(t *testing.T, s order.Service, u *api.User) {
	meal := "salad ceasar"
	order := s.AddMeal(u, meal)
	if len(order.Meal) != 1 {
		t.Error("meal was not added")
	}
	if order.User.ID != u.ID {
		t.Error("order has no user")
	}
}

func testRemoveMeal(t *testing.T, s order.Service, u *api.User) {
	meal := "salad ceasar"
	order := s.AddMeal(u, meal)
	prevLen := len(order.Meal)
	order = s.RemoveMeal(u, meal)
	if len(order.Meal) != prevLen-1 {
		t.Error("meal was not removed")
	}
}

func testSetTime(t *testing.T, s order.Service, u *api.User) {
	tt := time.Now().Add(1 * time.Hour)
	order := s.SetTime(u, tt)
	if *order.Time != tt {
		t.Error("time was not set")
	}
	if order.User.ID != u.ID {
		t.Error("order has no user")
	}
}

func testSetTakeaway(t *testing.T, s order.Service, u *api.User) {
	order := s.SetTakeaway(u, true)
	if order.Takeaway == nil || !*order.Takeaway {
		t.Error("takeway value is wrong")
	}
	if order.User.ID != u.ID {
		t.Error("order has no user")
	}
}

func testFinishOrder(t *testing.T, s order.Service, u *api.User) {
	_, err := s.FinishOrder(u)
	if err != nil {
		t.Error("order should be ready")
	}
}
