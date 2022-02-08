package factory

import (
	"github.com/awesomeProject/controller"
	"github.com/awesomeProject/dbconfig"
	"github.com/awesomeProject/service"
)

var (
	ProductService        service.ProductService
	ProductInsightService service.ProductInsightService
)

func Init() {

	sqlDb := dbconfig.CreatePostgresConn()

	productServiceControllerImpl := controller.ProductServiceControllerImpl{
		SqlDb:   sqlDb,
		RedisDb: dbconfig.CreateRedisConn(),
	}

	productInsightControllerImpl := controller.ProductInsightControllerImpl{SqlDb: sqlDb}

	ProductInsightService = service.ProductInsightServiceImpl{ProductInsightController: productInsightControllerImpl}

	ProductService = service.ProductServiceImpl{ProductController: productServiceControllerImpl}

}
