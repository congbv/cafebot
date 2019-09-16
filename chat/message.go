package chat

var (
	buttonText = map[string]string{
		"new_order": "Сделать заказ",
		"here":      "В кафе за столиком",
		"takeaway":  "На вынос с собой",

		"step_back": "Назад",
	}

	message = map[string]string{
		"help":  "Лучше позвоните им по телефону, так будет быстрее 😊",
		"start": "Привет 🖖. Я готов принять твой заказ 👌",
		"wrong": "Сделайте заказ по телефону, я не понимаю вас 😞",

		"where?": "В кафе или на вынос?",
		"when?":  "На который час готовить?",
		"what?":  "Что вы будете есть и пить?",
	}

	selectResponse = map[string]string{
		"takeaway_selected": "Выбрано: на вынос",
		"here_selected":     "Выбрано: в кафе",
	}
)
