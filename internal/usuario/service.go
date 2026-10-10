package usuario

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jarellano2x5/finance-go/internal/utils"
	"gorm.io/gorm"
)

type UsuarioService interface {
	Login(dto LoginDto) (SignDto, error)
	Register(dto UsuarioDto) (SignDto, error)
	Cancel(id uuid.UUID) error
}

type UsuarioServiceImpl struct {
	db *gorm.DB
}

func NewUsuarioService(db *gorm.DB) UsuarioService {
	return &UsuarioServiceImpl{db}
}

func (s *UsuarioServiceImpl) Login(dto LoginDto) (SignDto, error) {
	var u Usuario
	err := s.db.Where(`"user" = ?`, dto.User).First(&u).Error
	if err != nil {
		return SignDto{}, err
	}
	fmt.Println("###### UUID user:", u.Id.String())
	ck := utils.HashChk(dto.Pwd, u.Pwd)
	if !ck {
		return SignDto{}, errors.New("forbidden")
	}
	reg := SignDto{Id: u.Id, User: u.User}
	fmt.Println("###### UUID user:", reg.Id.String())
	return reg, nil
}

func (s *UsuarioServiceImpl) Register(dto UsuarioDto) (SignDto, error) {
	var ck bool
	if err := s.db.Model(&Usuario{}).Select("1").
		Where(`"user" = ?`, dto.User).Limit(1).
		Find(&ck).Error; err != nil {
		return SignDto{}, err
	}
	if ck {
		return SignDto{}, fmt.Errorf("Usuario alredy exists")
	}
	hsh, err := utils.HashPwd(dto.Pwd)
	if err != nil {
		return SignDto{}, err
	}
	mod := Usuario{
		Id: uuid.New(), User: dto.User,
		Pwd: hsh, Url: dto.Url,
		Expwd: dto.Expwd, Activo: true,
		CreadoEn: time.Now(),
	}
	if err := s.db.Create(&mod).Error; err != nil {
		return SignDto{}, err
	}
	reg := SignDto{
		Id: mod.Id, User: mod.User,
	}
	return reg, nil
}

func (s *UsuarioServiceImpl) Cancel(id uuid.UUID) error {
	var mod Usuario
	err := s.db.Where("id = ?", id).First(&mod)
	if err != nil {
		return err.Error
	}
	mod.Activo = false
	return s.db.Save(mod).Error
}
