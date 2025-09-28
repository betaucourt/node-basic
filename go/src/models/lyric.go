package models

// Lyric represents a lyric line
type Lyric struct {
	Song int    `json:"song"`
	Line int    `json:"line"`
	Text string `json:"text"`
}
