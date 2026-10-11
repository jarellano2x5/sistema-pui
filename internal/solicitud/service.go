package solicitud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jarellano2x5/finance-go/internal/reporte"
	"github.com/jarellano2x5/finance-go/internal/usuario"
	"gorm.io/gorm"
)

type SolicitudService interface {
	SendReport(reportId uuid.UUID) error
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
	var rep Envio
	if err = srv.db.Where(`"id" = ? and "pub" = ?`, mod.Id, mod.Pub).
		First(&rep).Error; err != nil {
		return err
	}
	rep.Cierre = time.Now()
	if err = srv.db.Save(&rep).Error; err != nil {
		return err
	}
	return nil
}

func (srv *SolicitudServiceImpl) SendReport(id uuid.UUID) error {
	var u usuario.Usuario
	if err := srv.db.Find(&u, id).Error; err != nil {
		return err
	}
	pay := reporte.SignX{Usu: "PUI", Cla: u.Expwd}
	jdta, err := json.Marshal(pay)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", u.Url+"/login", bytes.NewBuffer(jdta))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	cli := &http.Client{
		Timeout: 20 * time.Second,
	}
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyb, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("url return error code status %d: %s", resp.StatusCode, string(bodyb))
	}
	var auth reporte.TknX
	err = json.NewDecoder(resp.Body).Decode(&auth)
	if err != nil {
		return err
	}
	var mods []reporte.Reporte
	//srv.db.Where("usuario_id = ? and cierre is null", id).
	if err := srv.db.Where("usuario_id = ?", id).
		Find(&mods).Error; err != nil {
		return err
	}
	for _, r := range mods {
		end := ToEnvio(r)
		if ers := srv.db.Create(end).Error; ers != nil {
			return ers
		}
		env := reporte.ToSolicitudMinDto(r, end.Pub)
		fmt.Printf("%+v\n", env)
		ers := SendRpte(env, u.Url, auth.Tkn)
		if ers != nil {
			return ers
		}
	}
	return nil
}

func SendRpte(r reporte.SolicitudMinDto, url string, tkn string) error {
	jdta, err := json.Marshal(r)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url+"/activar-reporte", bytes.NewBuffer(jdta))
	if err != nil {
		return err
	}
	fmt.Println("url return error code status" + r.Id + ":" + r.Curp + ":" + *r.FechaDesaparicion + ":" + r.LugarNacimiento)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+tkn)
	cli := &http.Client{
		Timeout: 50 * time.Second,
	}
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyb, _ := io.ReadAll(resp.Body)
		fmt.Printf("url return error code status %d: %s\n", resp.StatusCode, string(bodyb))
		return fmt.Errorf("url return error code status %d: %s", resp.StatusCode, string(bodyb))
	}
	return nil
}
