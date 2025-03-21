package services

import (
	"github.com/google/uuid"
	"github.com/jmscatena/Fatec_Sert_SGCourse/handlers"
	"github.com/jmscatena/Fatec_Sert_SGCourse/infra"
	"log"
)

func New[T handlers.Tables](o handlers.PersistenceHandler[T], conn infra.Connection) (uuid.UUID, error) {
	/* metodo com devolucao do UUID */
	//db, err := config.InitDB()
	if conn.Db == nil {
		log.Fatalln("No connection Database")
		return uuid.Nil, nil // corrigir esse retorno
	}
	recid, err := o.Create(conn.Db)
	if err != nil {
		log.Fatalln(err)
		return uuid.Nil, err
	}
	return recid, nil
}

func Update[T handlers.Tables](o handlers.PersistenceHandler[T], uid uuid.UUID, conn infra.Connection) (*T, error) {
	//db, err := config.InitDB()
	if conn.Db == nil {
		log.Fatalln("No connection Database")
		return nil, nil
	}

	rec, err := o.Update(conn.Db, uid)
	if err != nil {
		//log.Fatalln(err)
		return nil, err
	}
	return rec, nil
}

func Del[T handlers.Tables](o handlers.PersistenceHandler[T], uid uuid.UUID, conn infra.Connection) (int64, error) {
	//db, err := config.InitDB()
	if conn.Db == nil {
		return -1, nil
	}
	rec, err := o.Delete(conn.Db, uid)
	if err != nil {
		//log.Fatalln(err)
		return 0, err
	}
	return rec, nil
}

/*func Get[T handlers.Tables](o handlers.PersistenceHandler[T], uid uuid.UUID) (*T, error) {
	db, err := database.Init()
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	rec, err := o.Find(db, uid)
	if err != nil {
		//log.Fatalln(err)
		return nil, err
	}
	return rec, nil
}
*/

func GetAll[T handlers.Tables](o handlers.PersistenceHandler[T], conn infra.Connection) (*[]T, error) {
	//db, err := config.InitDB()
	if conn.Db == nil {
		return nil, nil
	}
	var rec *[]T
	rec, err := o.List(conn.Db)
	if err != nil {
		//log.Fatalln(err)
		return nil, err
	}
	return rec, nil
}
func Get[T handlers.Tables](o handlers.PersistenceHandler[T], param string, values string, conn infra.Connection) (*T, error) {
	//db, err := config.InitDB()
	if conn.Db == nil {
		return nil, nil
	}
	rec, err := o.Find(conn.Db, param, values)
	if err != nil {
		//log.Fatalln(err)
		return nil, err
	}
	return rec, nil
}
