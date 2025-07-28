package dtos

type CreateUserInput struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=member manager"`
	Password string `json:"password" binding:"required,min=6"`
}
