package cursos

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Tipo string

const (
	pdf Tipo = "pdf"
	img Tipo = "img"
	doc Tipo = "doc"
)

type Documento struct {
	gorm.Model
	TipoID    uint
	ID        uint64    `gorm:"unique;primaryKey;autoIncrement" json:"ID"`
	Titulo    string    `gorm:"size:255;not null;unique" json:"titulo"`
	Tipo      Tipo      `gorm:"foreignKey:TipoID;references:ID" json:"tipo" validate:"required"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Ativo     bool      `gorm:"default:True;" json:"ativo"`
}

func (p *Documento) Validate() error {
	return nil
}

func (p *Documento) Create(db *gorm.DB) (uint64, error) {
	if verr := p.Validate(); verr != nil {
		return 0, verr
	}
	err := db.Debug().Omit("ID").Create(&p).Error
	if err != nil {
		return 0, err
	}
	return p.ID, nil
}

func (p *Documento) Update(db *gorm.DB, id uint64) (*Documento, error) {
	db = db.Debug().Model(Documento{}).Where("id = ?", id).Updates(Documento{
		Titulo: p.Titulo,
		Tipo:   p.Tipo,
		Ativo:  p.Ativo,
	})

	if db.Error != nil {
		return nil, db.Error
	}
	return p, nil
}

func (p *Documento) List(db *gorm.DB) (*[]Documento, error) {
	Documentos := []Documento{}
	err := db.Debug().Model(&Documento{}).Limit(100).Find(&Documentos).Error
	//result := db.Find(&Documentos)
	if err != nil {
		return nil, err
	}
	return &Documentos, nil
}

func (u *Documento) Find(db *gorm.DB, param string, id string) (*Documento, error) {
	err := db.Debug().Model(Documento{}).Where(param, id).Take(&u).Error
	if err != nil {
		return &Documento{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Documento{}, errors.New("Documento Inexistente")
	}
	return u, nil
}

/*
	func (p *Documento) Find(db *gorm.DB, id uint64) (*Documento, error) {
		err := db.Debug().Model(&Documento{}).Where("id = ?", id).Take(&p).Error
		if err != nil {
			return &Documento{}, err
		}
		return p, nil
	}

	func (p *Documento) FindBy(db *gorm.DB, param string, id ...interface{}) (*[]Documento, error) {
		Documentos := []Documento{}
		params := strings.Split(param, ";")
		ids := id[0].([]interface{})
		if len(params) != len(ids) {
			return nil, errors.New("condição inválida")
		}
		result := db.Where(strings.Join(params, " AND "), ids...).Find(&Documentos)
		if result.Error != nil {
			return nil, result.Error
		}
		return &Documentos, nil
	}
*/
func (p *Documento) Delete(db *gorm.DB, id uint64) (int64, error) {
	db = db.Delete(&Documento{}, "id = ? ", id)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (p *Documento) DeleteBy(db *gorm.DB, cond string, id uint64) (int64, error) {
	result := db.Delete(&Documento{}, cond+" = ?", id)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
