package models

import "time"

type Book struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Author     string    `json:"author"`
	Published  int       `json:"published,omitempty"`
	ISBN       string    `json:"isbn,omitempty"`
	CheckedOut bool      `json:"checked_out"`
	DueDate    time.Time `json:"due_date,omitempty"`
}

type Loan struct {
	BookID   string    `json:"book_id"`
	UserID   string    `json:"user_id"`
	DueDate  time.Time `json:"due_date"`
	Returned bool      `json:"returned"`
}

var (
	Books = make(map[string]Book)
	Loans = make(map[string]Loan)
)
