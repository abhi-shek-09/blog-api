package services

import (
	"blog-api/models"
	"blog-api/repository"
	"errors"
	"fmt"
)

type PostService struct {
	PostRepository *repository.PostRepository
}

func NewPostService(repo *repository.PostRepository) *PostService {
    return &PostService{
        PostRepository: repo,
    }
}

func (s *PostService) CreatePost(post *models.Post) (int, error) {
	if post.Title == "" {
		return 0, errors.New("title cannot be empty")
	}

	if post.Content == "" {
		return 0, errors.New("content cannot be empty")
	}

	return s.PostRepository.CreatePost(post)
}

func (s *PostService) FetchPost(id int) (*models.Post, error){
	if id <= 0 {
		return nil, errors.New("invalid post ID")
	}
	post, err := s.PostRepository.FetchPost(id)
	if err != nil {
		return nil, err
	}

	if post == nil {
		return nil, errors.New("post not found")
	}

	return post, nil
}

func (s *PostService) FetchPosts(term string) ([]*models.Post, error){
	var posts []*models.Post
	var err error

	posts, err = s.PostRepository.FetchPosts(term)
	if posts == nil{
		return []*models.Post{}, fmt.Errorf("error fetching posts: %v", err)
	}
	return posts, nil
}

func (s *PostService) UpdatePost(post *models.Post) error{
	if post.ID <= 0 {
        return fmt.Errorf("invalid post ID")
    }
	err := s.PostRepository.UpdatePost(post)
	if err != nil {
		return fmt.Errorf("error updating the post :%v", err)
	}
	return nil
}

func (s *PostService) DeletePost(id int) error{
	if id <= 0 {
        return fmt.Errorf("invalid post ID")
    }
	err := s.PostRepository.DeletePost(id)
	if err != nil {
		return fmt.Errorf("error deleting the post :%v", err)
	}
	return nil
}
