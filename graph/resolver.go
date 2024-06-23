package graph

import (
	"context"
	"fmt"
	"os"
	"sync"
	"t-ozon/graph/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

//go:generate mockgen -package=your_package_name -destination=mock_resolver.go -source=resolver.go

type Resolver struct {
	DB           InMemoryDB
	ChatComments []*model.Comment
	mu           sync.Mutex
	Pool         *pgxpool.Pool
}

func NewResolver(typeMemory string) (*Resolver, error) {
	resolver := &Resolver{}
	var err error
	if typeMemory == "true" {
		config := map[string]string{
			"host":     os.Getenv("POSTGRES_HOST"),
			"port":     os.Getenv("POSTGRES_PORT"),
			"user":     os.Getenv("POSTGRES_USER"),
			"password": os.Getenv("POSTGRES_PASSWORD"),
			"database": os.Getenv("POSTGRES_DB"),
		}
		resolver.Pool, err = ConnectToPostgres(config)
		if err != nil {
			return nil, fmt.Errorf("не удалось приконнектиться: %w", err)
		}
	} else {
		NewInMemoryDB()
	}
	resolver.ChatComments = make([]*model.Comment, 0)
	resolver.mu = sync.Mutex{}
	return resolver, nil
}

type Subscription struct {
	NewComment <-chan *model.Comment
}

type PostgresDB struct {
	pool *pgxpool.Pool
}

type InMemoryDB struct {
	Posts    map[string]*model.Post
	Users    map[string]*model.User
	Comments map[string]*model.Comment
	NextID   int
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		Posts:    make(map[string]*model.Post),
		Users:    make(map[string]*model.User),
		Comments: make(map[string]*model.Comment),
		NextID:   1,
	}
}

func (db *InMemoryDB) CreatePost(ctx context.Context, post *model.Post) error {
	if db.Posts == nil {
		db.Posts = make(map[string]*model.Post)
	}
	post.ID = fmt.Sprintf("%d", db.NextID)
	db.NextID++
	db.Posts[post.ID] = post
	return nil
}

func (db *InMemoryDB) CreateUser(ctx context.Context, User *model.User) error {
	if db.Users == nil {
		db.Users = make(map[string]*model.User)
	}
	db.Users[User.ID] = User
	return nil
}

func (db *InMemoryDB) GetPosts(ctx context.Context) ([]*model.Post, error) {
	posts := make([]*model.Post, 0, len(db.Posts))
	for _, post := range db.Posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (db *InMemoryDB) GetUsers(ctx context.Context) ([]*model.User, error) {
	Users := make([]*model.User, 0, len(db.Users))
	for _, User := range db.Users {
		Users = append(Users, User)
	}
	return Users, nil
}

func (db *InMemoryDB) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	if db.Users != nil {
		user, ok := db.Users[userID]
		if !ok {
			return nil, fmt.Errorf("юзер не найден: %s", userID)
		}
		return user, nil
	}
	return nil, nil
}

func (db *InMemoryDB) GetPostByID(ctx context.Context, postID string) (*model.Post, error) {
	post, ok := db.Posts[postID]
	if !ok {
		return nil, fmt.Errorf("пост не найден: %s", postID)
	}
	return post, nil
}

func (db *InMemoryDB) CreateComment(ctx context.Context, comment *model.Comment) error {
	if db.Comments == nil {
		db.Comments = make(map[string]*model.Comment)
	}
	db.Comments[comment.ID] = comment
	return nil
}

func (db *InMemoryDB) UpdatePost(ctx context.Context, post *model.Post) error {
	db.Posts[post.ID] = post
	return nil
}

func (db *InMemoryDB) GetComments(ctx context.Context, postID string) ([]*model.Comment, error) {
	comments := make([]*model.Comment, 0, len(db.Posts[postID].Comments))
	for _, comment := range db.Posts[postID].Comments {
		comments = append(comments, comment)
	}
	return comments, nil
}
