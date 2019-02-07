package inmemory

import (
	"fmt"

	"github.com/gerlacdt/prom_example/pkg/domain"
)

// PostService implements the corresponding service
type PostService struct {
	posts map[int]*domain.Post
}

// Post function finds the post with the given id
func (service *PostService) Post(id int) (*domain.Post, error) {
	if post, ok := service.posts[id]; ok {
		return post, nil
	}
	return nil, fmt.Errorf("Find Post: given id: %d does not exist", id)
}

// Posts returns all existing posts
func (service *PostService) Posts() ([]*domain.Post, error) {
	var posts []*domain.Post
	for _, p := range service.posts {
		posts = append(posts, p)
	}
	return posts, nil
}

// CreatePost creates a post
func (service *PostService) CreatePost(p *domain.Post) error {
	service.posts[p.ID] = p
	return nil
}

// DeletePost ...
func (service *PostService) DeletePost(id int) error {
	delete(service.posts, id)
	return nil
}
