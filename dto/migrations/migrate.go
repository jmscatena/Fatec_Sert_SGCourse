package migrations

import (
	"fmt"
	admin "github.com/jmscatena/Fatec_Sert_SGCourse/dto/models/administrativo"
	curso "github.com/jmscatena/Fatec_Sert_SGCourse/dto/models/cursos"
	"gorm.io/gorm"
)

func RunMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&admin.Usuario{},
		&curso.Curso{},
		&curso.Disciplina{},
		&curso.Documento{},
		&curso.Solicitacao_Doc{},
	)
	if err != nil {
		fmt.Println("Migrating database erro:", err)
		return
	}

}
