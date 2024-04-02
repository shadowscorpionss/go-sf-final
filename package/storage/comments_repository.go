package storage

import (
	"ApiGate/package/models"
	"context"
	"fmt"
	"time"
)

func (db *DB) NewComment(nc models.NewComment) (*models.Comment, error) {
	if nc.PostId < 0 {
		return nil, fmt.Errorf("wrong postId")
	}
	if len(nc.AuthorName) == 0 {
		return nil, fmt.Errorf("authorName is empty")
	}

	res := models.Comment{
		Id:         0,
		PostId:     nc.PostId,
		ParentId:   nc.ParentId,
		Content:    nc.Content,
		AuthorName: nc.AuthorName,
		PubTime:    time.Now().Unix(),
	}

	rw := db.pool.QueryRow(context.Background(), `
		INSERT INTO news_comments(pub_time, content, author_name, post_id, parent_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		res.PubTime,
		res.Content,
		res.AuthorName,
		res.PostId,
		NullInt32(int32(res.ParentId)),
	)

	var id int
	err := rw.Scan(&id)
	if err != nil {
		return nil, err
	}
	res.Id = id
	return &res, nil
}

func (db *DB) Comments(postId int) ([]models.Comment, error) {

	rws, err := db.pool.Query(context.Background(), `
		SELECT id, post_id, content, pub_time, author_name, coalesce(parent_id, -1) FROM news_comments
		WHERE post_id = $1
		ORDER BY pub_time DESC
		`,
		postId,
	)
	if err != nil {
		return nil, err
	}
	res := []models.Comment{}
	for rws.Next() {
		var c models.Comment
		err = rws.Scan(
			&c.Id,
			&c.PostId,
			&c.Content,
			&c.PubTime,
			&c.AuthorName,
			&c.ParentId,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}
