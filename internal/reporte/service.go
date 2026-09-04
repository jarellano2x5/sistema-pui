package reporte

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jarellano2x5/finance-go/internal/usuario"
	regen "github.com/zach-klippenstein/goregen"
	"gorm.io/gorm"
)

type ReporteService interface {
	GetAll(entidadId uuid.UUID) ([]ReporteDto, error)
	CreateReport(entidadId uuid.UUID, curps []string) error
	Generate(quantity int) ([]string, error)
	SendReport(reportId uuid.UUID) error
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

func (srv *ReporteServiceImpl) SendReport(id uuid.UUID) error {
	var u usuario.Usuario
	if err := srv.db.Find(&u, id).Error; err != nil {
		return err
	}
	pay := SignX{Usu: "PUI", Cla: u.Expwd}
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
	var auth TknX
	err = json.NewDecoder(resp.Body).Decode(&auth)
	if err != nil {
		return err
	}
	var mods []Reporte
	//srv.db.Where("usuario_id = ? and cierre is null", id).
	if err := srv.db.Where("usuario_id = ?", id).
		Find(&mods).Error; err != nil {
		return err
	}
	for _, r := range mods {
		ers := SendRpte(ToSolicitudMinDto(r), u.Url, auth.Tkn)
		if ers != nil {
			return ers
		}
	}
	return nil
}

func SendRpte(r SolicitudMinDto, url string, tkn string) error {
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

func (srv *ReporteServiceImpl) Generate(quantity int) ([]string, error) {
	var rgs []string
	pttr := `^[A-Z]{1}[AEIOUX]{1}[A-Z]{2}[0-9]{2}(0[1-9]|1[0-2])(0[1-9]|[12][0-9]|3[01])[HM]{1}(AS|BC|BS|CC|CH|CH|CL|CM|CS|DF|DG|GR|GT|HG|JC|MC|MN|MS|NT|NL|OC|PL|QT|QR|SP|SL|SR|TC|TS|TL|VZ|YN|ZS|NE)[B-DF-HJ-NP-TV-Z]{3}[0-9A-Z]{1}[0-9]{1}$`
	i := 0
	for i < quantity {
		str, err := regen.Generate(pttr)
		if err != nil {
			return nil, err
		} else {
			rgs = append(rgs, str)
		}
		i++
	}
	return rgs, nil
}
