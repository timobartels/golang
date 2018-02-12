package model

// Person defines the structure of a person record
type Person struct {
	ID        string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
}

// People stores the person records
type People struct {
	Store []Person
}
