package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type ListVideosRequest struct {
	// Limit, sort and offset are required
	Limit  *int
	Offset *int
	Sort   *string

	// Flags
	Statuses *[]int

	// Search
	Search     *string
	PlaylistID *int
	ID         *int
}

func ListVideos(db db.Queryable, req ListVideosRequest) ([]*structs.Video, error) {
	// check required fields
	if req.Limit == nil {
		return nil, fmt.Errorf("no limit provided")
	}

	if req.Offset == nil {
		return nil, fmt.Errorf("no offset provided")
	}

	// check sort formatting
	if req.Sort != nil {
		if *req.Sort != "desc" && *req.Sort != "asc" {
			return nil, fmt.Errorf("invalid sort provided")
		}
	} else {
		req.Sort = new(string)
		*req.Sort = "desc"
	}

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
		OrderBy("videos.id " + *req.Sort).
		Limit(uint64(*req.Limit)).
		Offset(uint64(*req.Offset))

	if req.PlaylistID != nil {
		q = q.LeftJoin("playlist_associations ON videos.id = playlist_associations.video_id")
		q = q.Where(sq.Eq{"playlist_associations.playlist_id": *req.PlaylistID})
	}

	if req.ID != nil {
		q = q.Where(sq.Eq{"videos.id": *req.ID})
	}

	if req.Statuses != nil {
		q = q.Where(sq.Eq{"videos.status": *req.Statuses})
	}

	if req.Search != nil {
		q = q.Where(sq.Or{
			sq.Like{"videos.title": "%" + *req.Search + "%"},
			sq.Like{"videos.description": "%" + *req.Search + "%"},
		})
	}

	query, args, err := q.ToSql()
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
