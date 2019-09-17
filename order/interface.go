// order holds user order related stuff
package order

import (
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Service describes order creation operations
type Service interface {
	InitOrder(*api.User) *Order
	AddMeal(*api.User, string) *Order
	RemoveMeal(*api.User, string) *Order
	SetTime(*api.User, time.Time) *Order
	SetTakeaway(*api.User, bool) *Order
	FinishOrder(*api.User) (*Order, error)
}
