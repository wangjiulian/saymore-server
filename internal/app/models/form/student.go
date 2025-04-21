package form

type (
	CodeForm struct {
		Code string `json:"code"`
	}

	LoginForm struct {
		EncryptedData string `json:"encrypted_data"`
		Iv            string `json:"iv"`
		SessionKey    string `json:"session_key"`
	}

	StudentEditForm struct {
		Avatar          string `json:"avatar"`
		Nickname        string `json:"nickname"`
		Gender          int    `json:"gender"`
		BirthDate       string `json:"birth_date"`
		StudentType     int    `json:"student_type"`
		LearningPurpose int    `json:"learning_purpose"`
		EnglishLevel    int    `json:"english_level"`
	}
)
