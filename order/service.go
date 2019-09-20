package order

import (
	"bytes"
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

func (o Order) IsReady() bool {
	return len(o.Meal) > 0 && o.Time != nil && o.Takeaway != nil
}

// inMemService is the most naive in-memory
// implementation of the Service
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

func (s *inMemService) InitOrder(u *api.User) *Order {
	order, ok := s.orders[u.ID]
	if ok {
		return order
	}
	order = &Order{
		User: u,
		Meal: make([]string, 0, 1),
	}

	log.Debugf("initializing order: uid: %v", order.User.ID)

	s.mu.Lock()
	s.orders[u.ID] = order
	s.mu.Unlock()

	return order
}

func (s *inMemService) Get(u *api.User) *Order {
	return s.get(u)
}

func (s *inMemService) AddMeal(u *api.User, meal string) *Order {
	order := s.get(u)
	order.Meal = append(order.Meal, meal)

	log.Debugf("adding meal: %s, order: %+v", meal, order)

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

	log.Debugf("removing meal: %s, order: %+v", order)

	return order
}

func (s *inMemService) SetTime(u *api.User, t time.Time) *Order {
	order := s.get(u)
	order.Time = &t

	log.Debugf(
		"setting order time: %v, order: %+v",
		t.Format("15:04"),
		order,
	)

	return order
}

func (s *inMemService) SetTakeaway(u *api.User, takeaway bool) *Order {
	order := s.get(u)
	order.Takeaway = &takeaway

	log.Debugf("setting takeaway: %v, order: %+v", takeaway, order)

	return order
}

type ErrNotComplete struct {
	noMeal, noTime, noTakeaway bool
}

func (e ErrNotComplete) Error() string {
	if !e.noMeal && !e.noTime && !e.noTakeaway {
		return ""
	}
	buf := new(bytes.Buffer)
	buf.WriteString("order has no ")
	if e.noMeal {
		buf.WriteString("meal ")
	}
	if e.noTime {
		buf.WriteString("time ")
	}
	if e.noTakeaway {
		buf.WriteString("takeaway")
	}
	return buf.String()
}

func (e ErrNotComplete) OrderNotComplete() bool {
	return e.noMeal || e.noTime || e.noTakeaway
}

func (s *inMemService) FinishOrder(u *api.User) (*Order, bool) {
	if u == nil {
		return nil, false
	}

	s.mu.RLock()
	order, ok := s.orders[u.ID]
	s.mu.RUnlock()
	if !ok {
		return nil, false
	}

	var err ErrNotComplete
	if len(order.Meal) == 0 {
		err.noMeal = true
	}
	if order.Time == nil {
		err.noTime = true
	}
	if order.Takeaway == nil {
		err.noTakeaway = true
	}

	s.mu.Lock()
	delete(s.orders, u.ID)
	s.mu.Unlock()

	log.Debugf("finishing order: %+v", order)

	if err.OrderNotComplete() {
		return nil, false
	}

	return order, true
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
