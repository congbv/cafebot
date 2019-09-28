package chat

import (
	"bytes"
	"fmt"

	"cafebot/order"
)

var (
	buttonText = map[string]string{
		"new":      "📝 Make an order",
		"here":     "🏠 In the cafe",
		"takeaway": "🚶‍♂️ Takeaway",
		"back":     "<< Back",
		"preview":  "🛒 My order",
		"send":     "📨 submit",
		"selected": "✅",
	}

	text = map[string]string{
		"help":    "Better call them on the phone, it will be faster 😊",
		"start":   "Hello 🖖. I am ready to accept your order 👌",
		"wrong":   "Make an order by phone, I do not understand you 😞",
		"where?":  "In a cafe or takeaway?",
		"when?":   "What time to cook?",
		"what?":   "What will you eat?",
		"preview": "Check your order then click 'Submit'",
		"sent":    "Your order has been successfully sent to the cafe 😉",

		"err_internal":    "An error has occurred 😞 make an order by phone please",
		"err_no_username": "Add your username in the settings and try again",

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
	buf.WriteRune('@')
	buf.WriteString(o.User.UserName)
	return buf.String()
}

func boldText(text string) string {
	return fmt.Sprintf("<b>%s</b>", text)
}
