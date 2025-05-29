package cursos

import (
	"errors"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const SharedFolderPath = "./static/requests/"

type Entrega_Doc struct {
	gorm.Model
	SolicitacaoID uint            `json:"solicitacaoID"`
	ID            uint            `gorm:"unique;primaryKey;autoIncrement" json:"ID"`
	Solicitacao   Solicitacao_Doc `json:"solicitacao"`
	Arquivo       string          `gorm:"type:varchar(255)" json:"arq"`
}

func (p *Entrega_Doc) Validate() error {
	if p.SolicitacaoID == 0 {
		return errors.New("obrigatório: SolicitacaoID")
	}
	if p.Arquivo == "" {
		return errors.New("obrigatório: Nome do arquivo")
	}
	return nil
}
func (p *Entrega_Doc) Prepare() (err error) {
	p.Solicitacao = p.Solicitacao
	p.Arquivo = p.Arquivo
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return
}
func (p *Entrega_Doc) SaveFile(fileHeader *multipart.FileHeader, fileName string) (string, error) {

	fullName := fileHeader.Filename
	// 1. Check if the file extension is allowed
	ext := filepath.Ext(fullName)
	allowedExt := map[string]bool{
		".pdf":  true,
		".jpg":  true,
		".png":  true,
		".docx": true,
		".doc":  true,
		".xls":  true,
		".xlsx": true,
		".ppt":  true,
		".pptx": true,
	}
	if !allowedExt[strings.ToLower(ext)] {
		return "", errors.New("arquivo: extensão não permitida, use apenas PDF, JPG, PNG, DOC, DOCX, PPT, PPTX, XLS, XLSX")
	}
	fullName = fileName + ext
	// 2. Open the uploaded file (Gin provides an *os.File-like interface)
	uploadedFile, err := fileHeader.Open()
	if err != nil {
		return "", errors.New("arquivo: erro ao abrir o arquivo")
	}
	defer uploadedFile.Close() // Ensure the file is closed

	// 3. Define the upload directory and create it if it doesn't exist
	uploadDir := "../static/requests/files" // A different directory to distinguish
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err = os.Mkdir(uploadDir, 0755); err != nil {
			return "", errors.New("diretorio: falha ao criar diretorio")
		}
	}

	// Use filepath.Base to prevent directory traversal attacks
	filePath := filepath.Join(uploadDir, filepath.Base(fullName))

	// 5. Create the destination file on disk
	dstFile, err := os.Create(filePath)
	if err != nil {
		return "", errors.New("erro: ao criar arquivo")
	}
	defer dstFile.Close() // Ensure the destination file is closed

	// 6. Copy the content from the uploaded file to the destination file
	// Use io.Copy to efficiently copy large files without loading entirely into memory
	_, err = io.Copy(dstFile, uploadedFile)
	if err != nil {
		return "", errors.New("erro: ao salvar arquivo")
	}
	return dstFile.Name(), nil
}
func (p *Entrega_Doc) Create(db *gorm.DB) (uint, error) {
	if verr := p.Validate(); verr != nil {
		return 0, verr
	}
	p.Prepare()
	err := db.Omit("ID").Create(&p).Error
	if err != nil {
		return 0, err
	}
	err = db.Model(&Solicitacao_Doc{}).Where("id = ?", p.SolicitacaoID).Updates(
		Solicitacao_Doc{
			Ativo:   true,
			Entrega: true,
		}).Error
	if err != nil {
		return 0, err
	}
	return p.ID, nil
}
func (p *Entrega_Doc) Update(db *gorm.DB, id uint) (*Entrega_Doc, error) {
	p.Prepare()
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
	//err := db.Model(&Entrega_Doc{}).Limit(100).Find(&Entrega_Docs).Error
	//result := db.Find(&Entrega_Docs)
	//err := db.Model(&Entrega_Doc{}).Preload("Materiais").Find(&Entrega_Docs).Error
	err := db.Preload("Curso").Preload("Solicitacao").Find(&Entrega_Docs).Error

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
