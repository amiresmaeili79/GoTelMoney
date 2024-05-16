package repositories

type Registry struct {
	userRepository        UserRepository
	expenseTypeRepository ExpenseTypeRepository
	expenseRepository     ExpenseRepository
}

func NewRegistry(userRep UserRepository, eTypeRep ExpenseTypeRepository, expenseRep ExpenseRepository) *Registry {
	return &Registry{
		userRepository:        userRep,
		expenseTypeRepository: eTypeRep,
		expenseRepository:     expenseRep,
	}
}

func (r *Registry) UserRepository() UserRepository {
	return r.userRepository
}
func (r *Registry) ExpenseTypeRepository() ExpenseTypeRepository {
	return r.expenseTypeRepository
}
func (r *Registry) ExpenseRepository() ExpenseRepository {
	return r.expenseRepository
}
