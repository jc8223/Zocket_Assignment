 
package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Mobile    int    `json:"mobile"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}
