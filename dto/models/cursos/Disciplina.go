package cursos

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jmscatena/Fatec_Sert_SGCourse/dto/models/administrativo"
	"gorm.io/gorm"
	"html"
	"strings"
	"time"
)

type Disciplina struct {
	gorm.Model
	UsuarioID uint
	CursoID   uint
	UID       uuid.UUID              `gorm:"type:uuid;default:uuid_generate_v4()" json:"ID"`
	Nome      string                 `gorm:"size:255;not null;unique" json:"nome"`
	Curso     Curso                  `gorm:"foreignKey:CursoID;references:ID" json:"curso"`
	Semestre  int                    `gorm:"default:-1" json:"semestre"`
	Professor administrativo.Usuario `gorm:"foreignkey:UsuarioID;references:ID, constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"professor"`
	Ativo     bool                   `gorm:"default:True;" json:"ativo"`
}

func (p *Disciplina) Validate() error {

	if p.Nome == "" || p.Nome == "null" {
		return errors.New("obrigatório: Nome")
	}
	if p.Semestre == 0 {
		return errors.New("obrigatório: Semestre de computadores")
	}
	return nil
}
func (p *Disciplina) Prepare(db *gorm.DB) (err error) {
	p.Nome = html.EscapeString(strings.TrimSpace(p.Nome))
	p.Semestre = int(p.Semestre)
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	return
}

func (p *Disciplina) Create(db *gorm.DB) (uuid.UUID, error) {
	if verr := p.Validate(); verr != nil {
		return uuid.Nil, verr
	}
	p.Prepare(db)
	err := db.Debug().Omit("ID").Create(&p).Error
	if err != nil {
		return uuid.Nil, err
	}
	return p.UID, nil
}

func (p *Disciplina) Update(db *gorm.DB, uid uuid.UUID) (*Disciplina, error) {
	p.Prepare(db)
	//err := db.Debug().Model(&Disciplina{}).Where("id = ?", uid).Take(&Disciplina{}).UpdateColumns(
	//	map[string]interface{}
	db = db.Model(Disciplina{}).Where("id = ?", uid).Updates(
		Disciplina{
			Nome:     p.Nome,
			Semestre: p.Semestre,
			Ativo:    p.Ativo,
			Curso:    p.Curso,
		})
	if db.Error != nil {
		return &Disciplina{}, db.Error
	}
	return p, nil
}

func (p *Disciplina) List(db *gorm.DB) (*[]Disciplina, error) {
	Disciplinas := []Disciplina{}
	//err := db.Debug().Model(&Disciplina{}).Limit(100).Find(&Disciplinas).Error
	//result := db.Find(&Disciplinas)
	err := db.Model(&Disciplina{}).Preload("CreatedBy").Preload("UpdatedBy").Preload("Materiais").Find(&Disciplinas).Error
	if err != nil {
		return nil, err
	}
	return &Disciplinas, nil
}
func (u *Disciplina) Find(db *gorm.DB, param string, uid string) (*Disciplina, error) {
	err := db.Debug().Model(Disciplina{}).Where(param, uid).Take(&u).Error
	if err != nil {
		return &Disciplina{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Disciplina{}, errors.New("Laboratorio Inexistente")
	}
	return u, nil
}

/*
	func (p *Disciplina) Find(db *gorm.DB, uid uuid.UUID) (*Disciplina, error) {
		err := db.Debug().Model(&Disciplina{}).Preload("CreatedBy").Preload("UpdatedBy").Preload("Materiais").Where("id = ?", uid).Take(&p).Error
		if err != nil {
			return &Disciplina{}, err
		}
		return p, nil
	}

	func (p *Disciplina) FindBy(db *gorm.DB, param string, uid ...interface{}) (*[]Disciplina, error) {
		Disciplinas := []Disciplina{}
		params := strings.Split(param, ";")
		uids := uid[0].([]interface{})
		if len(params) != len(uids) {
			return nil, errors.New("condição inválida")
		}
		result := db.Model(&Disciplina{}).Preload("CreatedBy").Preload("UpdatedBy").Preload("Materiais").Where(strings.Join(params, " AND "), uids...).Find(&Disciplinas)
		//result := db.Joins("CreatedBy", db.Where(strings.Join(params, " AND "), uids...)).Find(&Disciplinas)
		if result.Error != nil {
			return nil, result.Error
		}
		return &Disciplinas, nil
	}
*/
func (p *Disciplina) Delete(db *gorm.DB, uid uuid.UUID) (int64, error) {
	db = db.Delete(&Disciplina{}, "id = ? ", uid)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
