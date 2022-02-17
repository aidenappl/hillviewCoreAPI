package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type ListVideosRequest struct {
	Limit           *uint64
	IncludeArchived bool
	IncludeDrafts   bool
}

func ListVideos(db db.Queryable, req ListVideosRequest) ([]*structs.Video, error) {

	q := sq.Select(
		"videos.id",
		"videos.title",
		"videos.description",
		"videos.thumbnail",
		"videos.url",
		"videos.inserted_at",

		"video_statuses.id",
		"video_statuses.name",
		"video_statuses.short_name",
	).From("videos").
		LeftJoin("video_statuses ON videos.status = video_statuses.id").
		OrderBy("videos.id DESC").
		Limit(*req.Limit).
		Where(sq.Eq{"video_statuses.id": 1})

	if req.IncludeArchived {
		q = q.Where(sq.Or{
			sq.Eq{"video_statuses.id": 4},
			sq.Eq{"video_statuses.id": 1},
		})
	}

	if req.IncludeDrafts {
		q = q.Where(sq.Or{
			sq.Eq{"video_statuses.id": 2},
			sq.Eq{"video_statuses.id": 1},
		})
	}

	if req.IncludeDrafts && req.IncludeArchived {
		q = q.Where(sq.Or{
			sq.Eq{"video_statuses.id": 2},
			sq.Eq{"video_statuses.id": 4},
			sq.Eq{"video_statuses.id": 1},
		})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to create query: %w", err)
	}

	fmt.Println(query)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	defer rows.Close()

	var videos []*structs.Video

	for rows.Next() {
		var video structs.Video
		var status structs.GeneralNSN

		err = rows.Scan(
			&video.ID,
			&video.Title,
			&video.Description,
			&video.Thumbnail,
			&video.URL,
			&video.InsertedAt,

			&status.ID,
			&status.Name,
			&status.ShortName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		video.Status = &status

		videos = append(videos, &video)

	}

	return videos, nil
}

type EditVideoRequest struct {
	ID            *int                `json:"id"`
	Modifications *VideoModifications `json:"modifications"`
}

type VideoModifications struct {
	Title       *string `json:"title"`
	Thumbnail   *string `json:"thumbnail"`
	URL         *string `json:"url"`
	Description *string `json:"description"`
	Status      *int    `json:"status"`
}

func EditVideo(db db.Queryable, req EditVideoRequest) (*structs.Video, error) {
	if req.Modifications == nil {
		return nil, fmt.Errorf("no modifications provided")
	}

	if req.ID == nil {
		return nil, fmt.Errorf("no asset id provided")
	}

	dataToSet := map[string]interface{}{}

	if req.Modifications.Description != nil {
		dataToSet["description"] = *req.Modifications.Description
	}

	if req.Modifications.Thumbnail != nil {
		dataToSet["thumbnail"] = *req.Modifications.Thumbnail
	}

	if req.Modifications.Title != nil {
		dataToSet["title"] = *req.Modifications.Title
	}

	if req.Modifications.URL != nil {
		dataToSet["url"] = *req.Modifications.URL
	}

	if req.Modifications.Status != nil {
		dataToSet["status"] = *req.Modifications.Status
	}

	query, args, err := sq.Update("videos").
		SetMap(dataToSet).
		Where(sq.Eq{"id": *req.ID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sql query: %w", err)
	}

	video, err := GetVideo(db, GetVideoRequest{
		ID: req.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}

	return video, nil
}

type GetVideoRequest struct {
	ID *int `json:"id"`
}

func GetVideo(db db.Queryable, req GetVideoRequest) (*structs.Video, error) {
	if req.ID == nil || *req.ID == 0 {
		return nil, fmt.Errorf("no id provided")
	}

	query, args, err := sq.Select(
		"videos.id",
		"videos.title",
		"videos.description",
		"videos.thumbnail",
		"videos.url",
		"videos.inserted_at",
	).
		From("videos").
		Where(sq.Eq{"id": *req.ID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sql query: %w", err)
	}

	if !rows.Next() {
		return nil, nil
	}

	defer rows.Close()

	video := structs.Video{}

	err = rows.Scan(
		&video.ID,
		&video.Title,
		&video.Description,
		&video.Thumbnail,
		&video.URL,
		&video.InsertedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return &video, nil
}
