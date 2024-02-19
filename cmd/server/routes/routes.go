package routes

import (
	"database/sql"

	inboudorder "github.com/davidop97/apiGo/internal/inboudOrder"

	"github.com/davidop97/apiGo/cmd/server/handler"

	"github.com/davidop97/apiGo/internal/batch"

	"github.com/davidop97/apiGo/internal/locality"

	"github.com/davidop97/apiGo/internal/buyer"
	"github.com/davidop97/apiGo/internal/carries"
	"github.com/davidop97/apiGo/internal/employee"
	"github.com/davidop97/apiGo/internal/product"
	"github.com/davidop97/apiGo/internal/purchase_order"
	"github.com/davidop97/apiGo/internal/section"

	"github.com/davidop97/apiGo/internal/seller"
	"github.com/davidop97/apiGo/internal/warehouse"
	"github.com/gin-gonic/gin"

	//import docs for swagger
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	//import docs for swagger
	_ "github.com/davidop97/apiGo/docs"
)

type Router interface {
	MapRoutes()
}

type router struct {
	eng *gin.Engine
	rg  *gin.RouterGroup
	db  *sql.DB
}

func NewRouter(eng *gin.Engine, db *sql.DB) Router {
	return &router{eng: eng, db: db}
}

func (r *router) MapRoutes() {
	r.setGroup()

	r.buildSellerRoutes()
	r.buildlocalityRoutes()
	r.buildProductRoutes()
	r.buildSectionRoutes()
	r.buildWarehouseRoutes()
	r.buildEmployeeRoutes()
	r.buildBuyerRoutes()
	r.buildCarriesRoutes()
	r.buildInboudOrderRoutes()
	r.buildBatchRoutes()
	r.buildPORoutes()
}

func (r *router) setGroup() {
	r.rg = r.eng.Group("/api/v1")
}

func (r *router) buildSellerRoutes() {
	// Example
	repo := seller.NewRepository(r.db)
	service := seller.NewService(repo)
	handler := handler.NewSeller(service)
	r.rg.GET("/seller", handler.GetAll())
	r.rg.GET("/seller/:id", handler.Get())
	r.rg.DELETE("/seller/:id", handler.Delete())
	r.rg.POST("/seller", handler.Create())
	r.rg.PATCH("/seller/:id", handler.Update())

}

func (r *router) buildlocalityRoutes() {
	repo := locality.NewRepository(r.db)
	service := locality.NewService(repo)
	handler := handler.NewLocality(service)
	r.rg.GET("/localities/:id", handler.GetLocalityById())
	r.rg.GET("/localities/", handler.GetAll())
	r.rg.POST("/localities/", handler.Create())
	r.rg.GET("/localities/reportSellers", handler.GetReportSellers())
}

func (r *router) buildProductRoutes() {
	repo := product.NewRepository(r.db)
	service := product.NewService(repo)
	handler := handler.NewProduct(service)
	prodGroup := r.rg.Group("/products")
	prodGroup.GET("/", handler.GetAll())
	prodGroup.GET("/:id", handler.Get())
	prodGroup.POST("/", handler.Create())
	prodGroup.PATCH("/:id", handler.Update())
	prodGroup.DELETE("/:id", handler.Delete())
	r.rg.GET("/ping", handler.Ping())

	//routes for productRecords
	r.rg.POST("/productRecords", handler.CreateProductRecord())
	prodGroup.GET("/reportRecords", handler.GetProductRecord())

	r.rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (r *router) buildSectionRoutes() {
	repo := section.NewRepository(r.db)
	service := section.NewService(repo)
	handler := handler.NewSection(service)
	sectGroup := r.rg.Group("/sections")
	sectGroup.GET("/", handler.GetAll())
	sectGroup.GET("/:id", handler.Get())
	sectGroup.POST("/", handler.Create())
	sectGroup.DELETE("/:id", handler.Delete())
	sectGroup.PATCH("/:id", handler.Update())
	sectGroup.GET("/reportProducts", handler.ProductCount())
}

func (r *router) buildWarehouseRoutes() {
	repo := warehouse.NewRepository(r.db)
	service := warehouse.NewService(repo)
	warehouseHandler := handler.NewWarehouse(service)
	warehouseRouter := r.rg.Group("/warehouses")
	warehouseRouter.GET("/", warehouseHandler.GetAll())
	warehouseRouter.GET("/:id", warehouseHandler.Get())
	warehouseRouter.POST("/", warehouseHandler.Create())
	warehouseRouter.PATCH("/:id", warehouseHandler.Update())
	warehouseRouter.DELETE("/:id", warehouseHandler.Delete())
}

func (r *router) buildEmployeeRoutes() {
	repo := employee.NewRepository(r.db)
	service := employee.NewService(repo)
	handler := handler.NewEmployee(service)
	r.rg.GET("/employees", handler.GetAll())
	r.rg.GET("/employees/:id", handler.Get())
	r.rg.POST("/employees", handler.Create())
	r.rg.PATCH("/employees/:id", handler.Update())
	r.rg.DELETE("/employees/:id", handler.Delete())
}

func (r *router) buildBuyerRoutes() {
	repo := buyer.NewRepository(r.db)
	service := buyer.NewService(repo)
	handler := handler.NewBuyer(service)
	//r.rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.rg.GET("/buyers", handler.GetAll())
	r.rg.GET("/buyers/:id", handler.Get())
	r.rg.DELETE("/buyers/:id", handler.Delete())
	r.rg.POST("/buyers", handler.Create())
	r.rg.PATCH("/buyers/:id", handler.Update())
}

func (r *router) buildCarriesRoutes() {
	repo := carries.NewRepository(r.db)
	service := carries.NewService(repo)
	handler := handler.NewCarry(service)
	//r.rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.rg.GET("/carries", handler.GetAll())
	r.rg.POST("/carries", handler.Save())
	r.rg.GET("/localities/reportCarries", handler.GetCarriesByLocality())
}
func (r *router) buildInboudOrderRoutes() {
	repo := inboudorder.NewRepository(r.db)
	service := inboudorder.NewService(repo)
	handler := handler.NewInboudOrder(service)
	r.rg.GET("/employees/reportInboundOrders", handler.GenerateReport())
	r.rg.GET("/employees/reportInboundOrder", handler.GetAllReports())
	r.rg.POST("/inboundOrders", handler.CreateInboundOrder())
}
func (r *router) buildBatchRoutes() {
	repo := batch.NewRepository(r.db)
	service := batch.NewService(repo)
	handler := handler.NewProductBatch(service)
	batchGroup := r.rg.Group("/productBatches")
	batchGroup.GET("/", handler.GetAll())
	batchGroup.POST("/", handler.Create())
}

// purchase order route
func (r *router) buildPORoutes() {
	repo := purchase_order.NewRepository(r.db)
	service := purchase_order.NewService(repo)
	handler := handler.NewPurchaseOrder(service)
	r.rg.POST("/purchaseOrders", handler.Create())
	r.rg.GET("/buyers/reportPurchaseOrders", handler.ReportPurchaseOrdersByBuyer())
}
