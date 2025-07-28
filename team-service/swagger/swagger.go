package swagger

// This package is intended for Swagger documentation integration.
// Add your Swagger setup or generated code here.

import (
	_ "team-service/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupSwagger thiết lập route swagger
func SetupSwagger(r *gin.Engine) {
	url := ginSwagger.URL("http://localhost:8081/swagger/doc.json") // tùy chỉnh URL swagger json
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
