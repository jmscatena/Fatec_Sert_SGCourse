package laboratorios

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jmscatena/Fatec_Sert_SGLab/dto/models/administrativo"
	"gorm.io/gorm"
	"html"
	"log"
	"strings"
	"time"
)

type Reservas struct {
	gorm.Model
	UID           uuid.UUID              `gorm:"type:uuid;default:uuid_generate_v4()" json:"ID"`
	LaboratorioID uint64                 `gorm:"default:0" json:"labid"`
	Laboratorio   Laboratorios           `gorm:"foreignKey: LaboratorioID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"laboratorio"`
	DataInicial   time.Time              `gorm:"type:DATE;default:CURRENT_TIMESTAMP" json:"data_inicio"`
	DataFinal     time.Time              `gorm:"type:DATE;default:CURRENT_TIMESTAMP" json:"data_fim"`
	HoraInicial   time.Time              `gorm:"type:TIME;default:CURRENT_TIMESTAMP" json:"hora_inicio"`
	HoraFinal     time.Time              `gorm:"type:TIME;default:CURRENT_TIMESTAMP" json:"hora_fim"`
	DiaSemana     string                 `gorm:"not null;" json:"dia_semana"`
	Rotativo      bool                   `gorm:"default:false" json:"rotativo"`
	Autorizado    bool                   `gorm:"default:false" json:"autorizado"`
	AutorizadoID  *uint64                `gorm:"default:null" json:"autorizadoid"`
	AutorizadoBy  administrativo.Usuario `gorm:"foreignKey: AutorizadoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"autorizado_por"`
	AutorizadoAt  time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"autorizado_em"`
	SolicitadoID  uint64                 `gorm:"default:0" json:"solicitadoid"`
	SolicitadoBy  administrativo.Usuario `gorm:"foreignKey: SolicitadoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"solicitado_por"`
	SolicitadoAt  time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"solicitado_em"`
	Ativa         bool                   `gorm:"default:false" json:"ativa"`
}

func (p *Reservas) Validate() error {

	if p.DiaSemana == "" || p.DiaSemana == "null" {
		return errors.New("obrigatório: dia da semana")
	}
	if p.LaboratorioID <= 0 {
		return errors.New("obrigatório: laboratório")
	}
	if p.SolicitadoID <= 0 {
		return errors.New("obrigatório: usuário para solicitação")
	}
	if p.DataInicial.IsZero() {
		return errors.New("obrigatório: data inicial")
	}
	if p.DataFinal.IsZero() {
		return errors.New("obrigatório: data final")
	}
	if p.HoraInicial.IsZero() {
		return errors.New("obrigatório: hora inicial")
	}
	if p.HoraFinal.IsZero() {
		return errors.New("obrigatório: hora final")
	}
	return nil
}
func (p *Reservas) Prepare(db *gorm.DB) (err error) {
	//usuario := administrativo.Usuario{}
	//laboratorio := Laboratorios{}

	p.DiaSemana = html.EscapeString(strings.TrimSpace(p.DiaSemana))
	//p.DataInicial, _ = time.Parse("2006-01-02", string(p.DataInicial))
	//p.DataFinal, _ = time.Parse("2006-01-02", p.DataFinal)
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	if p.AutorizadoID != nil {
		err = db.Model(&administrativo.Usuario{}).Where("id = ?", p.AutorizadoID).Take(p.AutorizadoBy).Error
		//p.AutorizadoBy = usuario
		p.AutorizadoAt = time.Now()
	} else {
		p.AutorizadoID = nil
	}

	err = db.Model(&administrativo.Usuario{}).Where("id = ?", p.SolicitadoID).Take(&p.SolicitadoBy).Error
	//p.SolicitadoBy = usuario
	p.SolicitadoAt = time.Now()

	err = db.Model(&Laboratorios{}).Preload("Materiais").Where("id = ?", p.LaboratorioID).Take(&p.Laboratorio).Error
	//p.Laboratorio = laboratorio

	if err != nil {
		log.Fatalf("Error during preparation:%v", err)
	}
	return
}

func (p *Reservas) Create(db *gorm.DB) (uuid.UUID, error) {
	if verr := p.Validate(); verr != nil {
		return uuid.Nil, verr
	}
	p.Prepare(db)
	err := db.Debug().Omit("ID").Create(&p).Error
	if err != nil {
		return uuid.Nil, err
	}
	return p.UID, nil
}

func (p *Reservas) Update(db *gorm.DB, uid uuid.UUID) (*Reservas, error) {
	db = db.Debug().Model(Reservas{}).Where("id = ?", uid).Updates(Reservas{
		Laboratorio:  p.Laboratorio,
		DataInicial:  p.DataInicial,
		DataFinal:    p.DataFinal,
		HoraInicial:  p.HoraInicial,
		HoraFinal:    p.HoraFinal,
		DiaSemana:    p.DiaSemana,
		Rotativo:     p.Rotativo,
		Autorizado:   p.Autorizado,
		AutorizadoBy: p.AutorizadoBy,
		SolicitadoBy: p.SolicitadoBy,
		Ativa:        p.Ativa})

	if db.Error != nil {
		return nil, db.Error
	}
	return p, nil
}

func (p *Reservas) List(db *gorm.DB) (*[]Reservas, error) {
	Reservass := []Reservas{}
	err := db.Debug().Model(&Reservas{}).
		Preload("Laboratorio").
		Preload("AutorizadoBy").
		Preload("SolicitadoBy").
		Limit(100).Find(&Reservass).Error
	//result := db.Find(&Reservass)
	if err != nil {
		return nil, err
	}
	return &Reservass, nil
}
func (u *Reservas) Find(db *gorm.DB, param string, uid string) (*Reservas, error) {
	err := db.Debug().Model(Reservas{}).Where(param, uid).Take(&u).Error
	if err != nil {
		return &Reservas{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Reservas{}, errors.New("Reserva Inexistente")
	}
	return u, nil
}

/*
	func (p *Reservas) Find(db *gorm.DB, uid uuid.UUID) (*Reservas, error) {
		err := db.Debug().Model(&Reservas{}).
			Preload("Laboratorio").
			Preload("AutorizadoBy").
			Preload("SolicitadoBy").
			Where("id = ?", uid).Take(&p).Error
		if err != nil {
			return &Reservas{}, err
		}
		return p, nil
	}

	func (p *Reservas) FindBy(db *gorm.DB, param string, uid ...interface{}) (*[]Reservas, error) {
		Reservass := []Reservas{}
		params := strings.Split(param, ";")
		uids := uid[0].([]interface{})
		if len(params) != len(uids) {
			return nil, errors.New("condição inválida")
		}
		result := db.
			Preload("Laboratorio").
			Preload("AutorizadoBy").
			Preload("SolicitadoBy").
			Where(strings.Join(params, " AND "), uids...).Find(&Reservass)
		if result.Error != nil {
			return nil, result.Error
		}
		return &Reservass, nil
	}
*/
func (p *Reservas) Delete(db *gorm.DB, uid uuid.UUID) (int64, error) {
	db = db.Delete(&Reservas{}, "id = ? ", uid)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (p *Reservas) DeleteBy(db *gorm.DB, cond string, uid uuid.UUID) (int64, error) {
	result := db.Delete(&Reservas{}, cond+" = ?", uid)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
