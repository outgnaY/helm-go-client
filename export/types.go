package export

import "time"

// ChartInfo to avoid import circles
type ChartInfo struct {
	Ref       string
	Name      string
	Version   string
	Digest    string
	Size      int64
	CreatedAt time.Time
}
