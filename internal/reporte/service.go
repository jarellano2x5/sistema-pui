package reporte

import (
	"log"

	"github.com/google/uuid"
	regen "github.com/zach-klippenstein/goregen"
	"gorm.io/gorm"
)

type ReporteService interface {
	GetAll(entidadId uuid.UUID) ([]ReporteDto, error)
	CreateReport(entidadId uuid.UUID, curps []string) error
	Generate(quantity int) ([]string, error)
}

type ReporteServiceImpl struct {
	db *gorm.DB
}

func NewReporteService(db *gorm.DB) ReporteService {
	return &ReporteServiceImpl{db}
}

func (srv *ReporteServiceImpl) GetAll(entidadId uuid.UUID) ([]ReporteDto, error) {
	var mods []Reporte
	if err := srv.db.Where("usuario_id = ?", entidadId).
		Find(&mods).Error; err != nil {
		return nil, err
	}
	var ld []ReporteDto
	for _, r := range mods {
		ld = append(ld, ToReporteDto(r))
	}
	return ld, nil
}

func (srv *ReporteServiceImpl) CreateReport(entidadId uuid.UUID, curps []string) error {
	var mods []Reporte
	for _, d := range curps {
		mods = append(mods, ToReporteString(d, entidadId))
	}
	if err := srv.db.Create(&mods).Error; err != nil {
		return err
	}
	return nil
}

func (srv *ReporteServiceImpl) Generate(quantity int) ([]string, error) {
	var rgs []string
	pttr := "^[A-Z][AEIOUX][A-Z]{2}\\d{2}(?:0[1-9]|1[0-2])(?:0[1-9]|[12]\\d|3[01])[HM](?:AS|B[CS]|C[CLMSH]|D[FG]|G[TR]|HG|JC|M[CNS]|N[ETL]|OC|PL|Q[TR]|S[PLR]|T[CSL]|VZ|YN|ZS)[B-DF-HJ-NP-TV-Z]{3}[A-Z\\d]\\d$"
	i := 0
	for i < quantity {
		str, err := regen.Generate(pttr)
		if err != nil {
			log.Fatalf("Error")
		} else {
			rgs = append(rgs, str)
		}
		i++
	}
	return rgs, nil
}
