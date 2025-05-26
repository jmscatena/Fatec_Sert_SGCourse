package cursos

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"time"
)

const SharedFolderPath = "/shared/uploads/"

type Entrega_Doc struct {
	gorm.Model
	SolicitacaoID uint            `json:"solicitacaoID"`
	ID            uint            `gorm:"unique;primaryKey;autoIncrement" json:"ID"`
	Solicitacao   Solicitacao_Doc `json:"solicitacao"`
	Arquivo       string          `gorm:"type:varchar(255)" json:"arq"`
	FileName      string          `gorm:"type:varchar(255)" json:"fileName"`
}

func (p *Entrega_Doc) Validate() error {
	if p.SolicitacaoID == 0 {
		return errors.New("obrigatório: SolicitacaoID")
	}
	if p.FileName == "" {
		return errors.New("obrigatório: Nome do arquivo")
	}
	return nil
}
func (p *Entrega_Doc) Prepare(db *gorm.DB) (err error) {
	p.Solicitacao = p.Solicitacao
	p.Arquivo = p.Arquivo
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return
}
func (p *Entrega_Doc) SaveFile(fileBytes []byte, originalName string) (string, error) {
	ext := filepath.Ext(originalName)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(SharedFolderPath, fileName)

	err := os.MkdirAll(SharedFolderPath, 0755)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(filePath, fileBytes, 0644)
	if err != nil {
		return "", err
	}

	return fileName, nil
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
func (u *Entrega_Doc) Find(db *gorm.DB, params map[string]interface{}) (*Entrega_Doc, error) {
	var err error
	query := db.Model(&Entrega_Doc{})
	query = query.
		Preload("Solicitacao", func(db *gorm.DB) *gorm.DB {
			return db.Select("id,documento,disciplina,entrega,prazo,ativo").Omit("CreatedAt", "UpdatedAt", "DeletedAt")
		}).
		Preload("Disciplina.Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id,nome").Omit("CreatedAt", "UpdatedAt", "DeletedAt") }).
		Preload("Disciplina.Curso", func(db *gorm.DB) *gorm.DB {
			return db.Select("id,nome,periodo").Omit("CreatedAt", "UpdatedAt", "DeletedAt")
		}).
		Preload("Documento", func(db *gorm.DB) *gorm.DB {
			return db.Select("id,titulo,tipo").Omit("CreatedAt", "UpdatedAt", "DeletedAt")
		})

	if params != nil {
		for key, value := range params {
			if key == "email" {
				query = query.Preload("Solicitacao.Disciplina.Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id,nome,email").Where("email = ?", value) })
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

func (p *Entrega_Doc) FindAll(db *gorm.DB, param map[string]interface{}) (*[]Entrega_Doc, error) {
	var err error
	var Entrega_Docs []Entrega_Doc
	query := db.Model(&[]Entrega_Doc{})
	query = query.
		Preload("Solicitacao", func(db *gorm.DB) *gorm.DB {
			return db.Select("id,documento,disciplina,entrega,prazo,ativo").Omit("CreatedAt", "UpdatedAt", "DeletedAt")
		}).
		Preload("Disciplina.Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id,nome").Omit("CreatedAt", "UpdatedAt", "DeletedAt") }).
		Preload("Disciplina.Curso", func(db *gorm.DB) *gorm.DB {
			return db.Select("id,nome,periodo").Omit("CreatedAt", "UpdatedAt", "DeletedAt")
		}).
		Preload("Documento", func(db *gorm.DB) *gorm.DB {
			return db.Select("id,titulo,tipo").Omit("CreatedAt", "UpdatedAt", "DeletedAt")
		})

	if param != nil {
		for key, value := range param {
			if key == "email" {
				query = query.Preload("Solicitacao.Disciplina.Usuario", func(db *gorm.DB) *gorm.DB { return db.Select("id,nome,email").Where("email = ?", value) })
			} else {
				query = query.Where(key, value)
			}

		}
	}
	err = query.Omit("CreatedAt", "UpdatedAt", "DeletedAt").Find(&Entrega_Docs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Solicitação de Documento Inexistente")
		}
		return nil, err // Return the original error if it's not RecordNotFound
	}
	return &Entrega_Docs, nil
}

func (p *Entrega_Doc) Delete(db *gorm.DB, id uint) (int64, error) {
	db = db.Delete(&Entrega_Doc{}, "id = ? ", id)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
