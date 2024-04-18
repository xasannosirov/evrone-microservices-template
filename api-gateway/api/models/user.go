package models

type (
	User struct {
		Id        string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Bio       string `json:"bio"`
		Website   string `json:"website"`
	}

	UserList struct {
		Users []*User `json:"users"`
	}

	GetByID struct {
		ID string `json:"id"`
	}

	DelResp struct {
		Status bool `json:"status"`
	}

	GetList struct {
		Page  uint64 `json:"page"`
		Limit uint64 `json:"limit"`
	}
)
