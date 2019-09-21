package chat

import (
	"bytes"
	"fmt"

	"github.com/yarikbratashchuk/cafebot/order"
)

var (
	buttonText = map[string]string{
		"new":      "📝 Сделать заказ",
		"here":     "🏠 В кафе",
		"takeaway": "🚶‍♂️ На вынос",
		"back":     "<< Назад",
		"preview":  "🛒 Мой заказ",
		"send":     "📨 Отправить",
		"selected": "✅",
	}

	text = map[string]string{
		"help":  "Лучше позвоните им по телефону, так будет быстрее 😊",
		"start": "Привет 🖖. Я готов принять твой заказ 👌",
		"wrong": "Сделайте заказ по телефону, я не понимаю вас 😞",

		"where?":  "В кафе или на вынос?",
		"when?":   "На который час готовить?",
		"what?":   "Что вы будете есть?",
		"preview": "Проверьте свой заказ затем нажмите 'Отправить'",
		"sent":    "Ваш заказ успешно отправлен в кафе 😉",

		"error": "Произошла ошибка 😞 cделайте заказ по телефону пожалуйста",

		"time_preview_prefix":     "🕑",
		"order_preview_prefix":    "🥘",
		"here_preview_prefix":     "🏠",
		"takeaway_preview_prefix": "🚶‍♂️",
	}
)

func generatePreviewText(o order.Order) string {
	buf := new(bytes.Buffer)

	if *o.Takeaway {
		buf.WriteString(buttonText["takeaway"])
	} else {
		buf.WriteString(buttonText["here"])
	}
	buf.WriteString("\n")

	buf.WriteString(text["time_preview_prefix"])
	buf.WriteString(" ")
	buf.WriteString(o.Time.Format("15:04"))
	buf.WriteString("\n")

	buf.WriteString(text["order_preview_prefix"])
	buf.WriteString(" ")
	buf.WriteString(o.Meal[0])
	buf.WriteString("\n")
	if len(o.Meal) > 1 {
		for i := 1; i < len(o.Meal); i++ {
			buf.WriteString("      ")
			buf.WriteString(o.Meal[i])
			buf.WriteRune('\n')
		}
	}

	return buf.String()
}

func generateUserNameText(o order.Order) string {
	buf := new(bytes.Buffer)
	buf.WriteString("<b>")
	buf.WriteString(o.User.FirstName)
	buf.WriteRune(' ')
	buf.WriteString(o.User.LastName)
	buf.WriteString("</b>")
	buf.WriteRune('\n')
	if o.User.UserName != "" {
		buf.WriteRune('@')
		buf.WriteString(o.User.UserName)
		buf.WriteRune('\n')
	}
	buf.WriteRune('\n')
	return buf.String()
}

func boldText(text string) string {
	return fmt.Sprintf("<b>%s</b>", text)
}
