package solicitud

import (
	"time"

	"github.com/google/uuid"
)

type Solicitud struct {
	Id                     uuid.UUID `gorm:"type:uuid;primaryKey"`
	Pub                    uuid.UUID `gorm:"type:uuid;primaryKey"` // REFERENCES Reportes(id) ON DELETE CASCADE ON UPDATE CASCADE
	Curp                   string    `gorm:"size:18;not null"`
	Nombre                 string    `gorm:"size:50"`
	PrimerApellido         string    `gorm:"size:50"`
	SegundoApellido        string    `gorm:"size:50"`
	FechaNacimiento        time.Time
	LugarNacimiento        string   `gorm:"size:20"`
	SexoAsignado           string   `gorm:"size:1"`
	Telefono               string   `gorm:"size:15"`
	Correo                 string   `gorm:"size:50"`
	Direccion              string   `gorm:"size:500"`
	Calle                  string   `gorm:"size:50"`
	Numero                 string   `gorm:"size:20"`
	Colonia                string   `gorm:"size:50"`
	CodigoPostal           string   `gorm:"size:5"`
	MunicipioOAlcaldia     string   `gorm:"size:100"`
	EntidadFederativa      string   `gorm:"size:40"`
	Fotos                  []string `gorm:"type:varchar[]"`
	FormatoFotos           string   `gorm:"size:20"`
	Huellas                []string `gorm:"seralizer:json"`
	FormatoHuellas         string   `gorm:"size:50"`
	InstitucionId          string   `gorm:"size:13;not null"`
	TipoEvento             string   `gorm:"size:500"`
	FechaEvento            time.Time
	DescripcionLugarEvento string `gorm:"size:500"`
	DireccionEvento        string `gorm:"size:500"`
	FaseBusqueda           string `gorm:"size:1"`
	Fecha                  time.Time
}
