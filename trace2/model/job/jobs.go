package job

type GetJobRequest struct {
	ID string `json:"id"`
}

type GetJobResponse struct {
	ID         uint     `json:"id"`
	Name       string   `json:"name"`
	StrategyID uint     `json:"strategy_id"`
	TimeOut    int64    `json:"timeOut"`
	Content    string   `json:"content"`
	FilePath   string   `json:"filePath"`
	Target     []string `json:"target"`
}
