package connection

import (
	"database/sql/driver"
	"encoding/json"
)

func (ics *InactiveChatSessions) Scan(src interface{}) error {
	return json.Unmarshal([]byte(src.(string)), &ics)
}

func (ics *InactiveChatSessions) Value() (driver.Value, error) {
	val, err := json.Marshal(ics.ChatSessionsId)
	return string(val), err
}
