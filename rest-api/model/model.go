package model

// Person struct to store the data for the API
type Person struct {
	ID        string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
}

type People struct {
	Store []Person
}

func NewStore() *People {
	people := &People{}
	return people
}
