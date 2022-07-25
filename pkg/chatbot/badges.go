package chatbot

func (bot *ChatBot) hasBadges(badges map[string]int, acceptable []string) bool {
	for _, key := range acceptable {
		if _, found := badges[key]; found {
			return true
		}
	}
	return false
}

func (bot *ChatBot) IsModerator(badges map[string]int) bool {
	return bot.hasBadges(badges, []string{"broadcaster", "moderator"})
}

func (bot *ChatBot) IsVIP(badges map[string]int) bool {
	acceptable := []string{"vip"}
	return bot.IsModerator(badges) || bot.hasBadges(badges, acceptable)
}

func (bot *ChatBot) IsSubscriber(badges map[string]int) bool {
	acceptable := []string{"subscriber", "founder", "premium"}
	return bot.IsModerator(badges) || bot.hasBadges(badges, acceptable)
}
