package chatbot

// HasBadges checks if the user has any of the acceptable badges.
func (bot *ChatBot) HasBadges(badges map[string]int, acceptable []string) bool {
	for _, key := range acceptable {
		if _, found := badges[key]; found {
			return true
		}
	}
	return false
}

// IsModerator checks if the user is a moderator or broadcaster
func (bot *ChatBot) IsModerator(badges map[string]int) bool {
	return bot.HasBadges(badges, []string{"broadcaster", "moderator"})
}

// IsVIP checks if the user is a VIP or above
func (bot *ChatBot) IsVIP(badges map[string]int) bool {
	acceptable := []string{"vip"}
	return bot.IsModerator(badges) || bot.HasBadges(badges, acceptable)
}

// IsSubscriber checks if the user is a subscriber (founder, subscriber, premium) or above (mod, broadcaster)
func (bot *ChatBot) IsSubscriber(badges map[string]int) bool {
	acceptable := []string{"subscriber", "founder", "premium"}
	return bot.IsModerator(badges) || bot.HasBadges(badges, acceptable)
}
