package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type GetSpotlightRequest struct {
	Rank *int
}

func GetSpotlight(db db.Queryable, req GetSpotlightRequest) (*structs.Spotlight, error) {
	//  check required fields
	if req.Rank == nil {
		return nil, fmt.Errorf("required field rank is nil")
	}

	// build query

	q := sq.Select(
		"spotlight.rank",
		"spotlight.video_id",
		"spotlight.inserted_at",
		"spotlight.updated_at",

		"videos.id",
		"videos.uuid",
		"videos.title",
		"videos.description",
		"videos.thumbnail",
		"videos.url",
		"videos.download_url",
		"videos.allow_downloads",
		"videos.inserted_at",
		`(
			SELECT COUNT(video_views.id) FROM video_views WHERE video_views.video_id = videos.id
		) as views`,

		"video_statuses.id",
		"video_statuses.name",
		"video_statuses.short_name",
	)

	q = q.From("spotlight").
		Join("videos ON spotlight.video_id = videos.id").
		Join("video_statuses ON videos.status = video_statuses.id").
		Where(sq.Eq{"spotlight.rank": *req.Rank})

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	row := db.QueryRow(query, args...)
	s := &structs.Spotlight{}
	v := &structs.Video{}
	v.Status = &structs.GeneralNSN{}

	err = row.Scan(
		&s.Rank,
		&s.VideoID,
		&s.InsertedAt,
		&s.UpdatedAt,

		&v.ID,
		&v.UUID,
		&v.Title,
		&v.Description,
		&v.Thumbnail,
		&v.URL,
		&v.DownloadURL,
		&v.AllowDownloads,
		&v.InsertedAt,
		&v.Views,

		&v.Status.ID,
		&v.Status.Name,
		&v.Status.ShortName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %v", err)
	}

	s.Video = v

	return s, nil

}
