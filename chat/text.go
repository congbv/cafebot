package chat

import (
	"bytes"

	"github.com/yarikbratashchuk/cafebot/order"
)

var (
	buttonText = map[string]string{
		"new_order":     "📝 Сделать заказ",
		"here":          "🏠 В кафе",
		"takeaway":      "🚶‍♂️ На вынос",
		"back":          "⬅️ Назад",
		"preview_order": "🛒 Посмотреть свой заказ",
		"send_order":    "💌 Отправить свой заказ",
		"selected":      "✅",
	}

	text = map[string]string{
		"help":  "Лучше позвоните им по телефону, так будет быстрее 😊",
		"start": "Привет 🖖. Я готов принять твой заказ 👌",
		"wrong": "Сделайте заказ по телефону, я не понимаю вас 😞",

		"where?":        "В кафе или на вынос?",
		"when?":         "На который час готовить?",
		"what?":         "Что вы будете есть?",
		"preview_order": "Посмотрите свой заказ.",

		"time_preview_prefix":     "🕑",
		"order_preview_prefix":    "🥘",
		"here_preview_prefix":     "🏠",
		"takeaway_preview_prefix": "🚶‍♂️",
	}
)

func generatePreviewOrderText(o order.Order) string {
	buf := new(bytes.Buffer)

	if *o.Takeaway {
		buf.WriteString(text["here_preview_prefix"])
		buf.WriteRune(' ')
		buf.WriteString(buttonText["here"])
	} else {
		buf.WriteString(text["here_preview_prefix"])
		buf.WriteRune(' ')
		buf.WriteString(buttonText["takeaway"])
	}
	buf.WriteRune('\n')

	buf.WriteString(text["time_preview_prefix"])
	buf.WriteRune(' ')
	buf.WriteString(o.Time.Format("15:04"))
	buf.WriteRune('\n')

	buf.WriteString(text["order_preview_prefix"])
	buf.WriteRune(' ')
	buf.WriteString(o.Meal[0])
	buf.WriteRune('\n')
	if len(o.Meal) > 1 {
		for i := 1; i < len(o.Meal); i++ {
			buf.WriteString("   ")
			buf.WriteString(o.Meal[0])
			buf.WriteRune('\n')
		}
	}

	return buf.String()
}
