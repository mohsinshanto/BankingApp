package user

import "gorm.io/gorm"

type Module struct {
	Controller *UserController
}

func NewModule(db *gorm.DB) *Module {
	repo := NewUserRepository(db)
	service := NewUserService(repo)
	controller := NewUserController(service)

	return &Module{
		Controller: controller,
	}
}
