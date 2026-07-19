package types

type GroupStats struct {
	GroupName  string
	GroupID    int64
	Done       int64
	Todo       int64
	InProgress int64
	Total      int64
}
