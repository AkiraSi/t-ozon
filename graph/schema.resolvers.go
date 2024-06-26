package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"t-ozon/graph/model"
	"time"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*model.Post, error) {
	if r.Pool == nil {
		user, err := r.DB.GetUserByID(ctx, input.UserID)
		if err != nil {
			return nil, fmt.Errorf("fail (user): %w", err)
		}
		post := &model.Post{
			ID:         strconv.Itoa(r.DB.NextID),
			Title:      input.Title,
			Content:    input.Content,
			Date:       time.Now(),
			Commenting: input.Commenting,
			UserID:     input.UserID,
			User:       user,
		}
		err = r.DB.CreatePost(ctx, post)
		if err != nil {
			return nil, fmt.Errorf("save (post): %w", err)
		}
		return post, nil
	} else {
		var count int
		err := r.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM public.post").Scan(&count)
		if err != nil {
			return nil, fmt.Errorf("create (post, count): %w", err)
		}
		post := &model.Post{
			ID:         strconv.Itoa(count + 1),
			Title:      input.Title,
			Content:    input.Content,
			Date:       time.Now(),
			Commenting: input.Commenting,
			UserID:     input.UserID,
		}
		_, err = r.Pool.Exec(ctx, "INSERT INTO public.post(id, title, content, \"Date\", user_id, commenting) VALUES ($1, $2, $3, $4, $5, $6)",
			post.ID, post.Title, post.Content, post.Date, post.UserID, strconv.FormatBool(post.Commenting))
		if err != nil {
			return nil, fmt.Errorf("create (post): %w", err)
		}
		return post, nil
	}
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	if r.Pool == nil {
		user := &model.User{
			ID:   input.ID,
			Name: input.Name,
		}
		err := r.DB.CreateUser(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("save (post): %w", err)
		}
		return user, nil
	} else {
		user := &model.User{
			ID:   input.ID,
			Name: input.Name,
		}
		_, err := r.Pool.Exec(ctx, "INSERT INTO public.\"user\"(id, name) VALUES ($1, $2)", user.ID, user.Name)
		if err != nil {
			return nil, fmt.Errorf("create (user): %w", err)
		}
		return user, nil
	}
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input model.NewComment) (*model.Comment, error) {
	if r.Pool == nil {
		user, err := r.DB.GetUserByID(ctx, input.UserID)
		if err != nil {
			return nil, fmt.Errorf("find (user): %w", err)
		}
		post, err := r.DB.GetPostByID(ctx, input.PostID)
		if err != nil {
			return nil, fmt.Errorf("find (post): %w", err)
		}
		comment := &model.Comment{
			Content: input.Content,
			User:    user,
			Date:    time.Now(),
			PostID:  input.PostID,
			ID:      strconv.Itoa(r.DB.NextID),
			UserID:  input.UserID,
		}
		if len(comment.Content) > 2000 {
			return nil, fmt.Errorf("2k")
		}
		if !post.Commenting {
			return nil, fmt.Errorf("Nelzya")
		}
		err = r.DB.CreateComment(ctx, comment)
		if err != nil {
			return nil, fmt.Errorf("save (comment): %w", err)
		}
		err = r.DB.UpdatePost(ctx, post)
		if err != nil {
			return nil, fmt.Errorf("update (comment): %w", err)
		}
		post.Comments = append(post.Comments, comment)
		r.ChatComments = append(r.ChatComments, comment)
		return comment, nil
	} else {
		var exists bool
		r.Pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM post WHERE id = $1)", input.PostID).Scan(&exists)
		if !exists {
			return nil, fmt.Errorf("no post")
		}
		var count int
		err := r.Pool.QueryRow(ctx, "SELECT count(*) FROM comment WHERE post_id = $1;", input.PostID).Scan(&count)
		if err != nil {
			return nil, fmt.Errorf("create (comment): %w", err)
		}
		comment := &model.Comment{
			Content: input.Content,
			Date:    time.Now(),
			PostID:  input.PostID,
			ID:      strconv.Itoa(count + 1),
			UserID:  input.UserID,
		}
		if len(comment.Content) > 2000 {
			return nil, fmt.Errorf("2k")
		}
		var commenting string
		r.Pool.QueryRow(ctx, "SELECT Commenting FROM post WHERE id = $1", input.PostID).Scan(&commenting)
		if commenting == "false" {
			return nil, fmt.Errorf("Nelzya")
		}
		_, err = r.Pool.Exec(ctx, "INSERT INTO public.comment(content, \"Date\", user_id, post_id, id) VALUES ($1, $2, $3, $4, $5)",
			comment.Content, comment.Date, comment.UserID, comment.PostID, comment.ID)
		if err != nil {
			return nil, fmt.Errorf("ne vishlo: %w", err)
		}
		r.ChatComments = append(r.ChatComments, comment)
		return comment, nil
	}
}

