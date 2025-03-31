package cursos

import (
	"errors"
	"github.com/jmscatena/Fatec_Sert_SGCourse/dto/models/administrativo"
	"gorm.io/gorm"
	"html"
	"strings"
	"time"
)

type Disciplina struct {
	gorm.Model
	ID        uint `gorm:"unique;primaryKey;autoIncrement" json:"ID"`
	UsuarioID uint
	CursoID   uint
	Nome      string                 `gorm:"size:255;not null;unique" json:"nome"`
	Curso     Curso                  `gorm:"foreignKey:CursoID;references:ID" json:"curso"`
	Semestre  int                    `gorm:"default:-1" json:"semestre"`
	Usuario   administrativo.Usuario `gorm:"foreignkey:UsuarioID;references:ID" json:"professor"`
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

func (p *Disciplina) Create(db *gorm.DB) (uint, error) {
	if verr := p.Validate(); verr != nil {
		return 0, verr
	}
	p.Prepare(db)
	err := db.Debug().Omit("ID").Create(&p).Error
	if err != nil {
		return 0, err
	}
	return p.ID, nil
}

func (p *Disciplina) Update(db *gorm.DB, id uint) (*Disciplina, error) {
	p.Prepare(db)
	//err := db.Debug().Model(&Disciplina{}).Where("id = ?", id).Take(&Disciplina{}).UpdateColumns(
	//	map[string]interface{}
	db = db.Model(Disciplina{}).Where("id = ?", id).Updates(
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
func (u *Disciplina) Find(db *gorm.DB, param string, ID uint) (*Disciplina, error) {
	err := db.Debug().Model(Disciplina{}).Where(param, ID).Take(&u).Error
	if err != nil {
		return &Disciplina{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Disciplina{}, errors.New("Laboratorio Inexistente")
	}
	return u, nil
}

/*
	func (p *Disciplina) Find(db *gorm.DB, id uint) (*Disciplina, error) {
		err := db.Debug().Model(&Disciplina{}).Preload("CreatedBy").Preload("UpdatedBy").Preload("Materiais").Where("id = ?", id).Take(&p).Error
		if err != nil {
			return &Disciplina{}, err
		}
		return p, nil
	}

	func (p *Disciplina) FindBy(db *gorm.DB, param string, id ...interface{}) (*[]Disciplina, error) {
		Disciplinas := []Disciplina{}
		params := strings.Split(param, ";")
		ids := id[0].([]interface{})
		if len(params) != len(ids) {
			return nil, errors.New("condição inválida")
		}
		result := db.Model(&Disciplina{}).Preload("CreatedBy").Preload("UpdatedBy").Preload("Materiais").Where(strings.Join(params, " AND "), ids...).Find(&Disciplinas)
		//result := db.Joins("CreatedBy", db.Where(strings.Join(params, " AND "), ids...)).Find(&Disciplinas)
		if result.Error != nil {
			return nil, result.Error
		}
		return &Disciplinas, nil
	}
*/
func (p *Disciplina) Delete(db *gorm.DB, id uint) (int64, error) {
	db = db.Delete(&Disciplina{}, "id = ? ", id)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
