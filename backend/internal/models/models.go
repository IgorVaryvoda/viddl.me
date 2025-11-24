package models

type VideoInfo struct {
	Title        string       `json:"title"`
	Thumbnail    string       `json:"thumbnail"`
	Duration     float64      `json:"duration"`
	Uploader     string       `json:"uploader"`
	Formats      []FormatInfo `json:"formats"`
	MultiVideos  []VideoEntry `json:"multi_videos,omitempty"`
	IsMultiVideo bool         `json:"is_multi_video"`
}

type VideoEntry struct {
	Index     int     `json:"index"`
	Title     string  `json:"title"`
	Thumbnail string  `json:"thumbnail"`
	Duration  float64 `json:"duration"`
}

type FormatInfo struct {
	FormatID string `json:"format_id"`
	Ext      string `json:"ext"`
	Quality  string `json:"quality"`
	Filesize int64  `json:"filesize"`
}

type DownloadRequest struct {
	URL        string `json:"url" binding:"required"`
	Format     string `json:"format"`
	VideoIndex int    `json:"video_index"`
}

type YtDlpFormat struct {
	FormatID   string  `json:"format_id"`
	Ext        string  `json:"ext"`
	FormatNote string  `json:"format_note"`
	Quality    float64 `json:"quality"`
	Filesize   int64   `json:"filesize"`
	VCodec     string  `json:"vcodec"`
	ACodec     string  `json:"acodec"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	Fps        float64 `json:"fps"`
	Resolution string  `json:"resolution"`
}

type YtDlpInfo struct {
	Title     string        `json:"title"`
	Thumbnail string        `json:"thumbnail"`
	Duration  float64       `json:"duration"`
	Uploader  string        `json:"uploader"`
	Formats   []YtDlpFormat `json:"formats"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version,omitempty"`
	Error   string `json:"error,omitempty"`
}
