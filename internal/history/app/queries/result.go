package queries

type GetHistoryResult struct {
	Items []*HistoryMessage
}

type HistoryMessage map[string]interface{}
