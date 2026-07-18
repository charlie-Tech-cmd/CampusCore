package api

import "campuscore/internal/models"

type mockUserRepository struct {
	createFunc          func(*models.User) error
	findByIDFunc        func(string) (*models.User, error)
	findByEmailFunc     func(string) (*models.User, error)
	updateLastLoginFunc func(string) error

	getProfileFunc    func(string) (*models.User, error)
	updateProfileFunc func(*models.User) error
}

func (m *mockUserRepository) Create(user *models.User) error {
	if m.createFunc != nil {
		return m.createFunc(user)
	}
	return nil
}

func (m *mockUserRepository) FindByID(id string) (*models.User, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
	return nil, nil
}

func (m *mockUserRepository) FindByEmail(email string) (*models.User, error) {
	if m.findByEmailFunc != nil {
		return m.findByEmailFunc(email)
	}
	return nil, nil
}

func (m *mockUserRepository) UpdateLastLogin(id string) error {
	if m.updateLastLoginFunc != nil {
		return m.updateLastLoginFunc(id)
	}
	return nil
}

func (m *mockUserRepository) GetProfile(id string) (*models.User, error) {
	if m.getProfileFunc != nil {
		return m.getProfileFunc(id)
	}
	return nil, nil
}

func (m *mockUserRepository) UpdateProfile(user *models.User) error {
	if m.updateProfileFunc != nil {
		return m.updateProfileFunc(user)
	}
	return nil
}
