package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmscatena/Fatec_Sert_SGCourse/config"
	"github.com/jmscatena/Fatec_Sert_SGCourse/handlers"
	"github.com/jmscatena/Fatec_Sert_SGCourse/services"
	"net/http"
	"reflect"
)

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "acesso " + http.StatusText(200)})

}

func Add[T handlers.Tables](c *gin.Context, o handlers.PersistenceHandler[T], conn config.Connection) {
	if reflect.TypeOf(o) != nil {
		if err := c.ShouldBindJSON(&o); err != nil {
			fmt.Println("ERRO:", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusText(http.StatusBadRequest), "data": err})
			return
		}
		var handler handlers.PersistenceHandler[T] = o
		code, cerr := services.New(handler, conn)

		if cerr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusText(http.StatusConflict), "data": cerr})
		}
		c.JSON(http.StatusCreated, gin.H{"status": http.StatusText(http.StatusCreated), "data": code})
	}
}

func Modify[T handlers.Tables](c *gin.Context, o handlers.PersistenceHandler[T], ID uint, conn config.Connection) {
	if reflect.TypeOf(o) != nil {
		if err := c.ShouldBindJSON(&o); err != nil {
			fmt.Println("ERRO:", err)
			//msg para deploy
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusText(http.StatusBadRequest), "data": "Erro de JSON"})
			//msg para dev
			//c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusText(http.StatusBadRequest), "data": err})
			return

		}
		var handler handlers.PersistenceHandler[T] = o
		code, cerr := services.Update(handler, ID, conn)

		if cerr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusText(http.StatusBadRequest), "data": cerr})
		}
		c.JSON(http.StatusCreated, gin.H{"status": http.StatusText(http.StatusAccepted), "data": code})
	}
}

func Erase[T handlers.Tables](c *gin.Context, o handlers.PersistenceHandler[T], ID uint, conn config.Connection) {
	if reflect.TypeOf(o) != nil {
		var handler handlers.PersistenceHandler[T] = o
		code, cerr := services.Del(handler, ID, conn)

		if cerr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusText(http.StatusBadRequest), "data": cerr})
		}
		c.JSON(http.StatusCreated, gin.H{"status": http.StatusText(http.StatusNoContent), "data": code})
	}
}

/*
	func Get[T handlers.Tables](c *gin.Context, o handlers.PersistenceHandler[T], ID uint) {
		if reflect.TypeOf(o) != nil {
			if ID == uID.Nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusText(http.StatusBadRequest), "data": "{}"})
				return
			}
			var handler handlers.PersistenceHandler[T] = o
			rec, cerr := services.Get(handler, ID)

			if cerr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusText(http.StatusBadRequest), "data": cerr})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": rec, "status": http.StatusText(http.StatusOK)})
		}
	}
*/
func GetAll[T handlers.Tables](c *gin.Context, o handlers.PersistenceHandler[T], conn config.Connection) {
	if reflect.TypeOf(o) != nil {
		var handler handlers.PersistenceHandler[T] = o
		rec, cerr := services.GetAll(handler, conn)
		if cerr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusText(http.StatusBadRequest),
				"data": cerr})
		}
		c.JSON(http.StatusOK, gin.H{"data": rec, "status": http.StatusText(http.StatusOK)})
	}
}

func Get[T handlers.Tables](c *gin.Context, o handlers.PersistenceHandler[T], param map[string]interface{}, conn config.Connection) {
	if reflect.TypeOf(o) != nil {
		if len(param) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusText(http.StatusNotFound), "data": "No Data"})
		}
		var handler handlers.PersistenceHandler[T] = o
		rec, cerr := services.Get(handler, param, conn)

		if cerr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusText(http.StatusBadRequest), "data": cerr})
		}
		c.JSON(http.StatusOK, gin.H{"data": rec, "status": http.StatusText(http.StatusOK)})
	}
}
