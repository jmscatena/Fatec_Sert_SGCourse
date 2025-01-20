package cursos

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"html"
	"log"
	"strings"
	"time"
)

type Gestao struct {
	gorm.Model
	UID         uuid.UUID  `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"ID"`
	Disciplina  Disciplina `gorm:"foreignkey:DisciplinaID" json:"disciplina"`
	TipoArquivo string     `gorm:"type:text" json:"tipoarquivo"`
	Arquivo     string     `gorm:"type:text" json:"arquivo"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Gestao) Validate() error {
	if p.TipoArquivo == "" || p.TipoArquivo == "null" {
		return errors.New("obrigatório: tipo de arquivo")
	}
	if p.Arquivo == "" || p.Arquivo == "null" {
		return errors.New("obrigatório: arquivo")
	}
	return nil
}
func (p *Gestao) Prepare() {
	p.Disciplina = p.Disciplina //realizar a varredura do registro
	p.TipoArquivo = html.EscapeString(strings.TrimSpace(p.TipoArquivo))
	p.Arquivo = html.EscapeString(strings.TrimSpace(p.TipoArquivo))
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	err := p.Validate()
	if err != nil {
		log.Fatalf("Error during validation:%v", err)
	}
}
func (p *Gestao) Create(db *gorm.DB) (uuid.UUID, error) {
	if verr := p.Validate(); verr != nil {
		return uuid.Nil, verr
	}
	p.Prepare()
	err := db.Debug().Omit("ID").Create(&p).Error
	if err != nil {
		return uuid.Nil, err
	}
	return p.UID, nil
}
func (p *Gestao) Update(db *gorm.DB, uid uuid.UUID) (*Gestao, error) {
	db = db.Debug().Model(&Gestao{}).Where("id = ?", uid).Updates(Gestao{
		Disciplina:  p.Disciplina,
		TipoArquivo: p.TipoArquivo,
		Arquivo:     p.Arquivo})
	if db.Error != nil {
		return nil, db.Error
	}
	return p, nil
}
func (p *Gestao) List(db *gorm.DB) (*[]Gestao, error) {
	Gestaos := []Gestao{}
	err := db.Debug().Model(&Gestao{}).Limit(100).Find(&Gestaos).Error
	//result := db.Find(&Gestaos)
	if err != nil {
		return nil, err
	}
	return &Gestaos, nil
}

func (u *Gestao) Find(db *gorm.DB, param string, uid string) (*Gestao, error) {
	err := db.Debug().Model(Gestao{}).Where(param, uid).Take(&u).Error
	if err != nil {
		return &Gestao{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Gestao{}, errors.New("Material Inexistente")
	}
	return u, nil
}

/*
func (p *Gestao) Find(db *gorm.DB, uid uuid.UUID) (*Gestao, error) {
	err := db.Debug().Model(&Gestao{}).Where("id = ?", uid).Take(&p).Error
	if err != nil {
		return &Gestao{}, err
	}
	return p, nil
}

func (p *Gestao) FindBy(db *gorm.DB, param string, uid ...interface{}) (*[]Gestao, error) {
	Gestaos := []Gestao{}
	params := strings.Split(param, ";")
	uids := uid[0].([]interface{})
	if len(params) != len(uids) {
		return nil, errors.New("condição inválida")
	}
	result := db.Where(strings.Join(params, " AND "), uids...).Find(&Gestaos)
	if result.Error != nil {
		return nil, result.Error
	}
	return &Gestaos, nil
}
*/

func (p *Gestao) Delete(db *gorm.DB, uid uuid.UUID) (int64, error) {
	db = db.Delete(&Gestao{}, "id = ? ", uid)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
func (p *Gestao) DeleteBy(db *gorm.DB, cond string, uid uuid.UUID) (int64, error) {
	result := db.Delete(&Gestao{}, cond+" = ?", uid)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
