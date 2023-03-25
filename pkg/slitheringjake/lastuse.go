package slitheringjake

import "time"

// updateLastUse updates the last use time for the given user
func (jake *SlitheringJake) updateLastUse(user string) {
	jake.Bot.Lock(lastuse_mutex)
	defer jake.Bot.Unlock(lastuse_mutex)

	jake.lastUse[user] = time.Now()
}

// deleteLastUse deletes the last use time for the given user
func (jake *SlitheringJake) deleteLastUse(user string) {
	jake.Bot.Lock(lastuse_mutex)
	defer jake.Bot.Unlock(lastuse_mutex)

	delete(jake.lastUse, user)
}

func (jake *SlitheringJake) checkLastUse(user string) bool {
	jake.Bot.Lock(lastuse_mutex)
	defer jake.Bot.Unlock(lastuse_mutex)

	lastUse, found := jake.lastUse[user]
	if !found {
		return true
	}

	useInterval := time.Duration(jake.Config.UseInterval) * time.Second
	return time.Since(lastUse) >= useInterval
}
