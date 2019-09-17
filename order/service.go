package order

import (
	"errors"
	"sync"
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yarikbratashchuk/cafebot/config"
)

type Order struct {
	User     *api.User
	Meal     []string
	Takeaway *bool
	Time     *time.Time
}

// inMemService is the most naive in memory implementation of the Service
type inMemService struct {
	conf   config.Config
	mu     *sync.RWMutex
	orders map[int]*Order
}

func NewInMemoryService(conf config.Config) Service {
	return &inMemService{
		conf:   conf,
		mu:     &sync.RWMutex{},
		orders: make(map[int]*Order, 10),
	}
}

func (s *inMemService) Get(u *api.User) *Order {
	return s.get(u)
}

func (s *inMemService) InitOrder(u *api.User) *Order {
	order, ok := s.orders[u.ID]
	if ok {
		return order
	}
	order = &Order{
		User: u,
		Meal: make([]string, 0, 1),
	}
	log.Debugf("InitOrder: %#v", order)
	s.mu.Lock()
	s.orders[u.ID] = order
	s.mu.Unlock()
	return order
}

func (s *inMemService) AddMeal(u *api.User, meal string) *Order {
	order := s.get(u)
	order.Meal = append(order.Meal, meal)
	log.Debugf("AddMeal: %#v", order)
	return order
}

func (s *inMemService) RemoveMeal(u *api.User, meal string) *Order {
	order := s.get(u)
	for i, m := range order.Meal {
		if m == meal {
			order.Meal = append(order.Meal[:i], order.Meal[i+1:]...)
			break
		}
	}
	log.Debugf("RemoveMeal: %#v", order)
	return order
}

func (s *inMemService) SetTime(u *api.User, t time.Time) *Order {
	order := s.get(u)
	order.Time = &t
	log.Debugf("SetTime: %#v", order)
	return order
}

func (s *inMemService) SetTakeaway(u *api.User, takeaway bool) *Order {
	order := s.get(u)
	order.Takeaway = &takeaway
	log.Debugf("Takeaway: %#v", order)
	return order
}

func (s *inMemService) FinishOrder(u *api.User) (*Order, error) {
	if u == nil {
		return nil, errors.New("nil user provided")
	}
	s.mu.RLock()
	order, ok := s.orders[u.ID]
	s.mu.RUnlock()
	if !ok {
		return nil, errors.New("no order for such user")
	}
	if len(order.Meal) == 0 {
		return nil, errors.New("order has no meal selected")
	}
	if order.Time == nil {
		return nil, errors.New("order time is not selected")
	}
	if order.Takeaway == nil {
		return nil, errors.New("order place is not selected")
	}
	s.mu.Lock()
	delete(s.orders, u.ID)
	s.mu.Unlock()
	log.Debugf("FinishOrder: %#v", order)
	return order, nil
}

func (s *inMemService) get(u *api.User) *Order {
	if u == nil {
		return nil
	}
	s.mu.RLock()
	order, ok := s.orders[u.ID]
	s.mu.RUnlock()
	if !ok {
		return s.InitOrder(u)
	}
	return order
}
