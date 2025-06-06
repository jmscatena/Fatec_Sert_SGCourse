package cursos

import (
	"errors"
	"github.com/jmscatena/Fatec_Sert_SGCourse/dto/models/administrativo"
	"gorm.io/gorm"
	"html"
	"strings"
	"time"
)

type Disciplina struct {
	gorm.Model
	ID        uint                   `gorm:"primaryKey;autoIncrement" json:"ID"`
	Nome      string                 `gorm:"size:255;not null;unique" json:"nome"`
	CursoID   uint                   `json:"cursoID"`
	Curso     Curso                  `json:"curso"`
	Semestre  int                    `gorm:"default:-1" json:"semestre"`
	UsuarioID uint                   `json:"usuarioID"`
	Usuario   administrativo.Usuario `json:"professor"`
	Ativo     bool                   `gorm:"default:true" json:"ativo"`
}

func (p *Disciplina) Validate() error {

	if p.Nome == "" || p.Nome == "null" {
		return errors.New("obrigatório: Nome")
	}
	if p.Semestre == 0 {
		return errors.New("obrigatório: Semestre de computadores")
	}
	return nil
}
func (p *Disciplina) Prepare() (err error) {
	p.Nome = html.EscapeString(strings.TrimSpace(p.Nome))
	p.Semestre = int(p.Semestre)
	p.CursoID = uint(p.CursoID)
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	return
}

func (p *Disciplina) Create(db *gorm.DB) (uint, error) {
	if verr := p.Validate(); verr != nil {
		return 0, verr
	}
	p.Prepare()
	err := db.Omit("ID").Create(&p).Error
	if err != nil {
		return 0, err
	}
	return p.ID, nil
}

func (p *Disciplina) Update(db *gorm.DB, id uint) (*Disciplina, error) {
	p.Prepare()
	//err := db.Model(&Disciplina{}).Where("id = ?", id).Take(&Disciplina{}).UpdateColumns(
	//	map[string]interface{}
	db = db.Model(Disciplina{}).Where("id = ?", id).Updates(
		Disciplina{
			Nome:     p.Nome,
			Semestre: p.Semestre,
			Ativo:    p.Ativo,
			CursoID:  p.CursoID,
		})
	if db.Error != nil {
		return &Disciplina{}, db.Error
	}
	return p, nil
}

func (p *Disciplina) List(db *gorm.DB) (*[]Disciplina, error) {
	Disciplinas := []Disciplina{}
	//err := db.Model(&Disciplina{}).Limit(100).Find(&Disciplinas).Error
	err := db.
		Preload("Curso").
		Preload("Usuario", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, nome").Omit("email", "professor", "coordenador", "diretor")
		}).
		Find(&Disciplinas).Error

	//result := db.Find(&Disciplinas)
	//err := db.Model(&Disciplina{}).Preload("CreatedBy").Preload("UpdatedBy").Preload("Materiais").Find(&Disciplinas).Error
	if err != nil {
		return nil, err
	}
	return &Disciplinas, nil
}

func (u *Disciplina) Find(db *gorm.DB, params map[string]interface{}) (*Disciplina, error) {
	var err error
	query := db.Model(&Curso{})
	if params != nil {
		for key, value := range params {
			query = query.Where(key, value)
		}
	}
	err = query.Preload("Curso").Preload("Usuario").Find(&u).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Usuario Inexistente")
		}
		return nil, err // Return the original error if it's not RecordNotFound
	}
	return u, nil
}
func (p *Disciplina) FindAll(db *gorm.DB, param map[string]interface{}) (*[]Disciplina, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Disciplina) Delete(db *gorm.DB, id uint) (int64, error) {
	db = db.Delete(&Disciplina{}, "id = ? ", id)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
