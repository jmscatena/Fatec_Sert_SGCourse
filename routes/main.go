package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uID"
	"github.com/jmscatena/Fatec_Sert_SGCourse/config"
	"github.com/jmscatena/Fatec_Sert_SGCourse/dto/models/administrativo"
	"github.com/jmscatena/Fatec_Sert_SGCourse/middleware"
	"github.com/jmscatena/Fatec_Sert_SGCourse/services"
)

func ConfigRoutes(router *gin.Engine, conn config.Connection, token config.SecretsToken) *gin.Engine {
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
				ID, _ := uID.Parse(context.Param("id"))
				condition := "Id=?"
				middleware.Get[administrativo.Usuario](context, &user, condition, ID.String(), conn)
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
				ID, _ := uID.Parse(context.Param("id"))
				middleware.Modify[administrativo.Usuario](context, &user, ID, conn)
			})
			userRoute.DELETE("/:id", func(context *gin.Context) {
				ID, _ := uID.Parse(context.Param("id"))
				middleware.Erase[administrativo.Usuario](context, &user, ID, conn)
			})

		}
	}
	/*
		matRoute := main.Group("materiais", services.Authenticate(conn, token))
		{
			var mat laboratorios2.Materiais
			matRoute.POST("/", func(context *gin.Context) {
				middleware.Add[laboratorios2.Materiais](context, &mat, conn)
			})
			matRoute.GET("/:id", func(context *gin.Context) {
				ID, _ := uID.Parse(context.Param("id"))
				condition := "Id=?"
				middleware.Get[laboratorios2.Materiais](context, &mat, condition, ID.String(), conn)
			})

			matRoute.GET("/", func(context *gin.Context) {
				middleware.GetAll[laboratorios2.Materiais](context, &mat, conn)
			})
			matRoute.PATCH("/:id", func(context *gin.Context) {
				ID, _ := uID.Parse(context.Param("id"))
				middleware.Modify[laboratorios2.Materiais](context, &mat, ID, conn)
			})
			matRoute.DELETE("/:id", func(context *gin.Context) {
				ID, _ := uID.Parse(context.Param("id"))
				middleware.Erase[laboratorios2.Materiais](context, &mat, ID, conn)
			})


				mat.GET("/admin/", func(context *gin.Context) {
					//colocar as configuracoes para os params q virao do frontend
					params := "admin=?;ativo=?"
					controllers.GetBy[administrativo.Usuario](context, &obj, params, false, true)
				})


		}
		lab := main.Group("laboratorios", services.Authenticate(conn, token))
		{
			var obj laboratorios2.Laboratorios
			lab.POST("/", func(context *gin.Context) {
				middleware.Add[laboratorios2.Laboratorios](context, &obj, conn)
			})
			lab.GET("/:id", func(context *gin.Context) {
				ID, _ := uID.Parse(context.Param("id"))
				condition := "Id=?"
				middleware.Get[laboratorios2.Laboratorios](context, &obj, condition, ID.String(), conn)
			})

			lab.GET("/", func(context *gin.Context) {
				middleware.GetAll[laboratorios2.Laboratorios](context, &obj, conn)
			})
			lab.PATCH("/:id", func(context *gin.Context) {
				ID, _ := uID.Parse(context.Param("id"))
				middleware.Modify[laboratorios2.Laboratorios](context, &obj, ID, conn)
			})
			lab.DELETE("/:id", func(context *gin.Context) {
				ID, _ := uID.Parse(context.Param("id"))
				middleware.Erase[laboratorios2.Laboratorios](context, &obj, ID, conn)
			})
		}
		res := main.Group("reservas", services.Authenticate(conn, token))
		{
			var obj laboratorios2.Reservas
			res.POST("/", func(context *gin.Context) {
				middleware.Add[laboratorios2.Reservas](context, &obj, conn)
			})
			res.GET("/:id", func(context *gin.Context) {
				ID, _ := uID.Parse(context.Param("id"))
				condition := "Id=?"
				middleware.Get[laboratorios2.Reservas](context, &obj, condition, ID.String(), conn)
			})

			res.GET("/", func(context *gin.Context) {
				middleware.GetAll[laboratorios2.Reservas](context, &obj, conn)
			})
			res.PATCH("/:id", func(context *gin.Context) {
				ID, _ := uID.Parse(context.Param("id"))
				middleware.Modify[laboratorios2.Reservas](context, &obj, ID, conn)
			})
			res.DELETE("/:id", func(context *gin.Context) {
				ID, _ := uID.Parse(context.Param("id"))
				middleware.Erase[laboratorios2.Reservas](context, &obj, ID, conn)
			})
		}
		ges := main.Group("gestao", services.Authenticate(conn, token))
		{
			var obj laboratorios2.GestaoMateriais
			ges.POST("/", func(context *gin.Context) {
				middleware.Add[laboratorios2.GestaoMateriais](context, &obj, conn)
			})
			ges.GET("/:id", func(context *gin.Context) {
				ID, _ := uID.Parse(context.Param("id"))
				condition := "Id=?"
				middleware.Get[laboratorios2.GestaoMateriais](context, &obj, condition, ID.String(), conn)
			})

			ges.GET("/", func(context *gin.Context) {
				middleware.GetAll[laboratorios2.GestaoMateriais](context, &obj, conn)
			})
			ges.PATCH("/:id", func(context *gin.Context) {
				ID, _ := uID.Parse(context.Param("id"))
				middleware.Modify[laboratorios2.GestaoMateriais](context, &obj, ID, conn)
			})
			ges.DELETE("/:id", func(context *gin.Context) {
				ID, _ := uID.Parse(context.Param("id"))
				middleware.Erase[laboratorios2.GestaoMateriais](context, &obj, ID, conn)
			})
		}

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
					ID, _ := uID.Parse(context.Param("id"))
					controllers.Get[models.Instituicao](context, &obj, uID)
				})
				inst.PATCH("/:id", func(context *gin.Context) {
					ID, _ := uID.Parse(context.Param("id"))
					controllers.Modify[models.Instituicao](context, &obj, uID)
				})
				inst.DELETE("/:id", func(context *gin.Context) {
					ID, _ := uID.Parse(context.Param("id"))
					controllers.Erase[models.Instituicao](context, &obj, uID)
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
					ID, _ := uID.Parse(context.Param("id"))
					controllers.Get[models.Evento](context, &obj, uID)
				})
				event.PATCH("/:id", func(context *gin.Context) {
					ID, _ := uID.Parse(context.Param("id"))
					controllers.Modify[models.Evento](context, &obj, uID)
				})
				event.DELETE("/:id", func(context *gin.Context) {
					ID, _ := uID.Parse(context.Param("id"))
					controllers.Erase[models.Evento](context, &obj, uID)
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
					ID, _ := uID.Parse(context.Param("id"))
					controllers.Get[models.Certificado](context, &obj, uID)
				})
				cert.PATCH("/:id", func(context *gin.Context) {
					ID, _ := uID.Parse(context.Param("id"))
					controllers.Modify[models.Certificado](context, &obj, uID)
				})
				cert.DELETE("/:id", func(context *gin.Context) {
					ID, _ := uID.Parse(context.Param("id"))
					controllers.Erase[models.Certificado](context, &obj, uID)
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
					ID, _ := uID.Parse(context.Param("id"))
					controllers.Get[models.CertVal](context, &obj, uID)
				})
				certval.PATCH("/", func(context *gin.Context) {
					ID, _ := uID.Parse(context.Param("id"))
					controllers.Modify[models.CertVal](context, &obj, uID)
				})
				certval.DELETE("/:id", func(context *gin.Context) {
					ID, _ := uID.Parse(context.Param("id"))
					controllers.Erase[models.CertVal](context, &obj, uID)
				})
			}
	*/
	//main.POST("login", controllers.Login)

	return router
}
