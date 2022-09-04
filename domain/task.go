package domain

type (
	Task struct {
		ID   string
		Name string
	}

	TaskModel struct {
		ID   string
		Name string
	}

	TaskFetchRequest struct {
		IDs []string `json:"ids"`
	}

	TaskFetchResponse struct {
		Names []string `json:"names"`
	}
)
