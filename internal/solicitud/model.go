package solicitud

import (
	"log"
	"time"

	"github.com/google/uuid"
)

type Huellas struct {
	Rone  string `json:"rone"`
	Rtwo  string `json:"rtwo"`
	Lpalm string `json:"lpalm"`
}

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
	Huellas                Huellas  `gorm:"type:json;seralizer:json"`
	FormatoHuellas         string   `gorm:"size:50"`
	InstitucionId          string   `gorm:"size:13;not null"`
	TipoEvento             string   `gorm:"size:500"`
	FechaEvento            time.Time
	DescripcionLugarEvento string `gorm:"size:500"`
	DireccionEvento        string `gorm:"size:500"`
	FaseBusqueda           string `gorm:"size:1"`
	Fecha                  time.Time
}

type Cierre struct {
	Id            uuid.UUID
	Pub           uuid.UUID
	InstitucionId string
}

type SolicitudDto struct {
	Id                     string   `json:"id" validate:"required,min=36,max=75"`
	Curp                   string   `json:"curp" validate:"required,len=18"`
	NombreCompleto         string   `json:"nombre_completo"`
	Nombre                 string   `json:"nombre" validate:"max=50"`
	PrimerApellido         string   `json:"primer_apellido" validate:"max=50"`
	SegundoApellido        string   `json:"segundo_apellido" validate:"max=50"`
	FechaNacimiento        string   `json:"fecha_nacimiento" validate:"len=10"`
	LugarNacimiento        string   `json:"lugar_nacimiento" validate:"required,max=20"`
	SexoAsignado           string   `json:"sexo_asignado" validate:"len=1"`
	Telefono               string   `json:"telefono" validate:"max=15"`
	Correo                 string   `json:"correo" validate:"max=50"`
	Domicilio              string   `json:"domicilio"`
	Direccion              string   `json:"direccion" validate:"max=500"`
	Calle                  string   `json:"calle" validate:"max=50"`
	Numero                 string   `json:"numero" validate:"max=20"`
	Colonia                string   `json:"colonia" validate:"max=50"`
	CodigoPostal           string   `json:"codigo_postal" validate:"max=5"`
	MunicipioOAlcaldia     string   `json:"municipio_o_alcaldia" validate:"max=100"`
	EntidadFederativa      string   `json:"entidad_federativa" validate:"max=40"`
	Fotos                  []string `json:"fotos"`
	FormatoFotos           string   `json:"formato_fotos" validate:"max=20"`
	Huellas                Huellas  `json:"huellas"`
	FormatoHuellas         string   `json:"formato_huellas" validate:"max=50"`
	InstitucionId          string   `json:"institucion_id" validate:"min=4,max=13"`
	TipoEvento             string   `json:"tipo_evento" validate:"max=500"`
	FechaEvento            string   `json:"fecha_evento" validate:"len=10"`
	DescripcionLugarEvento string   `json:"descripcion_lugar_evento" validate:"max=500"`
	DireccionEvento        string   `json:"direccion_evento" validate:"max=500"`
	FaseBusqueda           string   `json:"fase_busqueda" validate:"required,len=1"`
}

type CierreDto struct {
	Id            string `json:"id" validate:"required,min=36,max=75"`
	InstitucionId string `json:"institucion_id" validate:"required,min=4,max=13"`
}

func ToSolicitud(dto SolicitudDto) (Solicitud, error) {
	var ly string = "2006-01-02"
	id, err := uuid.Parse(dto.Id[:36])
	if err != nil {
		log.Fatalf("failed to parse UUID: %v", err)
	}
	pub, err := uuid.Parse(dto.Id[37:])
	if err != nil {
		log.Fatalf("failed to parse pub UUID: %v", err)
	}
	fn, err := time.Parse(ly, dto.FechaNacimiento)
	if err != nil {
		log.Fatalf("failed to parse fnac UUID: %v", err)
	}
	fe, err := time.Parse(ly, dto.FechaEvento)
	if err != nil {
		log.Fatalf("failed to parse feve UUID: %v", err)
	}
	return Solicitud{
		Id:                     id,
		Pub:                    pub,
		Curp:                   dto.Curp,
		Nombre:                 dto.Nombre,
		PrimerApellido:         dto.PrimerApellido,
		SegundoApellido:        dto.SegundoApellido,
		FechaNacimiento:        fn,
		LugarNacimiento:        dto.LugarNacimiento,
		SexoAsignado:           dto.SegundoApellido,
		Telefono:               dto.Telefono,
		Correo:                 dto.Correo,
		Direccion:              dto.Direccion,
		Calle:                  dto.Calle,
		Numero:                 dto.Numero,
		Colonia:                dto.Colonia,
		CodigoPostal:           dto.CodigoPostal,
		MunicipioOAlcaldia:     dto.MunicipioOAlcaldia,
		EntidadFederativa:      dto.EntidadFederativa,
		Fotos:                  dto.Fotos,
		FormatoFotos:           dto.FormatoFotos,
		Huellas:                dto.Huellas,
		FormatoHuellas:         dto.FormatoHuellas,
		InstitucionId:          dto.InstitucionId,
		TipoEvento:             dto.TipoEvento,
		FechaEvento:            fe,
		DescripcionLugarEvento: dto.DescripcionLugarEvento,
		DireccionEvento:        dto.DireccionEvento,
		FaseBusqueda:           dto.FaseBusqueda,
		Fecha:                  time.Now(),
	}, nil
}

func ToCierre(dto CierreDto) (Cierre, error) {
	id, err := uuid.Parse(dto.Id[:36])
	if err != nil {
		log.Fatalf("failed to parse UUID: %v", err)
	}
	pub, err := uuid.Parse(dto.Id[37:])
	if err != nil {
		log.Fatalf("failed to parse pub UUID: %v", err)
	}
	return Cierre{
		Id:            id,
		Pub:           pub,
		InstitucionId: dto.InstitucionId,
	}, nil
}
