package administrativo

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"html"
	"log"
	"strings"
	"time"
)

type Usuario struct {
	gorm.Model
	ID          uint   `gorm:"unique;primaryKey;autoIncrement" json:"ID"`
	Nome        string `gorm:"size:255;not null;unique" json:"nome"`
	Email       string `gorm:"unique;size:100;not null,email;" json:"email"`
	Senha       string `gorm:"size:1024;not null;" json:"senha"`
	Ativo       bool   `gorm:"default:true;" json:"ativo"`
	Diretor     bool   `gorm:"default:false" json:"diretor"`
	Coordenador bool   `gorm:"default:false" json:"coordenador"`
	Professor   bool   `gorm:"default:false" json:"professor"`
}

func (u *Usuario) Create(db *gorm.DB) (uint, error) {
	if verr := u.Validate("insert"); verr != nil {
		return 0, verr
	}
	u.Prepare()
	err := db.Debug().Omit("ID").Create(&u).Error
	if err != nil {
		return 0, err
	}
	return u.ID, nil
}

func (u *Usuario) Update(db *gorm.DB, ID uint) (*Usuario, error) {
	if verr := u.Validate("insert"); verr != nil {
		println(verr)
		return nil, verr
	}
	u.Prepare()
	db = db.Model(Usuario{}).Where("id = ?", ID).Updates(Usuario{
		Senha: u.Senha,
		Nome:  u.Nome,
		Email: u.Email,
	})

	/*db = db.Debug().Model(&Usuario{}).Where("id = ?", ID).Take(&Usuario{}).UpdateColumns(
		map[string]interface{}{
			"Senha": u.Senha,
			"Nome":  u.Nome,
			"Email": u.Email,
			//"updated_at": time.Now(),
		},
	)*/
	if db.Error != nil {
		return &Usuario{}, db.Error
	}
	err := db.Debug().Model(&Usuario{}).Where("id = ?", ID).Take(&u).Error
	if err != nil {
		return &Usuario{}, err
	}
	return u, nil
}

func (u *Usuario) List(db *gorm.DB) (*[]Usuario, error) {
	Usuarios := []Usuario{}
	err := db.Debug().Model(&Usuario{}).Find(&Usuarios).Error
	if err != nil {
		return nil, err
	}
	return &Usuarios, err
}

func (u *Usuario) Find(db *gorm.DB, params map[string]interface{}) (*Usuario, error) {
	var err error
	query := db.Model(&Usuario{})
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

func (u *Usuario) Delete(db *gorm.DB, ID uint) (int64, error) {
	db = db.Debug().Where("id = ?", ID).Delete(&Usuario{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (u *Usuario) DeleteBy(db *gorm.DB, cond string, ID interface{}) (int64, error) {
	result := db.Delete(&Usuario{}, cond+" = ?", ID)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (u *Usuario) Validate(action string) error {
	/*
		if u.Nome == "" {
			return errors.New("obrigatório: nome do usuário")
		}
		if u.Email == "" {
			return errors.New("obrigatório: email")
		}
		if u.Senha == "" {
			return errors.New("obrigatório: senha")
		}
		return nil
	*/
	return nil

}

func Hash(Senha string) []byte {
	hash, _ := bcrypt.GenerateFromPassword([]byte(Senha), bcrypt.DefaultCost)
	return hash
}

func VerifyPassword(hashedSenha string, senha string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedSenha), []byte(senha))
}

func (u *Usuario) Prepare() {
	u.Nome = html.EscapeString(strings.TrimSpace(u.Nome))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))

	u.Senha = string(Hash(u.Senha))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	err := u.Validate("padrao")
	if err != nil {
		log.Fatalf("Error during validation:%v", err)
	}
}
