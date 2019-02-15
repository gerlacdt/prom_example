package domain

// Message contains the contents of a blog post
type Message struct {
	Body string `json:"body"`
}

// Post entry in the database or redis or in-memory
type Post struct {
	Message
	ID int `json:"id"`
}

// PostService for creating, finding and deleting posts
type PostService interface {
	Post(id int) (*Post, error)
	Posts() ([]*Post, error)
	CreatePost(m *Message) (*Post, error)
	DeletePost(id int) error
}
