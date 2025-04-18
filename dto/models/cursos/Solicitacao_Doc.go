package cursos

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Solicitacao_Doc struct {
	// Esta faltando os materiais
	gorm.Model
	DocumentoID  uint       `json:"documentoID"`
	DisciplinaID uint       `json:"disciplinaID"`
	ID           uint       `gorm:"unique;primaryKey;autoIncrement" json:"ID"`
	Documento    Documento  `json:"documento"`
	Disciplina   Disciplina `json:"disciplina"`
	Entrega      bool       `gorm:"size:255;not null;unique" json:"entrega"`
	Prazo        time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"prazo"`
	Ativo        bool       `gorm:"default:True;" json:"ativo"`
	CreatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Solicitacao_Doc) Validate() error {
	return nil
}

func (p *Solicitacao_Doc) Create(db *gorm.DB) (uint, error) {
	if verr := p.Validate(); verr != nil {
		return 0, verr
	}
	err := db.Debug().Omit("ID").Create(&p).Error
	if err != nil {
		return 0, err
	}
	return p.ID, nil
}

func (p *Solicitacao_Doc) Update(db *gorm.DB, ID uint) (*Solicitacao_Doc, error) {
	db = db.Debug().Model(Solicitacao_Doc{}).Where("id = ?", ID).Updates(Solicitacao_Doc{
		Disciplina: p.Disciplina,
		Documento:  p.Documento,
		Entrega:    p.Entrega,
		Ativo:      p.Ativo,
	})

	if db.Error != nil {
		return nil, db.Error
	}
	return p, nil
}

func (p *Solicitacao_Doc) List(db *gorm.DB) (*[]Solicitacao_Doc, error) {
	Solicitacao_Docs := []Solicitacao_Doc{}
	err := db.Debug().Model(&Solicitacao_Doc{}).Limit(100).Preload("Disciplina").Preload("Documento").Find(&Solicitacao_Docs).Error
	//result := db.Find(&Solicitacao_Docs)
	if err != nil {
		return nil, err
	}
	return &Solicitacao_Docs, nil
}

func (u *Solicitacao_Doc) Find(db *gorm.DB, params map[string]interface{}) (*Solicitacao_Doc, error) {
	var err error
	query := db.Model(&Solicitacao_Doc{})
	if params != nil {
		for key, value := range params {
			query = query.Where(key, value)
		}
	}
	err = query.Find(&u).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Solicitação de Documento Inexistente")
		}
		return nil, err // Return the original error if it's not RecordNotFound
	}
	return u, nil
}

func (p *Solicitacao_Doc) Delete(db *gorm.DB, ID uint) (int64, error) {
	db = db.Delete(&Solicitacao_Doc{}, "id = ? ", ID)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (p *Solicitacao_Doc) DeleteBy(db *gorm.DB, cond string, ID uint) (int64, error) {
	result := db.Delete(&Solicitacao_Doc{}, cond+" = ?", ID)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
