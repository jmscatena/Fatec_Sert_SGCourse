package handlers

import (
	"github.com/jmscatena/Fatec_Sert_SGCourse/dto/models/administrativo"
	"github.com/jmscatena/Fatec_Sert_SGCourse/dto/models/cursos"
	"gorm.io/gorm"
)

type Tables interface {
	administrativo.Usuario | cursos.Curso | cursos.Disciplina | cursos.Gestao
}

type PersistenceHandler[T Tables] interface {
	Create(db *gorm.DB) (uint64, error)
	List(db *gorm.DB) (*[]T, error)
	Update(db *gorm.DB, uid uint64) (*T, error)
	Find(db *gorm.DB, param string, uid string) (*T, error)
	Delete(db *gorm.DB, uid uint64) (int64, error)
	//FindBy(db *gorm.DB, param string, uid interface{}) (*T, error)
	//DeleteBy(db *gorm.DB, cond string, uid uint64) (int64, error)
}
