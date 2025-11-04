package grist

type AccessUser struct {
	Users []User `json:"users"`
}

type User struct {
	ID      int        `json:"id"`
	Name    string     `json:"name"`
	Email   string     `json:"email,omitempty"`
	Access  AccessRole `json:"access,omitempty"`
	Picture string     `json:"picture,omitempty"`
	Ref     string     `json:"ref,omitempty"`
	Member  bool       `json:"isMember,omitempty"`
}
