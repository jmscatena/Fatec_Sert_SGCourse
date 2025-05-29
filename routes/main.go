package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmscatena/Fatec_Sert_SGCourse/config"
	"github.com/jmscatena/Fatec_Sert_SGCourse/dto/models/administrativo"
	curso "github.com/jmscatena/Fatec_Sert_SGCourse/dto/models/cursos"
	"github.com/jmscatena/Fatec_Sert_SGCourse/middleware"
	"github.com/jmscatena/Fatec_Sert_SGCourse/services"
	"log"
	"net/http"
	"strconv"
)

func formatBytes(bytes int64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
	)

	if bytes < kb {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < mb {
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(kb))
	} else if bytes < gb {
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(mb))
	} else {
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(gb))
	}
}
func ConfigRoutes(router *gin.Engine, conn config.Connection, token config.SecretsToken) *gin.Engine {
	main := router.Group("/")
	{
		main.GET("/", services.Authenticate(conn, token), func(context *gin.Context) {
			var openSolicitations curso.Solicitacao_Doc
			var closedSolicitations curso.Solicitacao_Doc
			var profUsers administrativo.Usuario

			middleware.FindAll[curso.Solicitacao_Doc](context,
				&closedSolicitations,
				map[string]interface{}{"entrega": true, "ativo": true},
				conn)
			middleware.FindAll[curso.Solicitacao_Doc](context,
				&openSolicitations,
				map[string]interface{}{"entrega": false, "ativo": true},
				conn)

			middleware.FindAll[administrativo.Usuario](context,
				&profUsers,
				map[string]interface{}{"professor": true, "ativo": true},
				conn)

			context.JSON(200, gin.H{
				"openSolicitations":   openSolicitations,
				"closedSolicitations": closedSolicitations,
				"professors":          profUsers,
			})
		})

		signupstatus := main.Group("status")
		{
			signupstatus.POST("/", services.SignupStatus(conn))
		}

		login := main.Group("login")
		{
			login.POST("/", services.Login(conn, token))
		}
		signup := main.Group("signup")
		{
			signup.POST("/", services.Signup(conn, token))
		}
		logoutRoute := main.Group("logout", services.Authenticate(conn, token))
		{
			logoutRoute.POST("/", services.Logout(conn, token))
		}
		userRoute := main.Group("user", services.Authenticate(conn, token))
		{
			var user administrativo.Usuario
			userRoute.POST("/", func(context *gin.Context) {
				middleware.Add[administrativo.Usuario](context, &user, conn)
			})
			userRoute.GET("/:id", func(context *gin.Context) {
				ID := context.Param("id")
				params := map[string]interface{}{"id": ID, "ativo": true}
				middleware.Get[administrativo.Usuario](context, &user, params, conn)
			})

			userRoute.GET("/admin/", func(context *gin.Context) {
				//colocar as configuracoes para os params q virao do frontend
				params := map[string]interface{}{"diretor": true, "ativo": true}
				middleware.Get[administrativo.Usuario](context, &user, params, conn)
			})
			userRoute.GET("/professors/", func(context *gin.Context) {
				//colocar as configuracoes para os params q virao do frontend
				params := map[string]interface{}{"professor": true, "ativo": true}
				middleware.FindAll[administrativo.Usuario](context, &user, params, conn)
			})

			userRoute.PATCH("/:id", func(context *gin.Context) {
				ID, _ := strconv.ParseUint(context.Param("id"), 10, 64)
				middleware.Modify[administrativo.Usuario](context, &user, uint(ID), conn)
			})
			userRoute.DELETE("/:id", func(context *gin.Context) {
				ID, _ := strconv.ParseUint(context.Param("id"), 10, 64)
				middleware.Erase[administrativo.Usuario](context, &user, uint(ID), conn)
			})

		}
		usersRoute := main.Group("users", services.Authenticate(conn, token))
		{
			var usuario administrativo.Usuario
			usersRoute.GET("/", func(context *gin.Context) {
				middleware.GetAll[administrativo.Usuario](context, &usuario, conn)
			})

		}
		courseRoute := main.Group("course", services.Authenticate(conn, token))
		{
			var course curso.Curso
			courseRoute.POST("/", func(context *gin.Context) {
				middleware.Add[curso.Curso](context, &course, conn)
			})
			courseRoute.GET("/:id", func(context *gin.Context) {
				ID := context.Param("id")
				params := map[string]interface{}{"id": ID, "ativo": true}
				middleware.Get[curso.Curso](context, &course, params, conn)
			})

			courseRoute.GET("/", func(context *gin.Context) {
				middleware.GetAll[curso.Curso](context, &course, conn)
			})

			courseRoute.PATCH("/:id", func(context *gin.Context) {
				ID, _ := strconv.ParseUint(context.Param("id"), 10, 64)
				middleware.Modify[curso.Curso](context, &course, uint(ID), conn)
			})
			courseRoute.DELETE("/:id", func(context *gin.Context) {
				ID, _ := strconv.ParseUint(context.Param("id"), 10, 64)
				middleware.Erase[curso.Curso](context, &course, uint(ID), conn)
			})

		}
		coursesRoute := main.Group("courses", services.Authenticate(conn, token))
		{
			var course curso.Curso
			coursesRoute.GET("/", func(context *gin.Context) {
				middleware.GetAll[curso.Curso](context, &course, conn)
			})
		}
		disciplineRoute := main.Group("discipline", services.Authenticate(conn, token))
		{
			var discipline curso.Disciplina
			disciplineRoute.POST("/", func(context *gin.Context) {
				middleware.Add[curso.Disciplina](context, &discipline, conn)
			})
			disciplineRoute.GET("/:id", func(context *gin.Context) {
				ID := context.Param("id")
				params := map[string]interface{}{"id": ID, "ativo": true}
				middleware.Get[curso.Disciplina](context, &discipline, params, conn)
			})
			disciplineRoute.PATCH("/:id", func(context *gin.Context) {
				ID, _ := strconv.ParseUint(context.Param("id"), 10, 64)
				middleware.Modify[curso.Disciplina](context, &discipline, uint(ID), conn)
			})
			disciplineRoute.DELETE("/:id", func(context *gin.Context) {
				ID, _ := strconv.ParseUint(context.Param("id"), 10, 64)
				middleware.Erase[curso.Disciplina](context, &discipline, uint(ID), conn)
			})

		}
		disciplinesRoute := main.Group("disciplines", services.Authenticate(conn, token))
		{
			var disciplines curso.Disciplina
			disciplinesRoute.GET("/", func(context *gin.Context) {
				middleware.GetAll[curso.Disciplina](context, &disciplines, conn)
			})
		}
		docRoute := main.Group("document", services.Authenticate(conn, token))
		{
			var document curso.Documento
			docRoute.POST("/", func(context *gin.Context) {
				middleware.Add[curso.Documento](context, &document, conn)
			})
			docRoute.GET("/:id", func(context *gin.Context) {
				ID := context.Param("id")
				params := map[string]interface{}{"id": ID, "ativo": true}
				middleware.Get[curso.Documento](context, &document, params, conn)
			})
			docRoute.PATCH("/:id", func(context *gin.Context) {
				ID, _ := strconv.ParseUint(context.Param("id"), 10, 64)
				middleware.Modify[curso.Documento](context, &document, uint(ID), conn)
			})
			docRoute.DELETE("/:id", func(context *gin.Context) {
				ID, _ := strconv.ParseUint(context.Param("id"), 10, 64)
				middleware.Erase[curso.Documento](context, &document, uint(ID), conn)
			})

		}
		docsRoute := main.Group("documents", services.Authenticate(conn, token))
		{
			var docs curso.Documento
			docsRoute.GET("/", func(context *gin.Context) {
				middleware.GetAll[curso.Documento](context, &docs, conn)
			})
		}
		requisitionRoute := main.Group("requisition", services.Authenticate(conn, token))
		{
			var requisition curso.Solicitacao_Doc
			requisitionRoute.POST("/", func(context *gin.Context) {
				middleware.Add[curso.Solicitacao_Doc](context, &requisition, conn)
			})
			requisitionRoute.GET("/:id", func(context *gin.Context) {
				ID := context.Param("id")
				params := map[string]interface{}{"id": ID, "ativo": true}
				middleware.Get[curso.Solicitacao_Doc](context, &requisition, params, conn)
			})
			requisitionRoute.GET("/professor/", func(context *gin.Context) {
				params := map[string]interface{}{"email": context.Request.Header.Get("Id")}
				middleware.FindAll[curso.Solicitacao_Doc](context, &requisition, params, conn)
			})
			requisitionRoute.PATCH("/:id", func(context *gin.Context) {
				ID, _ := strconv.ParseUint(context.Param("id"), 10, 64)
				middleware.Modify[curso.Solicitacao_Doc](context, &requisition, uint(ID), conn)
			})
			requisitionRoute.DELETE("/:id", func(context *gin.Context) {
				ID, _ := strconv.ParseUint(context.Param("id"), 10, 64)
				middleware.Erase[curso.Solicitacao_Doc](context, &requisition, uint(ID), conn)
			})

		}
		requisitionsRoute := main.Group("requisitions", services.Authenticate(conn, token))
		{
			var reqs curso.Solicitacao_Doc
			requisitionsRoute.GET("/", func(context *gin.Context) {
				middleware.GetAll[curso.Solicitacao_Doc](context, &reqs, conn)
			})
		}
		deliveryRoute := main.Group("delivery", services.Authenticate(conn, token))
		{
			const MAX_UPLOAD_SIZE = 100 * 1024 * 1024 // 100 MB
			var delivery curso.Entrega_Doc
			deliveryRoute.POST("/", func(context *gin.Context) {
				reqCtx := context.Request.Context()
				fileHeader, err := context.FormFile("file") // "file" is the name of the input field in the form
				if err != nil {
					log.Printf("Request Context %p: Error getting file from form: %v", reqCtx, err)
					context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from form", "details": err.Error()})
					return
				}
				fileSize := fileHeader.Size // O tamanho já está disponível no FileHeader
				if fileSize > MAX_UPLOAD_SIZE {
					log.Printf("Request Context %p: File size (%d bytes) exceeds limit of %d bytes", reqCtx, fileSize, MAX_UPLOAD_SIZE)
					context.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": fmt.Sprintf("File size exceeds the allowed limit of %s", formatBytes(MAX_UPLOAD_SIZE))})
					return
				}
				solicitacaoID := context.PostForm("solicitacaoID")
				if solicitacaoID == "" {
					log.Printf("Request Context %p: Error in Solicitation: %v", reqCtx, err)
					context.JSON(http.StatusInternalServerError, gin.H{"error": "Solicitation is not present", "details": err.Error()})
					return
				}
				var solicitacao curso.Solicitacao_Doc
				params := map[string]interface{}{"id": solicitacaoID, "ativo": true}
				middleware.Get[curso.Solicitacao_Doc](context, &solicitacao, params, conn)
				fileName := strconv.FormatUint(uint64(solicitacao.ID), 10) + " - " + solicitacao.Documento.Titulo

				filePathName, err := delivery.SaveFile(fileHeader, fileName)
				if err != nil {
					log.Printf("Request Context %p: Error saving file: %v", reqCtx, err)
					context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file", "details": err.Error()})
					return
				}
				delivery.Arquivo = filePathName
				delivery.SolicitacaoID = solicitacao.ID
				_, err = delivery.Create(conn.Db)
				if err != nil {
					log.Printf("Request Context %p: Error creating delivery: %v", reqCtx, err)
					context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create delivery", "details": err.Error()})
					return
				}
			})

			deliveryRoute.GET("/:id", func(context *gin.Context) {
				ID := context.Param("id")
				params := map[string]interface{}{"id": ID}
				middleware.Get[curso.Entrega_Doc](context, &delivery, params, conn)
			})
			deliveryRoute.GET("/professor/", func(context *gin.Context) {
				params := map[string]interface{}{"email": context.Request.Header.Get("Id")}
				middleware.FindAll[curso.Entrega_Doc](context, &delivery, params, conn)
			})
			deliveryRoute.PATCH("/:id", func(context *gin.Context) {
				ID, _ := strconv.ParseUint(context.Param("id"), 10, 64)
				middleware.Modify[curso.Entrega_Doc](context, &delivery, uint(ID), conn)
			})
			deliveryRoute.DELETE("/:id", func(context *gin.Context) {
				ID, _ := strconv.ParseUint(context.Param("id"), 10, 64)
				middleware.Erase[curso.Entrega_Doc](context, &delivery, uint(ID), conn)
			})

		}
		deliveriesRoute := main.Group("deliveries", services.Authenticate(conn, token))
		{
			var deliveries curso.Entrega_Doc
			deliveriesRoute.GET("/", func(context *gin.Context) {
				middleware.GetAll[curso.Entrega_Doc](context, &deliveries, conn)
			})
		}

	}
	return router
}
