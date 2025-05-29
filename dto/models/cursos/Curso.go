package cursos

import (
	"errors"
	"gorm.io/gorm"
)

type Curso struct {
	// Esta faltando os materiais
	gorm.Model
	ID      uint   `gorm:"unique;primaryKey;autoIncrement" json:"ID"`
	Nome    string `gorm:"size:255;not null;unique" json:"nome"`
	Periodo string `json:"periodo" validate:"required"`
	Ativo   bool   `gorm:"default:True;" json:"ativo"`
}

func (p *Curso) FindAll(db *gorm.DB, param map[string]interface{}) (*[]Curso, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Curso) Validate() error {
	return nil
}

func (p *Curso) Create(db *gorm.DB) (uint, error) {
	if verr := p.Validate(); verr != nil {
		return 0, verr
	}
	err := db.Omit("ID").Create(&p).Error
	if err != nil {
		return 0, err
	}
	return p.ID, nil
}

func (p *Curso) Update(db *gorm.DB, id uint) (*Curso, error) {
	db = db.Model(Curso{}).Where("id = ?", id).Updates(Curso{
		Nome:    p.Nome,
		Periodo: p.Periodo,
		Ativo:   p.Ativo,
	})
	if db.Error != nil {
		return nil, db.Error
	}
	return p, nil
}

func (p *Curso) List(db *gorm.DB) (*[]Curso, error) {
	Cursos := []Curso{}
	err := db.Model(&Curso{}).Limit(100).Find(&Cursos).Error
	//result := db.Find(&Cursos)
	if err != nil {
		return nil, err
	}
	return &Cursos, nil
}

/*
	func (u *Curso) Find(db *gorm.DB, param string, ID uint) (*Curso, error) {
		err := db.Model(Curso{}).Where(param, ID).Take(&u).Error
		if err != nil {
			return &Curso{}, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &Curso{}, errors.New("Gestao Material Inexistente")
		}
		return u, nil
	}
*/
func (u *Curso) Find(db *gorm.DB, params map[string]interface{}) (*Curso, error) {
	var err error
	query := db.Model(&Curso{})
	if params != nil {
		for key, value := range params {
			query = query.Where(key, value)
		}
	}
	err = query.Find(&u).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Usuario Inexistente")
		}
		return nil, err // Return the original error if it's not RecordNotFound
	}
	return u, nil
}

/*
	func (p *Curso) Find(db *gorm.DB, id uint) (*Curso, error) {
		err := db.Model(&Curso{}).Where("id = ?", id).Take(&p).Error
		if err != nil {
			return &Curso{}, err
		}
		return p, nil
	}

	func (p *Curso) FindBy(db *gorm.DB, param string, id ...interface{}) (*[]Curso, error) {
		Cursos := []Curso{}
		params := strings.Split(param, ";")
		ids := id[0].([]interface{})
		if len(params) != len(ids) {
			return nil, errors.New("condição inválida")
		}
		result := db.Where(strings.Join(params, " AND "), ids...).Find(&Cursos)
		if result.Error != nil {
			return nil, result.Error
		}
		return &Cursos, nil
	}
*/
func (p *Curso) Delete(db *gorm.DB, id uint) (int64, error) {
	db = db.Delete(&Curso{}, "id = ? ", id)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (p *Curso) DeleteBy(db *gorm.DB, cond string, id uint) (int64, error) {
	result := db.Delete(&Curso{}, cond+" = ?", id)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
