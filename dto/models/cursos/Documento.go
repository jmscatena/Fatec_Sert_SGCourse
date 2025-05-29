package cursos

import (
	"errors"
	"gorm.io/gorm"
)

type Documento struct {
	gorm.Model
	ID     uint   `gorm:"unique;primaryKey;autoIncrement" json:"ID"`
	Titulo string `gorm:"size:255;not null;" json:"titulo"`
	Tipo   string `json:"tipo" validate:"required"`
	Ativo  bool   `gorm:"default:True;" json:"ativo"`
}

func (p *Documento) Validate() error {
	return nil
}

func (p *Documento) Create(db *gorm.DB) (uint, error) {
	if verr := p.Validate(); verr != nil {
		return 0, verr
	}
	err := db.Omit("ID").Create(&p).Error
	if err != nil {
		return 0, err
	}
	return p.ID, nil
}

func (p *Documento) Update(db *gorm.DB, id uint) (*Documento, error) {
	db = db.Model(Documento{}).Where("id = ?", id).Updates(Documento{
		Titulo: p.Titulo,
		Tipo:   p.Tipo,
		Ativo:  p.Ativo,
	})

	if db.Error != nil {
		return nil, db.Error
	}
	return p, nil
}

func (p *Documento) List(db *gorm.DB) (*[]Documento, error) {
	Documentos := []Documento{}
	err := db.Model(&Documento{}).Limit(100).Find(&Documentos).Error
	//result := db.Find(&Documentos)
	if err != nil {
		return nil, err
	}
	return &Documentos, nil
}

func (u *Documento) Find(db *gorm.DB, params map[string]interface{}) (*Documento, error) {
	var err error
	query := db.Model(&Documento{})
	if params != nil {
		for key, value := range params {
			query = query.Where(key, value)
		}
	}
	err = query.Find(&u).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Documento Inexistente")
		}
		return nil, err // Return the original error if it's not RecordNotFound
	}
	return u, nil
}
func (p *Documento) FindAll(db *gorm.DB, param map[string]interface{}) (*[]Documento, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Documento) Delete(db *gorm.DB, id uint) (int64, error) {
	db = db.Delete(&Documento{}, "id = ? ", id)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (p *Documento) DeleteBy(db *gorm.DB, cond string, id uint) (int64, error) {
	result := db.Delete(&Documento{}, cond+" = ?", id)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
