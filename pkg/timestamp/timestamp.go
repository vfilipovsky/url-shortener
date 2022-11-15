package timestamp

import "time"

type Timestamp struct {
	// createdAt - unix format datetime
	CreatedAt time.Time `json:"created_at"`

	// updatedAt - unix format datetime
	UpdatedAt time.Time `json:"updated_at"`
}
