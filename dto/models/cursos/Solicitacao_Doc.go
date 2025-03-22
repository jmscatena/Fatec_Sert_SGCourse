package cursos

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Solicitacao_Doc struct {
	// Esta faltando os materiais
	gorm.Model
	DocumentoID  uint
	DisciplinaID uint
	UID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()" json:"ID"`
	Documento    Documento  `gorm:"foreignkey:DocumentoID,references:ID,constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"documento"`
	Disciplina   Disciplina `gorm:"foreignkey:DisciplinaID,references:ID,constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"disciplina"`
	Entrega      bool       `gorm:"size:255;not null;unique" json:"entrega"`
	Prazo        time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"prazo"`
	Ativo        bool       `gorm:"default:True;" json:"ativo"`
	CreatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Solicitacao_Doc) Validate() error {
	return nil
}

func (p *Solicitacao_Doc) Create(db *gorm.DB) (uuid.UUID, error) {
	if verr := p.Validate(); verr != nil {
		return uuid.Nil, verr
	}
	err := db.Debug().Omit("ID").Create(&p).Error
	if err != nil {
		return uuid.Nil, err
	}
	return p.UID, nil
}

func (p *Solicitacao_Doc) Update(db *gorm.DB, uid uuid.UUID) (*Solicitacao_Doc, error) {
	db = db.Debug().Model(Solicitacao_Doc{}).Where("id = ?", uid).Updates(Solicitacao_Doc{
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
	err := db.Debug().Model(&Solicitacao_Doc{}).Limit(100).Find(&Solicitacao_Docs).Error
	//result := db.Find(&Solicitacao_Docs)
	if err != nil {
		return nil, err
	}
	return &Solicitacao_Docs, nil
}

func (u *Solicitacao_Doc) Find(db *gorm.DB, param string, uid string) (*Solicitacao_Doc, error) {
	err := db.Debug().Model(Solicitacao_Doc{}).Where(param, uid).Take(&u).Error
	if err != nil {
		return &Solicitacao_Doc{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Solicitacao_Doc{}, errors.New("Gestao Material Inexistente")
	}
	return u, nil
}

func (p *Solicitacao_Doc) Delete(db *gorm.DB, uid uuid.UUID) (int64, error) {
	db = db.Delete(&Solicitacao_Doc{}, "id = ? ", uid)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (p *Solicitacao_Doc) DeleteBy(db *gorm.DB, cond string, uid uuid.UUID) (int64, error) {
	result := db.Delete(&Solicitacao_Doc{}, cond+" = ?", uid)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
