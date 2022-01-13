package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type ListVideosRequest struct {
	Limit *uint64
}

func ListVideos(db db.Queryable, req ListVideosRequest) ([]*structs.Video, error) {

	query, args, err := sq.Select(
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
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to create query: %w", err)
	}

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
