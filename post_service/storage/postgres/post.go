package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	p "github.com/microservice/post_service/genproto/post"
	"github.com/opentracing/opentracing-go"
)

func (r *PostRepo) CreatePost(ctx context.Context, post *p.PostRequest) (*p.PostResponse, error) {
	trace, ctx := opentracing.StartSpanFromContext(ctx, "UpdateAddress")
	defer trace.Finish()
	var res p.PostResponse
	err := r.db.QueryRow(`
		insert into 
			posts(title, description, user_id) 
		values
			($1, $2, $3) 
		returning 
			id, title, description, likes, user_id, created_at, updated_at`, post.Title, post.Description, post.UserId).Scan(&res.Id, &res.Title, &res.Description, &res.Likes, &res.UserId, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to create post")
		return &p.PostResponse{}, err
	}

	return &res, nil
}

func (r *PostRepo) GetPostById(ctx context.Context, post *p.IdRequest) (*p.PostResponse, error) {
	trace, ctx := opentracing.StartSpanFromContext(ctx, "UpdateAddress")
	defer trace.Finish()
	res := p.PostResponse{}
	err := r.db.QueryRow(`
		select 
			id, title, description, likes, user_id, created_at, updated_at 
		from 
			posts 
		where 
			id = $1 and deleted_at is null`, post.Id).Scan(&res.Id, &res.Title, &res.Description, &res.Likes, &res.UserId, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to get post")
		return &p.PostResponse{}, err
	}

	return &res, nil
}

func (r *PostRepo) GetPostByUserId(ctx context.Context, id *p.IdRequest) (*p.Posts, error) {
	trace, ctx := opentracing.StartSpanFromContext(ctx, "UpdateAddress")
	defer trace.Finish()
	res := p.Posts{}
	rows, err := r.db.Query(`
		select 
			id, title, description, likes, user_id, created_at, updated_at 
		from 
			posts 
		where 
			user_id = $1 and deleted_at is null`, id.Id)

	if err != nil {
		log.Println("failed to get post by user_id")
		return &p.Posts{}, err
	}

	for rows.Next() {
		post := p.PostResponse{}

		err = rows.Scan(
			&post.Id,
			&post.Title,
			&post.Description,
			&post.Likes,
			&post.UserId,
			&post.CreatedAt,
			&post.UpdatedAt,
		)

		if err != nil {
			log.Println("failed to scanning post")
			return &p.Posts{}, err
		}

		res.Posts = append(res.Posts, &post)
	}

	return &res, nil
}

func (r *PostRepo) GetPostForUser(ctx context.Context, id *p.IdRequest) (*p.Posts, error) {
	trace, ctx := opentracing.StartSpanFromContext(ctx, "UpdateAddress")
	defer trace.Finish()
	res := p.Posts{}
	rows, err := r.db.Query(`
		select 
			id, title, description, likes, user_id, created_at, updated_at 
		from 
			posts 
		where
			user_id = $1`, id.Id)

	if err != nil {
		log.Println("failed to get post for user")
		return &p.Posts{}, err
	}

	for rows.Next() {
		post := p.PostResponse{}

		err = rows.Scan(
			&post.Id,
			&post.Title,
			&post.Description,
			&post.Likes,
			&post.UserId,
			&post.CreatedAt,
			&post.UpdatedAt,
		)

		if err != nil {
			log.Println("failed to scanning post")
			return &p.Posts{}, nil
		}

		res.Posts = append(res.Posts, &post)
	}

	return &res, nil
}

func (r *PostRepo) GetPostForComment(ctx context.Context, post *p.IdRequest) (*p.PostResponse, error) {
	trace, ctx := opentracing.StartSpanFromContext(ctx, "UpdateAddress")
	defer trace.Finish()
	res := p.PostResponse{}
	err := r.db.QueryRow(`
		select 
			id, title, description, likes, user_id, created_at, updated_at 
		from 
			posts 
		where 
			id = $1 and deleted_at is null`, post.Id).Scan(&res.Id, &res.Title, &res.Description, &res.Likes, &res.UserId, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to get post")
		return &p.PostResponse{}, err
	}

	return &res, nil
}

func (r *PostRepo) SearchByTitle(ctx context.Context, title *p.Search) (*p.Posts, error) {
	trace, ctx := opentracing.StartSpanFromContext(ctx, "UpdateAddress")
	defer trace.Finish()
	res := p.Posts{}
	query := fmt.Sprint("select id, title, description, likes, user_id, created_at, updated_at from posts where title ilike '%" + title.Search + "%' and deleted_at is null")

	rows, err := r.db.Query(query)
	if err != nil {
		log.Println("failed to search post")
		return &p.Posts{}, nil
	}

	for rows.Next() {
		post := p.PostResponse{}

		err = rows.Scan(
			&post.Id,
			&post.Title,
			&post.Description,
			&post.Likes,
			&post.UserId,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			log.Println("failed to scanning post")
			return &p.Posts{}, nil
		}

		res.Posts = append(res.Posts, &post)
	}

	return &res, nil
}

func (r *PostRepo) LikePost(ctx context.Context, l *p.LikeRequest) (*p.PostResponse, error) {
	trace, ctx := opentracing.StartSpanFromContext(ctx, "UpdateAddress")
	defer trace.Finish()
	res := p.PostResponse{}
	if l.IsLiked {
		err := r.db.QueryRow(`
			update 
				posts 
			set 
				likes = likes + 1 
			where 
				id = $1 
			returning 
				id, title, description, likes, user_id, created_at, updated_at`, l.PostId).Scan(&res.Id, &res.Title, &res.Description, &res.Likes, &res.UserId, &res.CreatedAt, &res.UpdatedAt)
		if err != nil {
			log.Println("failed to like post")
			return &p.PostResponse{}, err
		}
	} else {
		err := r.db.QueryRow(`
			select 
				id, title, description, likes + 1, user_id, created_at, updated_at 
			from 
				posts 
			where 
				id = $1`, l.PostId).Scan(&res.Id, &res.Title, &res.Description, &res.Likes, &res.UserId, &res.CreatedAt, &res.UpdatedAt)

		if err != nil {
			log.Println("failed to like post")
			return &p.PostResponse{}, err
		}
	}

	return &res, nil
}

func (r *PostRepo) UpdatePost(ctx context.Context, post *p.UpdatePostRequest) error {
	trace, ctx := opentracing.StartSpanFromContext(ctx, "UpdateAddress")
	defer trace.Finish()
	res, err := r.db.Exec(`
		update
			posts 
		set 
			title = $1, description = $2, updated_at = $3 
		where 
			id = $4`, post.Title, post.Description, time.Now(), post.Id)
	if err != nil {
		log.Println("failed to update post")
		return err
	}

	fmt.Println(res.RowsAffected())

	return nil
}

func (r *PostRepo) DeletePost(ctx context.Context, id *p.IdRequest) (*p.PostResponse, error) {
	trace, ctx := opentracing.StartSpanFromContext(ctx, "UpdateAddress")
	defer trace.Finish()
	post := p.PostResponse{}
	err := r.db.QueryRow(`
		update 
			posts 
		set 
			deleted_at = $1 
		where 
			id = $2 
		returning 
			id, title, description, likes, user_id, created_at, updated_at`, time.Now(), id.Id).Scan(&post.Id, &post.Title, &post.Description, &post.Likes, &post.UserId, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		log.Println("failed to delete post")
		return &p.PostResponse{}, err
	}

	return &post, nil
}
