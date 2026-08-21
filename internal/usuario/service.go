package usuario

import (
	"time"

	"github.com/google/uuid"
)

type UsuarioService interface {
	GetAll() ([]SignDto, error)
	CreateUser(dto UsuarioDto) (SignDto, error)
	RepairUser(dto UsuarioDto) error
}

type UsuarioServiceImpl struct {
	sto UsuarioStore
}

func NewUsuarioService(sto UsuarioStore) UsuarioService {
	return &UsuarioServiceImpl{sto: sto}
}

func (s *UsuarioServiceImpl) GetAll() ([]SignDto, error) {
	lm, err := s.sto.Find()
	if err != nil {
		return nil, err
	}
	var ld []SignDto
	for _, r := range lm {
		ld = append(ld, ToSignDto(r))
	}
	return ld, nil
}

func (s *UsuarioServiceImpl) CreateUser(dto UsuarioDto) (SignDto, error) {
	mod := &Usuario{
		Id: uuid.New(), User: dto.User,
		Pwd: dto.Pwd, Url: dto.Url,
		Expwd: dto.Expwd, Activo: true,
		CreadoEn: time.Now(),
	}
	if err := s.sto.Insert(mod); err != nil {
		return SignDto{}, err
	}
	return ToSignDto(*mod), nil
}

func (s *UsuarioServiceImpl) RepairUser(dto UsuarioDto) error {
	mod := &Usuario{
		Id: dto.Id, User: dto.User,
		Pwd: dto.Pwd, Url: dto.Url,
		Expwd: dto.Expwd, Activo: true,
	}
	return s.sto.Update(mod)
}
