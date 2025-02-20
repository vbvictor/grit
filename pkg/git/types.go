package git

type ChurnChunk struct {
	File    string `json:"path"`
	Churn   int    `json:"changes"`
	Added   int    `json:"additions"`
	Removed int    `json:"deletions"`
	Commits int    `json:"commits"`
}
