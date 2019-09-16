package chat

import (
	"net/url"
	"testing"
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yarikbratashchuk/cafebot/config"
	"github.com/yarikbratashchuk/cafebot/order"
)

func TestServiceProcessOrder(t *testing.T) {
	t.Parallel()

	o := order.NewInMemoryService(config.Config{})
	s := service{order: o}

	at, _ := time.Parse("15:04", "15:00")

	user1 := &api.User{
		ID:       1,
		UserName: "@testuser1",
	}
	user2 := &api.User{
		ID:       2,
		UserName: "@testuser1",
	}

	cases := []struct {
		test string

		user          *api.User
		params        url.Values
		expectedOrder order.Order
		HasErr        bool
	}{{
		test: "valid case",

		user: user1,
		params: url.Values{
			"add_meal":     []string{"testMeal1"},
			"set_time":     []string{"15:00"},
			"set_takeaway": []string{"1"},
		},
		expectedOrder: order.Order{
			User: user1,

			Meal:     []string{"testMeal1"},
			Time:     &at,
			Takeaway: true,
		},
		HasErr: false,
	}, {
		test: "remove meal case",

		user: user2,
		params: url.Values{
			"remove_meal":  []string{"testMeal1"},
			"set_time":     []string{"15:00"},
			"set_takeaway": []string{"0"},
		},
		expectedOrder: order.Order{
			User: user2,

			Meal:     []string{},
			Time:     &at,
			Takeaway: false,
		},
		HasErr: true,
	}}
	for _, c := range cases {
		c := c
		t.Run(c.test, func(t *testing.T) {
			s.processOrder(c.user, c.params)

			actualOrder, err := s.order.FinishOrder(c.user)
			if err != nil && !c.HasErr {
				t.Fatalf("finish order: must have no err %s", err)
			} else if err == nil && c.HasErr {
				t.Fatal("finish order: err must be not nil")
			} else if err != nil && c.HasErr {
				return
			}

			if actualOrder.User != c.user {
				t.Error("user is wrong")
			}
			if len(actualOrder.Meal) != len(c.expectedOrder.Meal) {
				t.Fatal("meal is wrong")
			}
			if actualOrder.Meal[0] != c.expectedOrder.Meal[0] {
				t.Error("meal is wrong")
			}
			if *actualOrder.Time != *c.expectedOrder.Time {
				t.Error("time is wrong")
			}
			if actualOrder.Takeaway != c.expectedOrder.Takeaway {
				t.Error("takeaway is wrong")
			}
		})
	}
}
