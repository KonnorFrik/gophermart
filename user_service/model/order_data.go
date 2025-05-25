package model

type Order struct {
    ID         int64 `json:"-"`
    Accrual    int64 `json:"accrual"`
	Number     string `json:"number"`
	Status     string `json:"status"`
    UploadedAt string `json:"uploaded_at"`
	UserID     int64 `json:"-"`
}

func validByLUHN(numbers string) bool {
    return len(numbers) > 0
}