func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	if r.Pool == nil {
		posts, err := r.DB.GetPosts(ctx)
		if err != nil {
			return nil, fmt.Errorf("get (posts): %w", err)
		}
		return posts, nil
	} else {
		rows, err := r.Pool.Query(ctx, "SELECT id, title, content, commenting, \"Date\", user_id FROM public.post")
		if err != nil {
			return nil, fmt.Errorf("get rows (posts): %w", err)
		}
		var posts []*model.Post
		for rows.Next() {
			var post model.Post
			var date string
			var commenting string
			err := rows.Scan(&post.ID, &post.Title, &post.Content, &commenting, &date, &post.UserID)
			if err != nil {
				return nil, fmt.Errorf("scan (post): %w", err)
			}
			newDate := strings.Split(date, " ")
			newTime := newDate[0] + " " + newDate[1]
			post.Date, err = time.Parse("2006-01-02 15:04:05.999999999", newTime) // Парсинг даты в нужном формате
			if err != nil {
				return nil, fmt.Errorf("parse (post): %w", err)
			}
			post.Commenting = strconv.CanBackquote(commenting)
			posts = append(posts, &post)
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("iteraite fail: %w", err)
		}
		return posts, nil
	}
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	posts, err := r.DB.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail (user): %w", err)
	}
	return posts, nil
}

func (r *queryResolver) Comments(ctx context.Context, postID string, limit *int, offset *int) ([]*model.Comment, error) {
	var limitArg int = 10
	if limit != nil {
		limitArg = *limit
	} else {
		limitArg = 10
	}
	var offsetArg int = 0
	if offset != nil {
		offsetArg = *offset
	}
	if r.Pool == nil {
		comments, err := r.DB.GetComments(ctx, postID)
		if err != nil {
			return nil, fmt.Errorf("get (comment): %w", err)
		}
		if len(comments) > offsetArg+limitArg {
			comments = comments[offsetArg : offsetArg+limitArg]
		} else if len(comments) > offsetArg {
			comments = comments[offsetArg:]
		} else {
			comments = comments[:0]
		}
		return comments, nil
	} else {
		rows, err := r.Pool.Query(ctx, "SELECT id, content, \"Date\", user_id, post_id FROM public.comment WHERE post_id = $1 ORDER BY \"Date\" DESC LIMIT $2 OFFSET $3", postID, limitArg, offsetArg)
		if err != nil {
			return nil, fmt.Errorf("get rows (comment): %w", err)
		}
		defer rows.Close()
		var comments []*model.Comment
		for rows.Next() {
			var comment model.Comment
			var date string
			err := rows.Scan(&comment.ID, &comment.Content, &date, &comment.UserID, &comment.PostID)
			if err != nil {
				return nil, fmt.Errorf("scan (comment): %w", err)
			}
			newDate := strings.Split(date, " ")
			newTime := newDate[0] + " " + newDate[1]
			comment.Date, err = time.Parse("2006-01-02 15:04:05.999999999", newTime)
			if err != nil {
				return nil, fmt.Errorf("parse data(comment): %w", err)
			}
			comments = append(comments, &comment)
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("fail rows(comment): %w", err)
		}
		return comments, nil
	}
}

func (r *subscriptionResolver) NewComment(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	ch := make(chan *model.Comment)
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				r.mu.Lock()
				var latestComment *model.Comment
				for _, comment := range r.ChatComments {
					if comment.PostID == postID {
						latestComment = comment
					}
				}
				if latestComment != nil {
					select {
					case ch <- latestComment:
						r.ChatComments = nil
					case <-ctx.Done():
						return
					}
				}
				r.mu.Unlock()
				time.Sleep(3 * time.Second)
			}
		}
	}()
	return ch, nil
}

func (r *Resolver) Mutation() MutationResolver         { return &mutationResolver{r} }
func (r *Resolver) Query() QueryResolver               { return &queryResolver{r} }
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
