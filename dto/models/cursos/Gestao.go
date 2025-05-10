package cursos

import (
	"errors"
	"gorm.io/gorm"
	"html"
	"log"
	"strings"
	"time"
)

type Gestao struct {
	gorm.Model
	ID           uint       `gorm:"unique;primaryKey;autoIncrement" json:"ID"`
	TipoArquivo  string     `gorm:"type:text" json:"tipoarquivo"`
	Arquivo      string     `gorm:"type:text" json:"arquivo"`
	DisciplinaID uint       `json:"disciplinaID"`
	Disciplina   Disciplina `json:"disciplina"`
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
func (p *Gestao) Create(db *gorm.DB) (uint, error) {
	if verr := p.Validate(); verr != nil {
		return 0, verr
	}
	p.Prepare()
	err := db.Debug().Omit("ID").Create(&p).Error
	if err != nil {
		return 0, err
	}
	return p.ID, nil
}
func (p *Gestao) Update(db *gorm.DB, id uint) (*Gestao, error) {
	db = db.Debug().Model(&Gestao{}).Where("id = ?", id).Updates(Gestao{
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
	err := db.Debug().Model(&Gestao{}).Limit(100).Preload("Disciplina").Find(&Gestaos).Error
	//result := db.Find(&Gestaos)
	if err != nil {
		return nil, err
	}
	return &Gestaos, nil
}

func (u *Gestao) Find(db *gorm.DB, param string, ID uint) (*Gestao, error) {
	err := db.Debug().Model(Gestao{}).Where(param, ID).Preload("Disciplina").Take(&u).Error
	if err != nil {
		return &Gestao{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Gestao{}, errors.New("Material Inexistente")
	}
	return u, nil
}

/*
func (p *Gestao) Find(db *gorm.DB, id uint) (*Gestao, error) {
	err := db.Debug().Model(&Gestao{}).Where("id = ?", id).Take(&p).Error
	if err != nil {
		return &Gestao{}, err
	}
	return p, nil
}

func (p *Gestao) FindBy(db *gorm.DB, param string, id ...interface{}) (*[]Gestao, error) {
	Gestaos := []Gestao{}
	params := strings.Split(param, ";")
	ids := id[0].([]interface{})
	if len(params) != len(ids) {
		return nil, errors.New("condição inválida")
	}
	result := db.Where(strings.Join(params, " AND "), ids...).Find(&Gestaos)
	if result.Error != nil {
		return nil, result.Error
	}
	return &Gestaos, nil
}
*/

func (p *Gestao) Delete(db *gorm.DB, id uint) (int64, error) {
	db = db.Delete(&Gestao{}, "id = ? ", id)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
func (p *Gestao) DeleteBy(db *gorm.DB, cond string, id uint) (int64, error) {
	result := db.Delete(&Gestao{}, cond+" = ?", id)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
