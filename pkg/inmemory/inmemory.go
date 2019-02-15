package inmemory

import (
	"fmt"
	"sync"

	"github.com/gerlacdt/prom_example/pkg/domain"
)

// PostService implements the corresponding service
type PostService struct {
	posts  map[int]*domain.Post
	nextID int
	mutex  *sync.Mutex
}

// New create a new PostService
func New() *PostService {
	posts := make(map[int]*domain.Post)
	return &PostService{posts: posts, nextID: 1, mutex: &sync.Mutex{}}
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
func (service *PostService) CreatePost(m *domain.Message) (*domain.Post, error) {
	service.mutex.Lock()
	defer service.mutex.Unlock()
	p := &domain.Post{ID: service.nextID, Message: *m}
	service.nextID++
	service.posts[p.ID] = p
	return p, nil
}

// DeletePost ...
func (service *PostService) DeletePost(id int) error {
	delete(service.posts, id)
	return nil
}
