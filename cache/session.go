package cache

import "time"

// session inherits the expiry struct and is embedded in types UserSession and MerchantSession.
type session struct {
	uuid string
	expiry
}

//getKey returns the stored session UUID.
func (s *session) getKey() string {
	return s.uuid
}

// expiry is embedded in types session (tracks active user and merchant sessions) and cachedItem (tracks products stored in cache).
type expiry struct {
	expiry time.Time
}

//monitor sleeps for the difference between the expiration time and the current time.
//eg expires at 4pm, current time is 4.30pm. it sleeps for 30mins.
//after sleep ends, it checks if the expiration time has exceeded the current time it returns.
//or else it loops.
func (e *expiry) monitor() {
	for {
		sleeptime := time.Until(e.expiry)
		time.Sleep(sleeptime)
		if e.expiry.Before(time.Now()) {
			return
		}
	}
}

//updateExpiryTime updates the session's expiry time.
func (e *expiry) updateExpiryTime(updatedTime time.Time) {
	e.expiry = updatedTime
}
