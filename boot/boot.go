package boot

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Start(port, contextPath string, router func(webRouter *WebRouter)) {
	gin.SetMode(gin.ReleaseMode)
	ginServer := gin.New()
	ginServer.Use(gin.Logger())
	ginServer.Use(Recovery)
	ginServer.NoRoute(_404Handler)
	ginServer.NoMethod(_404Handler)
	apiRouter := ginServer.Group(contextPath)

	webRouter := WebRouter{Router: apiRouter}
	if router != nil {
		router(&webRouter)
	}

	log.Println("Web Server starting http://127.0.0.1:" + port + contextPath)
	http.ListenAndServe(":"+port, ginServer)
}
