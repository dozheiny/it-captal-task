package serializers

type State struct {
	Message string      `json:"message"`
	Code    int64       `json:"code"`
	Status  bool        `json:"status"`
	Details *string     `json:"details"`
	Counts  *Counts     `json:"counts"`
	Data    interface{} `json:"data"`
}

type Counts struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	TotalPage int `json:"total_pages"`
	Total     int `json:"total"`
}

func NewState() State {
	return State{}
}

func (s State) SetMessage(message string) State {
	s.Message = message
	return s
}

func (s State) SetCode(code int64) State {
	s.Code = code
	return s
}

func (s State) SetStatus(status bool) State {
	s.Status = status
	return s
}

func (s State) SetDetails(details string) State {
	s.Details = &details
	return s
}

func (s State) SetData(data interface{}) State {
	s.Data = data
	return s
}

func (s State) SetCounters(counts Counts) State {
	s.Counts = &counts
	return s
}
