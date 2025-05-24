package cursos

import (
	"errors"
	"gorm.io/gorm"
)

type Entrega_Doc struct {
	gorm.Model
	SolicitacaoID uint            `json:"solicitacaoID"`
	ID            uint            `gorm:"unique;primaryKey;autoIncrement" json:"ID"`
	Solicitacao   Solicitacao_Doc `json:"solicitacao"`
	Arquivo       string          `gorm:"type:text" json:"arquivo"`
}

func (p *Entrega_Doc) Validate() error {
	return nil
}
func (p *Entrega_Doc) Prepare(db *gorm.DB) (err error) {
	p.Solicitacao = p.Solicitacao
	p.Arquivo = p.Arquivo
	return
}

func (p *Entrega_Doc) Create(db *gorm.DB) (uint, error) {
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

func (p *Entrega_Doc) Update(db *gorm.DB, id uint) (*Entrega_Doc, error) {
	p.Prepare(db)
	//err := db.Debug().Model(&Entrega_Doc{}).Where("id = ?", id).Take(&Entrega_Doc{}).UpdateColumns(
	//	map[string]interface{}
	db = db.Model(Entrega_Doc{}).Where("id = ?", id).Updates(
		Entrega_Doc{
			Solicitacao: p.Solicitacao,
			Arquivo:     p.Arquivo,
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
	//err := db.Model(&Entrega_Doc{}).Preload("Materiais").Find(&Entrega_Docs).Error
	err := db.Debug().Preload("Curso").Preload("Solicitacao").Find(&Entrega_Docs).Error

	if err != nil {
		return nil, err
	}
	return &Entrega_Docs, nil
}
func (u *Entrega_Doc) Find(db *gorm.DB, param string, id uint) (*Entrega_Doc, error) {
	err := db.Debug().Model(Entrega_Doc{}).Where(param, id).Preload("Solicitacao").Take(&u).Error
	if err != nil {
		return &Entrega_Doc{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Entrega_Doc{}, errors.New("Laboratorio Inexistente")
	}
	return u, nil
}
func (p *Entrega_Doc) FindAll(db *gorm.DB, param map[string]interface{}) (*[]Entrega_Doc, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Entrega_Doc) Delete(db *gorm.DB, id uint) (int64, error) {
	db = db.Delete(&Entrega_Doc{}, "id = ? ", id)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
