package solicitud

import (
	"time"

	"github.com/jarellano2x5/finance-go/internal/reporte"
	"gorm.io/gorm"
)

type SolicitudService interface {
	Notificar(dto SolicitudDto) error
	Cierre(dto CierreDto) error
}

type SolicitudServiceImpl struct {
	db *gorm.DB
}

func NewSolicitudService(db *gorm.DB) SolicitudService {
	return &SolicitudServiceImpl{db}
}

func (srv *SolicitudServiceImpl) Notificar(dto SolicitudDto) error {
	mod, err := ToSolicitud(dto)
	if err != nil {
		return err
	}
	if err = srv.db.Create(&mod).Error; err != nil {
		return err
	}
	return nil
}

func (srv *SolicitudServiceImpl) Cierre(dto CierreDto) error {
	mod, err := ToCierre(dto)
	if err != nil {
		return err
	}
	var rep reporte.Reporte
	if err = srv.db.Where(`"id" = ? and "usuario_id" = ?`, mod.Id, mod.Pub).
		First(&rep).Error; err != nil {
		return err
	}
	rep.Cierre = time.Now()
	if err = srv.db.Save(&rep).Error; err != nil {
		return err
	}
	return nil
}
