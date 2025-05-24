package cursos

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Solicitacao_Doc struct {
	gorm.Model
	DocumentoID  uint       `json:"documentoID"`
	DisciplinaID uint       `json:"disciplinaID"`
	CursoID      uint       `json:"cursoID"`
	SemestreID   uint       `json:"semestreID"`
	ID           uint       `gorm:"unique;primaryKey;autoIncrement" json:"ID"`
	Documento    Documento  `json:"documento"`
	Disciplina   Disciplina `json:"disciplina"`
	Entrega      bool       `gorm:"default:false;not null" json:"entrega"`
	Prazo        time.Time  `gorm:"column:prazo;type:date;not null" json:"prazo,omitempty" time_format:"2006-01-02" example:"2006-01-02"`
	Ativo        bool       `gorm:"default:True;" json:"ativo"`
}

func (p *Solicitacao_Doc) Prepare() (err error) {
	if !p.Ativo {
		p.Ativo = true
	}
	p.UpdatedAt = time.Now()
	if !p.Prazo.IsZero() {
		prazoStr := p.Prazo.Format("2006-01-02")
		parsedTime, err := time.Parse("2006-01-02", prazoStr)
		if err != nil {
			return errors.New("invalid date format for Prazo, expected YYYY-MM-DD")
		}
		p.Prazo = parsedTime
	}
	return
}

func (p *Solicitacao_Doc) Validate(db *gorm.DB) error {
	if p.DocumentoID == 0 {
		return errors.New("DocumentoID is required")
	}
	var err error
	solicitacaoDocs := []Solicitacao_Doc{}
	query := db.Model(&Solicitacao_Doc{}).
		Where("documento_id = ?", p.DocumentoID).
		Where("disciplina_id = ?", p.DisciplinaID).
		Where("curso_id = ?", p.CursoID).
		Where("semestre_id = ?", p.SemestreID)
	err = query.Find(&solicitacaoDocs).Error
	if err != nil {
		return err
	}
	if len(solicitacaoDocs) > 0 {
		return errors.New("documento já solicitado para esta disciplina")
	}
	return nil
}

func (p *Solicitacao_Doc) Create(db *gorm.DB) (uint, error) {

	if verr := p.Validate(db); verr != nil {
		return 0, verr
	}
	perr := p.Prepare()
	if perr != nil {
		return 0, perr
	}
	disciplinas := []Disciplina{}
	var err error
	query := db.Debug().Model(&Disciplina{})

	if p.CursoID == 0 {
		if p.SemestreID > 0 {
			query = query.Where("semestre = ?", p.SemestreID)
		}
	} else if p.CursoID > 0 {
		query = query.Where("curso_id = ?", p.CursoID)
		if p.SemestreID > 0 {
			query = query.Where("semestre = ?", p.SemestreID)

			if p.DisciplinaID > 0 {
				query = query.Where("id = ?", p.DisciplinaID)
			}
		}
	}

	err = query.Preload("Curso").Preload("Usuario").Find(&disciplinas).Error
	if err != nil {
		return 0, err
	}
	fmt.Println(disciplinas)

	for _, disciplina := range disciplinas {
		p.DisciplinaID = disciplina.ID
		newSolicitation := *p
		err := db.Debug().Model(&Solicitacao_Doc{}).Omit("ID").Create(&newSolicitation).Error
		if err != nil {
			return 0, err
		}
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
	//[]Solicitacao_Doc{}
	err := db.Debug().
		Model(&Solicitacao_Doc{}).
		Limit(100).
		Preload("Documento", func(db *gorm.DB) *gorm.DB { return db.Select("id,titulo,tipo") }).
		Preload("Disciplina.Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id,nome") }).
		Preload("Disciplina.Curso", func(db *gorm.DB) *gorm.DB { return db.Select("id,nome,periodo") }).
		Select("id, documento_id, disciplina_id, curso_id, semestre_id, entrega, prazo, ativo").
		Omit("CreatedAt", "UpdatedAt", "DeletedAt").
		Find(&Solicitacao_Docs).Error
	fmt.Println(Solicitacao_Docs[0])
	if err != nil {
		return nil, err
	}

	return &Solicitacao_Docs, nil

}

func (u *Solicitacao_Doc) Find(db *gorm.DB, params map[string]interface{}) (*Solicitacao_Doc, error) {
	var err error
	query := db.Model(&Solicitacao_Doc{})
	if params != nil {
		query = query.Preload("Documento", func(db *gorm.DB) *gorm.DB {
			return db.Select("id,titulo,tipo").Omit("CreatedAt", "UpdatedAt", "DeletedAt")
		}).
			Preload("Disciplina.Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id,nome").Omit("CreatedAt", "UpdatedAt", "DeletedAt") }).
			Preload("Disciplina.Curso", func(db *gorm.DB) *gorm.DB {
				return db.Select("id,nome,periodo").Omit("CreatedAt", "UpdatedAt", "DeletedAt")
			}).
			Select("id, documento_id, disciplina_id, curso_id, semestre_id, entrega, prazo, ativo")
		for key, value := range params {
			if key == "email" {
				//query = query.Preload("Disciplina.Usuario").Where("disciplinas.usuario.email = ?", value)
				query = query.Preload("Disciplina.Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id,nome,email").Where("email = ?", value) })
				//Preload("Disciplina.Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id,nome") }).

			} else {
				query = query.Where(key, value)
			}

		}
	}

	err = query.Omit("CreatedAt", "UpdatedAt", "DeletedAt").Find(&u).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Solicitação de Documento Inexistente")
		}
		return nil, err // Return the original error if it's not RecordNotFound
	}
	return u, nil
}

func (u *Solicitacao_Doc) FindAll(db *gorm.DB, params map[string]interface{}) (*[]Solicitacao_Doc, error) {
	var err error
	Solicitacao_Docs := []Solicitacao_Doc{}
	query := db.Model(&Solicitacao_Doc{})
	if params != nil {
		query = query.Preload("Documento", func(db *gorm.DB) *gorm.DB {
			return db.Select("id,titulo,tipo").Omit("CreatedAt", "UpdatedAt", "DeletedAt")
		}).
			Preload("Disciplina.Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id,nome").Omit("CreatedAt", "UpdatedAt", "DeletedAt") }).
			Preload("Disciplina.Curso", func(db *gorm.DB) *gorm.DB {
				return db.Select("id,nome,periodo").Omit("CreatedAt", "UpdatedAt", "DeletedAt")
			}).
			Select("id, documento_id, disciplina_id, curso_id, semestre_id, entrega, prazo, ativo")
		for key, value := range params {
			if key == "email" {
				//query = query.Preload("Disciplina.Usuario").Where("disciplinas.usuario.email = ?", value)
				query = query.Preload("Disciplina.Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id,nome,email").Where("email = ?", value) })
				//Preload("Disciplina.Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id,nome") }).

			} else {
				query = query.Where(key, value)
			}

		}
	}

	err = query.Omit("CreatedAt", "UpdatedAt", "DeletedAt").Find(&Solicitacao_Docs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Solicitação de Documento Inexistente")
		}
		return nil, err // Return the original error if it's not RecordNotFound
	}
	return &Solicitacao_Docs, nil
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
