package services

import (
	"github.com/jmscatena/Fatec_Sert_SGCourse/config"
	"github.com/jmscatena/Fatec_Sert_SGCourse/handlers"
	"log"
)

func New[T handlers.Tables](o handlers.PersistenceHandler[T], conn config.Connection) (uint64, error) {
	/* metodo com devolucao do UID */
	//db, err := config.InitDB()
	if conn.Db == nil {
		log.Fatalln("No connection Database")
		return 0, nil // corrigir esse retorno
	}
	recid, err := o.Create(conn.Db)
	if err != nil {
		log.Fatalln(err)
		return 0, err
	}
	return recid, nil
}

func Update[T handlers.Tables](o handlers.PersistenceHandler[T], ID uint64, conn config.Connection) (*T, error) {
	//db, err := config.InitDB()
	if conn.Db == nil {
		log.Fatalln("No connection Database")
		return nil, nil
	}

	rec, err := o.Update(conn.Db, ID)
	if err != nil {
		//log.Fatalln(err)
		return nil, err
	}
	return rec, nil
}

func Del[T handlers.Tables](o handlers.PersistenceHandler[T], ID uint64, conn config.Connection) (int64, error) {
	//db, err := config.InitDB()
	if conn.Db == nil {
		return -1, nil
	}
	rec, err := o.Delete(conn.Db, ID)
	if err != nil {
		//log.Fatalln(err)
		return 0, err
	}
	return rec, nil
}

/*func Get[T handlers.Tables](o handlers.PersistenceHandler[T], ID uint64) (*T, error) {
	db, err := database.Init()
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	rec, err := o.Find(db, ID)
	if err != nil {
		//log.Fatalln(err)
		return nil, err
	}
	return rec, nil
}
*/

func GetAll[T handlers.Tables](o handlers.PersistenceHandler[T], conn config.Connection) (*[]T, error) {
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

func Get[T handlers.Tables](o handlers.PersistenceHandler[T], param string, values string, conn config.Connection) (*T, error) {
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
