package user

type RegistrationInput struct {
	Name     string `json:"name" binding:"required" example:"Mohsin Shanto"`
	Email    string `json:"email" binding:"required,email" example:"mohsin@example.com"`
	Password string `json:"password" binding:"required,min=6,max=32,excludes= " example:"password123"`
}
type RegisterResponse struct {
	ID    uint   `json:"id" example:"1"`
	Name  string `json:"name" example:"Mohsin"`
	Email string `json:"email" example:"mohsin@example.com"`
}
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=32,excludes= "`
}
type LoginResponse struct {
	Token string `json:"token"`
}
