package administrativo

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"html"
	"log"
	"strings"
	"time"
)

type Usuario struct {
	gorm.Model
	UID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"UID"`
	Nome      string    `gorm:"size:255;not null;unique" json:"nome"`
	Email     string    `gorm:"size:100;not null,email;" json:"email"`
	Senha     string    `gorm:"size:100;not null;" json:"-"`
	Ativo     bool      `gorm:"default:True;" json:"ativo"`
	Professor bool      `gorm:"default:False;"`
	Gestor    bool      `gorm:"default:False;"`
}

func (u *Usuario) Create(db *gorm.DB) (uuid.UUID, error) {
	if verr := u.Validate("insert"); verr != nil {
		return uuid.Nil, verr
	}
	u.Prepare()
	err := db.Debug().Omit("ID").Create(&u).Error
	if err != nil {
		return uuid.Nil, err
	}
	return u.UID, nil
}

func (u *Usuario) Update(db *gorm.DB, uid uuid.UUID) (*Usuario, error) {

	if verr := u.Validate("insert"); verr != nil {
		return nil, verr
	}
	u.Prepare()
	db = db.Model(Usuario{}).Where("id = ?", uid).Updates(Usuario{
		Senha:     u.Senha,
		Nome:      u.Nome,
		Email:     u.Email,
		Professor: u.Professor,
		Gestor:    u.Gestor,
	})

	/*db = db.Debug().Model(&Usuario{}).Where("id = ?", uid).Take(&Usuario{}).UpdateColumns(
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
	err := db.Debug().Model(&Usuario{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Usuario{}, err
	}
	return u, nil
}

func (u *Usuario) List(db *gorm.DB) (*[]Usuario, error) {
	Usuarios := []Usuario{}
	err := db.Debug().Model(&Usuario{}).Limit(100).Find(&Usuarios).Error
	if err != nil {
		return nil, err
	}
	return &Usuarios, err
}

/*
func (u *Usuario) Find(db *gorm.DB, uid uuid.UUID) (*Usuario, error) {

	err := db.Debug().Model(Usuario{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Usuario{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Usuario{}, errors.New("Usuario Inexistente")
	}
	return u, err
}
*/

func (u *Usuario) Find(db *gorm.DB, param string, uid string) (*Usuario, error) {
	err := db.Debug().Model(Usuario{}).Where(param, uid).Take(&u).Error
	if err != nil {
		return &Usuario{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Usuario{}, errors.New("Usuario Inexistente")
	}
	return u, nil
}

func (u *Usuario) Delete(db *gorm.DB, uid uuid.UUID) (int64, error) {
	db = db.Debug().Where("id = ?", uid).Delete(&Usuario{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (u *Usuario) DeleteBy(db *gorm.DB, cond string, uid interface{}) (int64, error) {
	result := db.Delete(&Usuario{}, cond+" = ?", uid)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (u *Usuario) Validate(action string) error {
	if u.Nome == "" || u.Nome == "null" {
		return errors.New("obrigatório: nome do usuário")
	}
	if u.Email == "" || u.Email == "null" {
		return errors.New("obrigatório: email")
	}
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
