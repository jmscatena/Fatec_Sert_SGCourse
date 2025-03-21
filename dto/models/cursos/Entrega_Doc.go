package cursos

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Entrega_Doc struct {
	gorm.Model
	UID         uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4()" json:"ID"`
	Solicitacao Solicitacao_Doc `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"json:"Solicitacao"`
	Arquivo     string          `gorm:"type:text" json:"arquivo"`
	CreatedAt   time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Entrega_Doc) Validate() error {
	return nil
}
func (p *Entrega_Doc) Prepare(db *gorm.DB) (err error) {
	p.Solicitacao = p.Solicitacao
	p.Arquivo = p.Arquivo
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return
}

func (p *Entrega_Doc) Create(db *gorm.DB) (uuid.UUID, error) {
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

func (p *Entrega_Doc) Update(db *gorm.DB, uid uuid.UUID) (*Entrega_Doc, error) {
	p.Prepare(db)
	//err := db.Debug().Model(&Entrega_Doc{}).Where("id = ?", uid).Take(&Entrega_Doc{}).UpdateColumns(
	//	map[string]interface{}
	db = db.Model(Entrega_Doc{}).Where("id = ?", uid).Updates(
		Entrega_Doc{
			Solicitacao: p.Solicitacao,
			Arquivo:     p.Arquivo,
			UpdatedAt:   time.Now(),
		})
	if db.Error != nil {
		return &Entrega_Doc{}, db.Error
	}
	return p, nil
}

func (p *Entrega_Doc) List(db *gorm.DB) (*[]Entrega_Doc, error) {
	Entrega_Docs := []Entrega_Doc{}
	//err := db.Debug().Model(&Entrega_Doc{}).Limit(100).Find(&Entrega_Docs).Error
	//result := db.Find(&Entrega_Docs)
	err := db.Model(&Entrega_Doc{}).Preload("CreatedBy").Preload("UpdatedBy").Preload("Materiais").Find(&Entrega_Docs).Error
	if err != nil {
		return nil, err
	}
	return &Entrega_Docs, nil
}
func (u *Entrega_Doc) Find(db *gorm.DB, param string, uid string) (*Entrega_Doc, error) {
	err := db.Debug().Model(Entrega_Doc{}).Where(param, uid).Take(&u).Error
	if err != nil {
		return &Entrega_Doc{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Entrega_Doc{}, errors.New("Laboratorio Inexistente")
	}
	return u, nil
}

/*
	func (p *Entrega_Doc) Find(db *gorm.DB, uid uuid.UUID) (*Entrega_Doc, error) {
		err := db.Debug().Model(&Entrega_Doc{}).Preload("CreatedBy").Preload("UpdatedBy").Preload("Materiais").Where("id = ?", uid).Take(&p).Error
		if err != nil {
			return &Entrega_Doc{}, err
		}
		return p, nil
	}

	func (p *Entrega_Doc) FindBy(db *gorm.DB, param string, uid ...interface{}) (*[]Entrega_Doc, error) {
		Entrega_Docs := []Entrega_Doc{}
		params := strings.Split(param, ";")
		uids := uid[0].([]interface{})
		if len(params) != len(uids) {
			return nil, errors.New("condição inválida")
		}
		result := db.Model(&Entrega_Doc{}).Preload("CreatedBy").Preload("UpdatedBy").Preload("Materiais").Where(strings.Join(params, " AND "), uids...).Find(&Entrega_Docs)
		//result := db.Joins("CreatedBy", db.Where(strings.Join(params, " AND "), uids...)).Find(&Entrega_Docs)
		if result.Error != nil {
			return nil, result.Error
		}
		return &Entrega_Docs, nil
	}
*/
func (p *Entrega_Doc) Delete(db *gorm.DB, uid uuid.UUID) (int64, error) {
	db = db.Delete(&Entrega_Doc{}, "id = ? ", uid)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
