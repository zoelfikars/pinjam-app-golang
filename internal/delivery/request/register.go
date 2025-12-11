package request

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=4"`
	Password string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
	NamaLengkap          string `json:"nama_lengkap" validate:"required"`
}
