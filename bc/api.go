package bc

import "time"

// Client is Beancounter Platform API client.
type Client interface {
	ListUsers(c *Cred) ([]*User, error)
	GetUserProfile(c *Cred, userId string) (*UserProfile, error)
	ListActivities(c *Cred, userId string, page int) ([]*Activity, error)
}

// Cred is a credentials object identifying a single customer.
type Cred struct {
	Customer string
	ApiKey   string
}

// User is a Beancounter user.
type User struct {
	Id       string         `json:"id"`
	Email    string         `json:"email,omitempty"`
	Name     string         `json:"name,omitempty"`
	Photo    string         `json:"photo,omitempty"`
	About    string         `json:"about,omitempty"`
	Age      int            `json:"age,omitempty"`
	Gender   string         `json:"gender,omitempty"`
	Location string         `json:"location,omitempty"`
	Services []*UserService `json:"services,omitempty"`
}

// UserService is a third-party social network service.
type UserService struct {
	Name  string `json:"name"`
	Id    string `json:"id"`
	Link  string `json:"link"`
	Photo string `json:"photo,omitempty"`
}

// UserProfile is profile of a single user.
type UserProfile struct {
	Updated time.Time `json:"updated"`
	Topics  []*Topic  `json:"topics"`
}

// Topic identifies a single category or interest (typically of a user).
type Topic struct {
	Kind       string   `json:"kind"`
	Resource   string   `json:"resource"`
	Label      string   `json:"label"`
	Weight     float32  `json:"weight"`
	Urls       []string `json:"urls,omitempty"`
	Activities []string `json:"activities,omitempty"`
}

// Activity is a single user activity.
type Activity struct {
	Id          string    `json:"id"`
	Verb        string    `json:"verb"`
	Kind        string    `json:"kind"`
	Service     string    `json:"service"`
	Timestamp   time.Time `json:"ts"`
	Url         string    `json:"url,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"desc,omitempty"`
}
