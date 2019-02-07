package domain

// Post entry in the database or redis or in-memory
type Post struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// PostService for creating, finding and deleting posts
type PostService interface {
	Post(id int) (*Post, error)
	Posts() ([]*Post, error)
	CreatePost(p *Post) error
	DeletePost(id int) error
}
