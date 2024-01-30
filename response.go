package go_cake

type ResponseJSON struct {
	Items []map[string]any `json:"_items"`
	Meta  MetaJSON         `json:"_meta"`
}

type MetaJSON struct {
	StatusCode      int     `json:"status_code"`
	StatusMessage   string  `json:"status_message"`
	Total           uint64  `json:"total"`
	TotalTimeMs     float64 `json:"total_time_ms"`
	Page            int64   `json:"page"`
	PerPage         int64   `json:"per_page"`
	RequestUniqueID string  `json:"request_unique_id"`
	Version         string  `json:"version"`
	Method          string  `json:"method"`
	URL             string  `json:"url"`
}

type ItemStatusJSON struct {
	Meta ItemStatusMetaJSON `json:"_meta"`
}

type ItemStatusMetaJSON struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

func (rj *ResponseJSON) ResetItems() {
	rj.Items = make([]map[string]any, 0)
}

func NewResponseJSON() *ResponseJSON {
	response := ResponseJSON{}
	response.ResetItems()

	return &response
}
