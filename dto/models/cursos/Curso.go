package cursos

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Periodo string

const (
	matutino   Periodo = "matutino"
	vespertino Periodo = "vespertino"
	noturno    Periodo = "noturno"
)

type Curso struct {
	// Esta faltando os materiais
	gorm.Model
	UID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"ID"`
	Nome    string    `gorm:"size:255;not null;unique" json:"nome"`
	Periodo Periodo   `json:"periodo" validate:"required"`
	Ativo   bool      `gorm:"default:True;" json:"ativo"`
}

func (p *Curso) Validate() error {
	return nil
}

func (p *Curso) Create(db *gorm.DB) (uuid.UUID, error) {
	if verr := p.Validate(); verr != nil {
		return uuid.Nil, verr
	}
	err := db.Debug().Omit("ID").Create(&p).Error
	if err != nil {
		return uuid.Nil, err
	}
	return p.UID, nil
}

func (p *Curso) Update(db *gorm.DB, uid uuid.UUID) (*Curso, error) {
	db = db.Debug().Model(Curso{}).Where("id = ?", uid).Updates(Curso{
		Nome:    p.Nome,
		Periodo: p.Periodo,
		Ativo:   p.Ativo,
	})

	if db.Error != nil {
		return nil, db.Error
	}
	return p, nil
}

func (p *Curso) List(db *gorm.DB) (*[]Curso, error) {
	Cursos := []Curso{}
	err := db.Debug().Model(&Curso{}).Limit(100).Find(&Cursos).Error
	//result := db.Find(&Cursos)
	if err != nil {
		return nil, err
	}
	return &Cursos, nil
}

func (u *Curso) Find(db *gorm.DB, param string, uid string) (*Curso, error) {
	err := db.Debug().Model(Curso{}).Where(param, uid).Take(&u).Error
	if err != nil {
		return &Curso{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Curso{}, errors.New("Gestao Material Inexistente")
	}
	return u, nil
}

/*
	func (p *Curso) Find(db *gorm.DB, uid uuid.UUID) (*Curso, error) {
		err := db.Debug().Model(&Curso{}).Where("id = ?", uid).Take(&p).Error
		if err != nil {
			return &Curso{}, err
		}
		return p, nil
	}

	func (p *Curso) FindBy(db *gorm.DB, param string, uid ...interface{}) (*[]Curso, error) {
		Cursos := []Curso{}
		params := strings.Split(param, ";")
		uids := uid[0].([]interface{})
		if len(params) != len(uids) {
			return nil, errors.New("condição inválida")
		}
		result := db.Where(strings.Join(params, " AND "), uids...).Find(&Cursos)
		if result.Error != nil {
			return nil, result.Error
		}
		return &Cursos, nil
	}
*/
func (p *Curso) Delete(db *gorm.DB, uid uuid.UUID) (int64, error) {
	db = db.Delete(&Curso{}, "id = ? ", uid)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (p *Curso) DeleteBy(db *gorm.DB, cond string, uid uuid.UUID) (int64, error) {
	result := db.Delete(&Curso{}, cond+" = ?", uid)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
