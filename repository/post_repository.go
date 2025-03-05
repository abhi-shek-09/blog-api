package repository

// CRUD Operations happen here, Encapsulates database queries
// Uses database/sql to interact with PostgreSQL

import (
	"blog-api/database"
	"blog-api/models"
	"database/sql"
	"fmt"
)

type PostRepository struct {
	DB *sql.DB
}

func NewRepository() *PostRepository {
	return &PostRepository{
		DB: database.DB,
	}
}

func (r *PostRepository) CreatePost(post *models.Post) (int, error) {
	query := `INSERT INTO posts (title, content, category, tags, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`
	// Uses placeholders ($1, $2, $3, $4) for safe parameterized queries.
	var id int
	err := r.DB.QueryRow(query, post.Title, post.Content, post.Category, post.Tags).Scan(&id)
	// .Scan(&id) to store the new post’s ID in id using the pointer to that variable
	if err != nil {
		return 0, fmt.Errorf("failed to create post: %v", err)
	}
	return id, nil
}

func (r *PostRepository) FetchPost(id int) (*models.Post, error) {
	query := `SELECT id, title, content, category, tags, created_at, updated_at FROM posts WHERE id = $1`

	post := &models.Post{}
	err := r.DB.QueryRow(query, id).Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.Tags, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to access the post: %v", err)
	}
	return post, nil
}

func (r *PostRepository) FetchPosts(term string) ([]*models.Post, error) {
	query := `SELECT id, title, content, category, tags, created_at, updated_at FROM posts`
	var rows *sql.Rows
	var err error

	if term != "" {
		query += ` WHERE title ILIKE $1 OR content ILIKE $1 OR category ILIKE $1`
		rows, err = r.DB.Query(query, "%"+term+"%")
	} else {
		rows, err = r.DB.Query(query)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %v", err)
	}
	defer rows.Close() // very imp

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.Tags, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				return []*models.Post{}, nil
			}
			return []*models.Post{}, fmt.Errorf("failed to access the post: %v", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepository) UpdatePost(post *models.Post) error {
	query := `UPDATE posts SET title=$1, content=$2, category=$3, tags=$4, updated_at=NOW() WHERE id=$5`

	result, err := r.DB.Exec(query, post.Title, post.Content, post.Category, post.Tags, post.ID)
	if err != nil {
		return fmt.Errorf("failed to update post: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("post not found")
	}
	return nil
}

func (r *PostRepository) DeletePost(id int) error {
	query := `DELETE FROM posts WHERE id=$1`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("post not found")
	}
	return nil
}
