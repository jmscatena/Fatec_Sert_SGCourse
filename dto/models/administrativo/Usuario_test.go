package administrativo

import (
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

func ShouldUsuarioCreateCorrect(t *testing.T) {
}

func TestUsuario_Create_Correct(t *testing.T) {
	type fields struct {
		Model     gorm.Model
		UID       uuid.UUID
		Nome      string
		Email     string
		Senha     string
		Ativo     bool
		Admin     bool
		Perfil    Perfil
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	type args struct {
		db *gorm.DB
	}
	err := godotenv.Load("../../../.env")
	if err != nil {

		log.Fatalf("Error Loading Configuration File")
	}
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbase := os.Getenv("DB")
	dbServer := os.Getenv("DBSERVER")
	dbPort := os.Getenv("DBPORT")
	dbURL := "postgres://" + dbUser + ":" + dbPass + "@" + dbServer + ":" + dbPort + "/" + dbase
	con, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "Teste de Usuário",
			fields:  fields{Nome: "Teste", Email: "Teste@email.com", Senha: "1234", Perfil: professor},
			args:    args{db: con},
			want:    0,
			wantErr: false,
		},
	}
	if err != nil {
		log.Fatalln(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usuario{
				UID:    tt.fields.UID,
				Nome:   tt.fields.Nome,
				Email:  tt.fields.Email,
				Senha:  tt.fields.Senha,
				Perfil: tt.fields.Perfil,
			}
			got, err := u.Create(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Error(err)
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsuario_Delete(t *testing.T) {
	type fields struct {
		Model     gorm.Model
		UID       uuid.UUID
		Nome      string
		Email     string
		Senha     string
		Ativo     bool
		Admin     bool
		Perfil    Perfil
		Tecnico   bool
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	type args struct {
		db   *gorm.DB
		uUID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usuario{
				Model:  tt.fields.Model,
				UID:    tt.fields.UID,
				Nome:   tt.fields.Nome,
				Email:  tt.fields.Email,
				Senha:  tt.fields.Senha,
				Ativo:  tt.fields.Ativo,
				Perfil: tt.fields.Perfil,
			}
			got, err := u.Delete(tt.args.db, tt.args.uUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsuario_DeleteBy(t *testing.T) {
	type fields struct {
		Model     gorm.Model
		UID       uuid.UUID
		Nome      string
		Email     string
		Senha     string
		Ativo     bool
		Admin     bool
		Perfil    Perfil
		Tecnico   bool
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	type args struct {
		db   *gorm.DB
		cond string
		uUID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usuario{
				Model:  tt.fields.Model,
				UID:    tt.fields.UID,
				Nome:   tt.fields.Nome,
				Email:  tt.fields.Email,
				Senha:  tt.fields.Senha,
				Ativo:  tt.fields.Ativo,
				Perfil: tt.fields.Perfil,
			}
			got, err := u.DeleteBy(tt.args.db, tt.args.cond, tt.args.uUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeleteBy() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsuario_Find(t *testing.T) {
	type fields struct {
		Model     gorm.Model
		ID        uuid.UUID
		Nome      string
		Email     string
		Senha     string
		Ativo     bool
		Admin     bool
		Perfil    Perfil
		Tecnico   bool
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	type args struct {
		db   *gorm.DB
		uUID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Usuario
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usuario{
				Model:  tt.fields.Model,
				UID:    tt.fields.UID,
				Nome:   tt.fields.Nome,
				Email:  tt.fields.Email,
				Senha:  tt.fields.Senha,
				Ativo:  tt.fields.Ativo,
				Perfil: tt.fields.Perfil,
			}
			got, err := u.Find(tt.args.db, tt.args.uUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsuario_FindBy(t *testing.T) {
	type fields struct {
		Model     gorm.Model
		ID        uuid.UUID
		Nome      string
		Email     string
		Senha     string
		Ativo     bool
		Admin     bool
		Perfil    Perfil
		Tecnico   bool
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	type args struct {
		db    *gorm.DB
		param string
		uUID  uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *[]Usuario
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usuario{
				Model:  tt.fields.Model,
				UID:    tt.fields.uUID,
				Nome:   tt.fields.Nome,
				Email:  tt.fields.Email,
				Senha:  tt.fields.Senha,
				Ativo:  tt.fields.Ativo,
				Perfil: tt.fields.Perfil,
			}
			got, err := u.FindBy(tt.args.db, tt.args.param, tt.args.uUID...)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindBy() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsuario_List(t *testing.T) {
	type fields struct {
		Model     gorm.Model
		ID        uuid.UUID
		Nome      string
		Email     string
		Senha     string
		Ativo     bool
		Admin     bool
		Perfil    Perfil
		Tecnico   bool
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *[]Usuario
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usuario{
				Model:  tt.fields.Model,
				UID:    tt.fields.UID,
				Nome:   tt.fields.Nome,
				Email:  tt.fields.Email,
				Senha:  tt.fields.Senha,
				Ativo:  tt.fields.Ativo,
				Perfil: tt.fields.Perfil,
			}
			got, err := u.List(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsuario_Prepare(t *testing.T) {
	type fields struct {
		Model     gorm.Model
		ID        uuid.UUID
		Nome      string
		Email     string
		Senha     string
		Ativo     bool
		Admin     bool
		Perfil    Perfil
		Tecnico   bool
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usuario{
				Model:  tt.fields.Model,
				UID:    tt.fields.UID,
				Nome:   tt.fields.Nome,
				Email:  tt.fields.Email,
				Senha:  tt.fields.Senha,
				Ativo:  tt.fields.Ativo,
				Perfil: tt.fields.Perfil,
			}
			u.Prepare()
		})
	}
}

func TestUsuario_Update(t *testing.T) {
	type fields struct {
		Model     gorm.Model
		UID       uuid.UUID
		Nome      string
		Email     string
		Senha     string
		Ativo     bool
		Perfil    Perfil
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	type args struct {
		db   *gorm.DB
		uUID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Usuario
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usuario{
				Model:  tt.fields.Model,
				UID:    tt.fields.UID,
				Nome:   tt.fields.Nome,
				Email:  tt.fields.Email,
				Senha:  tt.fields.Senha,
				Ativo:  tt.fields.Ativo,
				Perfil: tt.fields.Perfil,
			}
			got, err := u.Update(tt.args.db, tt.args.uUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsuario_ValUIDate(t *testing.T) {
	type fields struct {
		Model     gorm.Model
		UID       uuid.UUID
		Nome      string
		Email     string
		Senha     string
		Ativo     bool
		Admin     bool
		Perfil    Perfil
		Tecnico   bool
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	type args struct {
		action string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Usuario{
				Model:  tt.fields.Model,
				UID:    tt.fields.UID,
				Nome:   tt.fields.Nome,
				Email:  tt.fields.Email,
				Senha:  tt.fields.Senha,
				Ativo:  tt.fields.Ativo,
				Perfil: tt.fields.Perfil,
			}
			if err := u.ValUIDate(tt.args.action); (err != nil) != tt.wantErr {
				t.Errorf("ValUIDate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVerifySenha(t *testing.T) {
	type args struct {
		hashedSenha string
		Senha       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifySenha(tt.args.hashedSenha, tt.args.Senha); (err != nil) != tt.wantErr {
				t.Errorf("VerifySenha() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
