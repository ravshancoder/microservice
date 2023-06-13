package postgres

import (
	"context"
	"log"
	"time"

	c "github.com/microservice/comment_service/genproto/comment"
	"github.com/opentracing/opentracing-go"
)

func (r *CommentRepo) WriteComment(ctx context.Context, comment *c.CommentRequest) (*c.CommentResponse, error) {
	trace, ctx := opentracing.StartSpanFromContext(ctx, "UpdateAddress")
	defer trace.Finish()
	var res c.CommentResponse
	err := r.db.QueryRow(`
		insert into 
			comments(post_id, user_id, text)
		values
			($1, $2, $3) 
		returning 
			id, post_id, user_id, text, created_at`, comment.PostId, comment.UserId, comment.Text).Scan(&res.Id, &res.PostId, &res.UserId, &res.Text, &res.CreatedAt)

	if err != nil {
		log.Println("failed to create comment: ", err)
		return &c.CommentResponse{}, err
	}

	return &res, nil
}

func (r *CommentRepo) GetComments(ctx context.Context, com *c.GetAllCommentsRequest) (*c.Comments, error) {
	trace, ctx := opentracing.StartSpanFromContext(ctx, "UpdateAddress")
	defer trace.Finish()
	var res c.Comments
	rows, err := r.db.Query(`
		select 
			id, post_id, user_id, text, created_at
		from 
			comments 
		where 
			post_id = $1 and deleted_at is null`, com.PostId)

	if err != nil {
		log.Println("failed to get comment: ", err)
		return &c.Comments{}, nil
	}

	for rows.Next() {
		comment := c.CommentResponse{}

		err = rows.Scan(
			&comment.Id,
			&comment.PostId,
			&comment.UserId,
			&comment.Text,
			&comment.CreatedAt,
		)

		if err != nil {
			log.Println("failed to scanning comment: ", err)
			return &c.Comments{}, err
		}

		res.Comments = append(res.Comments, &comment)
	}

	return &res, nil
}

func (r *CommentRepo) GetCommentsForPost(ctx context.Context, com *c.GetAllCommentsRequest) (*c.Comments, error) {
	trace, ctx := opentracing.StartSpanFromContext(ctx, "UpdateAddress")
	defer trace.Finish()
	var res c.Comments
	rows, err := r.db.Query(`
		select 
			id, post_id, user_id, text, created_at 
		from 
			comments 
		where 
			post_id = $1 and deleted_at is null`, com.PostId)

	if err != nil {
		log.Println("failed to get comment for post: ", err)
		return &c.Comments{}, nil
	}

	for rows.Next() {
		comment := c.CommentResponse{}

		err = rows.Scan(
			&comment.Id,
			&comment.PostId,
			&comment.UserId,
			&comment.Text,
			&comment.CreatedAt,
		)

		if err != nil {
			log.Println("failed to scanning comment: ", err)
			return &c.Comments{}, err
		}

		res.Comments = append(res.Comments, &comment)
	}

	return &res, nil
}

func (r *CommentRepo) DeleteComment(ctx context.Context, id *c.IdRequest) (*c.CommentResponse, error) {
	trace, ctx := opentracing.StartSpanFromContext(ctx, "UpdateAddress")
	defer trace.Finish()
	var res c.CommentResponse
	err := r.db.QueryRow(`
		update 
			comments 
		set 
			deleted_at = $1 
		where 
			id = $2 
		returning 
			id, post_id, user_id, text, created_at`, time.Now(), id.Id).Scan(&res.Id, &res.PostId, &res.UserId, &res.Text, &res.CreatedAt)

	if err != nil {
		log.Println("failed to delete comment", err)
	}

	return &res, nil
}
