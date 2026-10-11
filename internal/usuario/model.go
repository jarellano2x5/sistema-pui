package usuario

import (
	"time"

	"github.com/google/uuid"
)

type Usuario struct {
	Id       uuid.UUID `gorm:"type:uuid;primaryKey"`
	User     string    `gorm:"size:18;not null"`
	Pwd      string    `gorm:"size:500;not null"`
	Url      string    `gorm:"size:200;not null"`
	Expwd    string    `gorm:"size:20;not null"`
	Activo   bool
	CreadoEn time.Time
}

type LoginDto struct {
	User string `json:"institucion_id" validate:"required,min=4,max=13"`
	Pwd  string `json:"clave" validate:"required,min=8,max=18"`
}

type UsuarioDto struct {
	User  string `json:"institucion_id" validate:"required,min=4,max=13"`
	Pwd   string `json:"clave" validate:"required,min=8,max=18"`
	Url   string `json:"url" validate:"required,min=20,max=200"`
	Expwd string `json:"clave_pui" validate:"required,min=8,max=20"`
}

type SignDto struct {
	Id   uuid.UUID
	User string
}

type TokenDto struct {
	Tkn string `json:"token"`
}

func ToSignDto(m Usuario) SignDto {
	return SignDto{
		Id: m.Id, User: m.User,
	}
}

func (r *UsuarioDto) Validate() map[string]string {
	errors := make(map[string]string)
	if r.User == "" {
		errors["usuario"] = "Usuario is required"
	}
	if len(r.User) > 18 {
		errors["usuario"] = "Usuario must be 15 chars or less"
	}
	return errors
}
