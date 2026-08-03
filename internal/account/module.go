package account

import "gorm.io/gorm"

type Module struct {
	Controller *AccountController
}

func NewModule(db *gorm.DB) *Module {
	repo := NewAccountRepository(db)
	service := NewAccountService(repo)
	controller := NewAccountController(service)

	return &Module{
		Controller: controller,
	}
}
