package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmscatena/Fatec_Sert_SGCourse/dto/models/administrativo"
	laboratorios2 "github.com/jmscatena/Fatec_Sert_SGCourse/dto/models/laboratorios"
	"github.com/jmscatena/Fatec_Sert_SGCourse/infra"
	"github.com/jmscatena/Fatec_Sert_SGCourse/middleware"
	"github.com/jmscatena/Fatec_Sert_SGCourse/services"
)

func ConfigRoutes(router *gin.Engine, conn infra.Connection, token infra.SecretsToken) *gin.Engine {
	main := router.Group("/")
	{
		login := main.Group("login")
		{
			login.POST("/", services.Login(conn, token))
		}
		userRoute := main.Group("user", services.Authenticate(conn, token))
		{
			var user administrativo.Usuario
			userRoute.POST("/", func(context *gin.Context) {
				middleware.Add[administrativo.Usuario](context, &user, conn)
			})
			userRoute.GET("/:id", func(context *gin.Context) {
				uid, _ := uuid.Parse(context.Param("id"))
				condition := "Id=?"
				middleware.Get[administrativo.Usuario](context, &user, condition, uid.String(), conn)
			})

			userRoute.GET("/", func(context *gin.Context) {
				middleware.GetAll[administrativo.Usuario](context, &user, conn)
			})
			userRoute.GET("/admin/", func(context *gin.Context) {
				//colocar as configuracoes para os params q virao do frontend
				params := "admin=?;ativo=?"
				middleware.Get[administrativo.Usuario](context, &user, params, "false; true", conn)
			})

			userRoute.PATCH("/:id", func(context *gin.Context) {
				uid, _ := uuid.Parse(context.Param("id"))
				middleware.Modify[administrativo.Usuario](context, &user, uid, conn)
			})
			userRoute.DELETE("/:id", func(context *gin.Context) {
				uid, _ := uuid.Parse(context.Param("id"))
				middleware.Erase[administrativo.Usuario](context, &user, uid, conn)
			})

		}
	}

	matRoute := main.Group("materiais", services.Authenticate(conn, token))
	{
		var mat laboratorios2.Materiais
		matRoute.POST("/", func(context *gin.Context) {
			middleware.Add[laboratorios2.Materiais](context, &mat, conn)
		})
		matRoute.GET("/:id", func(context *gin.Context) {
			uid, _ := uuid.Parse(context.Param("id"))
			condition := "Id=?"
			middleware.Get[laboratorios2.Materiais](context, &mat, condition, uid.String(), conn)
		})

		matRoute.GET("/", func(context *gin.Context) {
			middleware.GetAll[laboratorios2.Materiais](context, &mat, conn)
		})
		matRoute.PATCH("/:id", func(context *gin.Context) {
			uid, _ := uuid.Parse(context.Param("id"))
			middleware.Modify[laboratorios2.Materiais](context, &mat, uid, conn)
		})
		matRoute.DELETE("/:id", func(context *gin.Context) {
			uid, _ := uuid.Parse(context.Param("id"))
			middleware.Erase[laboratorios2.Materiais](context, &mat, uid, conn)
		})

		/*
			mat.GET("/admin/", func(context *gin.Context) {
				//colocar as configuracoes para os params q virao do frontend
				params := "admin=?;ativo=?"
				controllers.GetBy[administrativo.Usuario](context, &obj, params, false, true)
			})
		*/

	}
	lab := main.Group("laboratorios", services.Authenticate(conn, token))
	{
		var obj laboratorios2.Laboratorios
		lab.POST("/", func(context *gin.Context) {
			middleware.Add[laboratorios2.Laboratorios](context, &obj, conn)
		})
		lab.GET("/:id", func(context *gin.Context) {
			uid, _ := uuid.Parse(context.Param("id"))
			condition := "Id=?"
			middleware.Get[laboratorios2.Laboratorios](context, &obj, condition, uid.String(), conn)
		})

		lab.GET("/", func(context *gin.Context) {
			middleware.GetAll[laboratorios2.Laboratorios](context, &obj, conn)
		})
		lab.PATCH("/:id", func(context *gin.Context) {
			uid, _ := uuid.Parse(context.Param("id"))
			middleware.Modify[laboratorios2.Laboratorios](context, &obj, uid, conn)
		})
		lab.DELETE("/:id", func(context *gin.Context) {
			uid, _ := uuid.Parse(context.Param("id"))
			middleware.Erase[laboratorios2.Laboratorios](context, &obj, uid, conn)
		})
	}
	res := main.Group("reservas", services.Authenticate(conn, token))
	{
		var obj laboratorios2.Reservas
		res.POST("/", func(context *gin.Context) {
			middleware.Add[laboratorios2.Reservas](context, &obj, conn)
		})
		res.GET("/:id", func(context *gin.Context) {
			uid, _ := uuid.Parse(context.Param("id"))
			condition := "Id=?"
			middleware.Get[laboratorios2.Reservas](context, &obj, condition, uid.String(), conn)
		})

		res.GET("/", func(context *gin.Context) {
			middleware.GetAll[laboratorios2.Reservas](context, &obj, conn)
		})
		res.PATCH("/:id", func(context *gin.Context) {
			uid, _ := uuid.Parse(context.Param("id"))
			middleware.Modify[laboratorios2.Reservas](context, &obj, uid, conn)
		})
		res.DELETE("/:id", func(context *gin.Context) {
			uid, _ := uuid.Parse(context.Param("id"))
			middleware.Erase[laboratorios2.Reservas](context, &obj, uid, conn)
		})
	}
	ges := main.Group("gestao", services.Authenticate(conn, token))
	{
		var obj laboratorios2.GestaoMateriais
		ges.POST("/", func(context *gin.Context) {
			middleware.Add[laboratorios2.GestaoMateriais](context, &obj, conn)
		})
		ges.GET("/:id", func(context *gin.Context) {
			uid, _ := uuid.Parse(context.Param("id"))
			condition := "Id=?"
			middleware.Get[laboratorios2.GestaoMateriais](context, &obj, condition, uid.String(), conn)
		})

		ges.GET("/", func(context *gin.Context) {
			middleware.GetAll[laboratorios2.GestaoMateriais](context, &obj, conn)
		})
		ges.PATCH("/:id", func(context *gin.Context) {
			uid, _ := uuid.Parse(context.Param("id"))
			middleware.Modify[laboratorios2.GestaoMateriais](context, &obj, uid, conn)
		})
		ges.DELETE("/:id", func(context *gin.Context) {
			uid, _ := uuid.Parse(context.Param("id"))
			middleware.Erase[laboratorios2.GestaoMateriais](context, &obj, uid, conn)
		})
	}
	/*
		inst := main.Group("instituicao")
		{
			var obj models.Instituicao
			inst.POST("/", func(context *gin.Context) {
				controllers.Add[models.Instituicao](context, &obj)
			})
			inst.GET("/", func(context *gin.Context) {
				controllers.GetAll[models.Instituicao](context, &obj)
			})
			inst.GET("/:id", func(context *gin.Context) {
				uid, _ := uuid.Parse(context.Param("id"))
				controllers.Get[models.Instituicao](context, &obj, uuid)
			})
			inst.PATCH("/:id", func(context *gin.Context) {
				uid, _ := uuid.Parse(context.Param("id"))
				controllers.Modify[models.Instituicao](context, &obj, uuid)
			})
			inst.DELETE("/:id", func(context *gin.Context) {
				uid, _ := uuid.Parse(context.Param("id"))
				controllers.Erase[models.Instituicao](context, &obj, uuid)
			})
		}
		event := main.Group("evento")
		{
			var obj models.Evento
			event.POST("/", func(context *gin.Context) {
				controllers.Add[models.Evento](context, &obj)
			})
			event.GET("/", func(context *gin.Context) {
				controllers.GetAll[models.Evento](context, &obj)
			})
			event.GET("/:id", func(context *gin.Context) {
				uid, _ := uuid.Parse(context.Param("id"))
				controllers.Get[models.Evento](context, &obj, uuid)
			})
			event.PATCH("/:id", func(context *gin.Context) {
				uid, _ := uuid.Parse(context.Param("id"))
				controllers.Modify[models.Evento](context, &obj, uuid)
			})
			event.DELETE("/:id", func(context *gin.Context) {
				uid, _ := uuid.Parse(context.Param("id"))
				controllers.Erase[models.Evento](context, &obj, uuid)
			})
		}
		cert := main.Group("cert")
		{
			var obj models.Certificado
			cert.POST("/", func(context *gin.Context) {
				controllers.Add[models.Certificado](context, &obj)
			})
			cert.GET("/", func(context *gin.Context) {
				controllers.GetAll[models.Certificado](context, &obj)
			})
			cert.GET("/:id", func(context *gin.Context) {
				uid, _ := uuid.Parse(context.Param("id"))
				controllers.Get[models.Certificado](context, &obj, uuid)
			})
			cert.PATCH("/:id", func(context *gin.Context) {
				uid, _ := uuid.Parse(context.Param("id"))
				controllers.Modify[models.Certificado](context, &obj, uuid)
			})
			cert.DELETE("/:id", func(context *gin.Context) {
				uid, _ := uuid.Parse(context.Param("id"))
				controllers.Erase[models.Certificado](context, &obj, uuid)
			})
		}
		certval := main.Group("valida")
		{
			var obj models.CertVal
			certval.POST("/", func(context *gin.Context) {
				controllers.Add[models.CertVal](context, &obj)
			})
			certval.GET("/", func(context *gin.Context) {
				controllers.GetAll[models.CertVal](context, &obj)
			})
			certval.GET("/:id", func(context *gin.Context) {
				uid, _ := uuid.Parse(context.Param("id"))
				controllers.Get[models.CertVal](context, &obj, uuid)
			})
			certval.PATCH("/", func(context *gin.Context) {
				uid, _ := uuid.Parse(context.Param("id"))
				controllers.Modify[models.CertVal](context, &obj, uuid)
			})
			certval.DELETE("/:id", func(context *gin.Context) {
				uid, _ := uuid.Parse(context.Param("id"))
				controllers.Erase[models.CertVal](context, &obj, uuid)
			})
		}
	*/
	//main.POST("login", controllers.Login)

	return router
}
