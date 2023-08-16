package appeal

/// Структура для создания обращений \\\

type Appeal struct {
	ID          int64   `json:"id" example:"1567"`
	Email       string  `json:"email" example:"petrovmaksim1992@mail.ru"`
	PhoneNumber string  `json:"phone_number" example:"89656879175"`
	Nickname    string  `json:"nickname" example:"Petrov Maksim"`
	Subject     *string `json:"subject" example:"Service"`
	Message     string  `json:"message" example:"-"`
	Document    *string `json:"document" example:"-"`
}

type CreateAppealDTO struct {
	Email       string  `json:"email" example:"petrovmaksim1992@mail.ru"`
	PhoneNumber string  `json:"phone_number" example:"89656879175"`
	Nickname    string  `json:"nickname" example:"Petrov Maksim"`
	Subject     *string `json:"subject" example:"Service"`
	Message     string  `json:"message" example:"-"`
	Document    *string `json:"document" example:"-"`
}
