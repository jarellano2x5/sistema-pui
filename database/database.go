package database

import (
	"fmt"
	"log"

	"github.com/jarellano2x5/finance-go/config"
	"github.com/jarellano2x5/finance-go/internal/reporte"
	"github.com/jarellano2x5/finance-go/internal/solicitud"
	"github.com/jarellano2x5/finance-go/internal/usuario"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPwrd, cfg.DBName,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("database connected successfully")

	err = DB.AutoMigrate(&usuario.Usuario{},
		&reporte.Reporte{},
		&solicitud.Envio{},
		&solicitud.Solicitud{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("database migrated sucessfully")
}
