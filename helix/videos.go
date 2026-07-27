package helix

import "context"

type VideoPeriod string

const (
	VideoPeriodAll   VideoPeriod = "all"
	VideoPeriodDay   VideoPeriod = "day"
	VideoPeriodMonth VideoPeriod = "month"
	VideoPeriodWeek  VideoPeriod = "week"
)

type VideoSort string

const (
	VideoSortTime     VideoSort = "time"
	VideoSortTrending VideoSort = "trending"
	VideoSortViews    VideoSort = "views"
)

type VideoType string

const (
	VideoTypeAll       VideoType = "all"
	VideoTypeArchive   VideoType = "archive"
	VideoTypeHighlight VideoType = "highlight"
	VideoTypeUpload    VideoType = "upload"
)

type VideoViewable string

const VideoViewablePublic VideoViewable = "public"

type GetVideosRequest struct {
	IDs      []string     `query:"id,omitempty"`
	UserID   *string      `query:"user_id,omitempty"`
	GameID   *string      `query:"game_id,omitempty"`
	Language *string      `query:"language,omitempty"`
	Period   *VideoPeriod `query:"period,omitempty"`
	Sort     *VideoSort   `query:"sort,omitempty"`
	Type     *VideoType   `query:"type,omitempty"`
	First    *int         `query:"first,omitempty"`
	After    *string      `query:"after,omitempty"`
	Before   *string      `query:"before,omitempty"`
}

type VideoMutedSegment struct {
	Duration int `json:"duration"`
	Offset   int `json:"offset"`
}

type Video struct {
	ID            string              `json:"id"`
	StreamID      *string             `json:"stream_id"`
	UserID        string              `json:"user_id"`
	UserLogin     string              `json:"user_login"`
	UserName      string              `json:"user_name"`
	Title         string              `json:"title"`
	Description   string              `json:"description"`
	CreatedAt     Timestamp           `json:"created_at"`
	PublishedAt   Timestamp           `json:"published_at"`
	URL           string              `json:"url"`
	ThumbnailURL  string              `json:"thumbnail_url"`
	Viewable      VideoViewable       `json:"viewable"`
	ViewCount     int64               `json:"view_count"`
	Language      string              `json:"language"`
	Type          VideoType           `json:"type"`
	Duration      string              `json:"duration"`
	MutedSegments []VideoMutedSegment `json:"muted_segments"`
}

type GetVideosData []Video

type DeleteVideosRequest struct {
	IDs []string `query:"id"`
}

type DeleteVideosData []string

func (s *VideosService) GetVideos(ctx context.Context, req GetVideosRequest) (*Response[GetVideosData], error) {
	return executeTask26Endpoint[GetVideosData](s.client, ctx, "get-videos", req, nil, nil, "")
}

func (s *VideosService) GetVideosPager(req GetVideosRequest, opts ...PagerOption) (*Pager[GetVideosData], error) {
	plan, err := newPaginationPlan(req, "after")
	if err != nil {
		return nil, err
	}
	initialCursor := ""
	baseRequest := req
	if req.Before != nil {
		initialCursor = *req.Before
		baseRequest.Before = nil
	} else if req.After != nil {
		initialCursor = *req.After
		baseRequest.After = nil
	}
	return newPagerAt(func(ctx context.Context, cursor string) (*Response[GetVideosData], error) {
		if cursor == "" {
			return s.GetVideos(ctx, baseRequest)
		}
		next, err := plan.withCursor(baseRequest, cursor)
		if err != nil {
			return nil, err
		}
		request, ok := next.(GetVideosRequest)
		if !ok {
			return nil, &paginationRequestError{reason: "videos request has unexpected type"}
		}
		return s.GetVideos(ctx, request)
	}, initialCursor, opts...)
}

func (s *VideosService) DeleteVideos(ctx context.Context, req DeleteVideosRequest) (*Response[DeleteVideosData], error) {
	return executeTask26Endpoint[DeleteVideosData](s.client, ctx, "delete-videos", req, nil, nil, "")
}
