package handlers

import (
	"github.com/jmscatena/Fatec_Sert_SGCourse/dto/models/administrativo"
	"github.com/jmscatena/Fatec_Sert_SGCourse/dto/models/cursos"
	"gorm.io/gorm"
)

type Tables interface {
	administrativo.Usuario | cursos.Curso | cursos.Disciplina | cursos.Documento | cursos.Solicitacao_Doc
}

type PersistenceHandler[T Tables] interface {
	Create(db *gorm.DB) (uint, error)
	List(db *gorm.DB) (*[]T, error)
	Update(db *gorm.DB, ID uint) (*T, error)
	Find(db *gorm.DB, param map[string]interface{}) (*T, error)
	FindAll(db *gorm.DB, param map[string]interface{}) (*[]T, error)
	Delete(db *gorm.DB, ID uint) (int64, error)
	//FindBy(db *gorm.DB, param string, ID interface{}) (*T, error)
	//DeleteBy(db *gorm.DB, cond string, ID uint) (int64, error)
}
