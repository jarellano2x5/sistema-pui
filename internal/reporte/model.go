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
	UsuarioId uuid.UUID `gorm:"type:uuid;not null"` //;REFERENCES usuarios(id) ON DELETE CASCADE ON UPDATE CASCADE"`
	Cierre    time.Time
}

type ReporteDto struct {
	Id    string `json:"required"`
	Curp  string `json:"curp" validate:"required,min=18,max=18"`
	Sexo  string `json:"sexo_asignado"`
	Nac   string `json:"fecha_nacimiento"`
	Fecha time.Time
}

type SolicitudMinDto struct {
	Id                 string  `json:"id"`
	Curp               string  `json:"curp"`
	Nombre             *string `json:"nombre"`
	PrimerApellido     *string `json:"primer_apellido"`
	SegundoApellido    *string `json:"segundo_apellido"`
	FechaNacimiento    *string `json:"fecha_nacimiento"`
	FechaDesaparicion  *string `json:"fecha_desaparicion"`
	LugarNacimiento    string  `json:"lugar_nacimiento"`
	SexoAsignado       *string `json:"sexo_asignado"`
	Telefono           *string `json:"telefono"`
	Correo             *string `json:"correo"`
	Direccion          *string `json:"direccion"`
	Calle              *string `json:"calle"`
	Numero             *string `json:"numero"`
	Colonia            *string `json:"colonia"`
	CodigoPostal       *string `json:"codigo_postal"`
	MunicipioOAlcaldia *string `json:"municipio_o_alcaldia"`
	EntidadFederativa  *string `json:"entidad_federativa"`
}

type SignX struct {
	Usu string `json:"usuario"`
	Cla string `json:"clave"`
}

type TknX struct {
	Tkn string `json:"token"`
}

func ToReporteDto(m Reporte) ReporteDto {
	var sex string = m.Curp[10:11]
	var fnc string = "20" + m.Curp[4:6] + "-" + m.Curp[6:8] + "-" + m.Curp[8:10]
	return ReporteDto{
		Id:    strings.Join([]string{m.Id.String(), m.UsuarioId.String()}, "-"),
		Curp:  m.Curp,
		Fecha: m.Fecha,
		Sexo:  sex,
		Nac:   fnc,
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

func ToSolicitudMinDto(m Reporte) SolicitudMinDto {
	var f string
	f = m.Fecha.Format("2006-01-02")
	var sex string = m.Curp[10:11]
	var fnc string = "20" + m.Curp[4:6] + "-" + m.Curp[6:8] + "-" + m.Curp[8:10]
	return SolicitudMinDto{
		Id:                 strings.Join([]string{m.Id.String(), m.UsuarioId.String()}, "-"),
		Curp:               m.Curp,
		Nombre:             nil,
		PrimerApellido:     nil,
		SegundoApellido:    nil,
		FechaNacimiento:    &fnc,
		FechaDesaparicion:  &f,
		LugarNacimiento:    GetLugar(m.Curp[11:13]),
		SexoAsignado:       &sex,
		Telefono:           nil,
		Correo:             nil,
		Direccion:          nil,
		Calle:              nil,
		Numero:             nil,
		Colonia:            nil,
		CodigoPostal:       nil,
		MunicipioOAlcaldia: nil,
		EntidadFederativa:  nil,
	}
}

func GetLugar(e string) string {
	switch e {
	case "AS":
		return "AGUASCALIENTES"
	case "BC":
		return "BAJA CALIFORNIA"
	case "BS":
		return "BAJA CALIFORNIA SUR"
	case "CC":
		return "CAMPECHE"
	case "CS":
		return "CHIAPAS"
	case "CH":
		return "CHIHUAHUA"
	case "DF":
		return "CDMX"
	case "CL":
		return "COAHUILA"
	case "CM":
		return "COLIMA"
	case "DG":
		return "DURANGO"
	case "GT":
		return "GUANAJUATO"
	case "GR":
		return "GUERRERO"
	case "HG":
		return "HIDALGO"
	case "JC":
		return "JALISCO"
	case "MC":
		return "MÉXICO"
	case "MN":
		return "MICHOACÁN"
	case "MS":
		return "MORELOS"
	case "NT":
		return "NAYARIT"
	case "NL":
		return "NUEVO LEÓN"
	case "OC":
		return "OAXACA"
	case "PL":
		return "PUEBLA"
	case "QO":
		return "QUERÉTARO"
	case "QR":
		return "QUINTANA ROO"
	case "SP":
		return "SAN LUIS POTOSÍ"
	case "SL":
		return "SINALOA"
	case "SR":
		return "SONORA"
	case "TC":
		return "TABASCO"
	case "TS":
		return "TAMAULIPAS"
	case "TL":
		return "TLAXCALA"
	case "VZ":
		return "VERACRUZ"
	case "YN":
		return "YUCATÁN"
	case "ZS":
		return "ZACATECAS"
	case "NE":
		return "FORÁNEO"
	default:
		return "DESCONOCIDO"
	}
}
