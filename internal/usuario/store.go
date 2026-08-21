package usuario

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsuarioStore interface {
	Find() ([]Usuario, error)
	FindId(id uuid.UUID) (*Usuario, error)
	Insert(mod *Usuario) error
	Update(mod *Usuario) error
	Delete(id uuid.UUID) error
}

type UsuarioStoreImpl struct {
	db *gorm.DB
}

func NewUsuarioStore(db *gorm.DB) UsuarioStore {
	return &UsuarioStoreImpl{db}
}

func (s *UsuarioStoreImpl) Find() ([]Usuario, error) {
	var mods []Usuario
	rs := s.db.Find(&mods)
	if rs.Error != nil {
		return nil, rs.Error
	}
	return mods, nil
}

func (s *UsuarioStoreImpl) FindId(id uuid.UUID) (*Usuario, error) {
	var mod Usuario
	r := s.db.Find(&mod, id)
	if r.Error != nil {
		return nil, r.Error
	}
	return &mod, nil
}

func (s *UsuarioStoreImpl) Insert(mod *Usuario) error {
	return s.db.Create(mod).Error
}

func (s *UsuarioStoreImpl) Update(mod *Usuario) error {
	return s.db.Save(mod).Error
}

func (s *UsuarioStoreImpl) Delete(id uuid.UUID) error {
	return s.db.Delete(&Usuario{}, id).Error
}
