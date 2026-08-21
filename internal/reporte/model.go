package reporte

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Reporte struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Curp      string    `gorm:"size:18;not null"`
	Fecha     time.Time
	UsuarioId uuid.UUID `gorm:"type:uuid;not null REFERENCES Usuarios(id) ON DELETE CASCADE ON UPDATE CASCADE"`
}

type ReporteDto struct {
	Id    string `json:"required"`
	Curp  string `json:"curp" validate:"required,min=18,max=18"`
	Fecha time.Time
}

func ToReporteDto(m Reporte) ReporteDto {
	return ReporteDto{
		Id:    strings.Join([]string{m.Id.String(), m.UsuarioId.String()}, "-"),
		Curp:  m.Curp,
		Fecha: m.Fecha,
	}
}

func ToReporteString(curp string, entidadId uuid.UUID) Reporte {
	return Reporte{
		Id:        uuid.New(),
		Curp:      curp,
		Fecha:     time.Now(),
		UsuarioId: entidadId,
	}
}
