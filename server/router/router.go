package router

import (
	"shopify/handlers"
	"shopify/middleware"
	"shopify/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, sf *service.ServiceFactory, db *gorm.DB) {
	// 注入服务和数据库连接
	r.Use(middleware.InjectServices(sf, db))

	// 健康检查
	r.GET("/health", handlers.HealthCheck)

	// API v1 分组
	v1 := r.Group("/api/v1")
	{
		// 公开路由 - 不需要认证
		public := v1.Group("")
		{
			// 用户认证
			public.POST("/users/register", handlers.RegisterUser)
			public.POST("/users/login", handlers.LoginUser)
			public.POST("/users/password/reset/code", handlers.SendResetPasswordCode)
			public.POST("/users/password/reset", handlers.ResetPasswordByCode)

			// 公开的商品接口
			products := public.Group("/products")
			{
				products.GET("", handlers.ListProducts)                    // 获取商品列表
				products.GET("/:id", handlers.GetProduct)                  // 获取商品详情
				products.GET("/category/:category", handlers.ListProducts) // 按类别查询商品
				products.GET("/search", handlers.ListProducts)             // 搜索商品

				// 添加分类查询路由
				products.GET("/filter/category", handlers.ListProducts) // 按类别筛选
				products.GET("/filter/price", handlers.ListProducts)    // 按价格区间筛选
				products.GET("/filter/tags", handlers.ListProducts)     // 按标签筛选
				products.GET("/filter/keyword", handlers.ListProducts)  // 按关键词搜索

				// 商品评论
				products.GET("/:id/reviews", handlers.GetProductReviews) // 获取商品评论列表
			}

			// 公开的广告接口
			public.GET("/advertisements", handlers.ListAdvertisements)
			public.GET("/advertisements/position/:position", handlers.GetActiveAdvertisements)
			public.GET("/advertisements/:id", handlers.GetAdvertisement)

			// 支付相关路由
			payments := public.Group("/payments")
			{
				payments.POST("", handlers.CreatePayment)                // 创建支付
				payments.GET("/:id/status", handlers.QueryPaymentStatus) // 查询支付状态
			}
		}

		// 需要认证的路由
		authorized := v1.Group("")
		authorized.Use(middleware.AuthMiddleware())
		{
			// 用户相关
			users := authorized.Group("/users")
			{
				users.GET("/profile", handlers.GetUserProfile)
				users.PUT("/profile", handlers.UpdateUserProfile)

				users.POST("/addresses", handlers.AddUserAddress)
				users.GET("/addresses", handlers.ListUserAddresses)
				users.DELETE("/addresses/:id", handlers.DeleteUserAddress)
				users.PUT("/addresses/:id", handlers.UpdateUserAddress)
				users.PUT("/addresses/:id/default", handlers.SetDefaultAddresses)
			}

			// 订单相关
			orders := authorized.Group("/orders")
			{
				orders.POST("", handlers.CreateOrder)
				orders.GET("", handlers.ListOrders)
				orders.GET("/:id", handlers.GetOrder)
				orders.PUT("/:id/status", handlers.UpdateOrderStatus)
				orders.GET("/:id/logistics", handlers.GetLogistics)
			}

			// 购物车相关
			cart := authorized.Group("/cart")
			{
				cart.POST("/items", handlers.AddCartItem)
				cart.GET("/items", handlers.ListCartItems)
				cart.DELETE("/items/:id", handlers.RemoveCartItem)
				cart.PUT("/items/:id", handlers.UpdateCartItem)
				cart.PUT("/select-all", handlers.SelectAllCartItems)
				cart.PUT("/select-items", handlers.SelectCartItems) // 新增：选择部分商品
				cart.GET("/selected", handlers.GetSelectedCartItems)
			}

			// 评论相关
			reviews := authorized.Group("/reviews")
			{
				reviews.POST("", handlers.CreateReview)       // 创建评论
				reviews.GET("/me", handlers.GetUserReviews)   // 获取用户的评论列表
				reviews.DELETE("/:id", handlers.DeleteReview) // 删除评论
			}

			// 管理员路由组
			admin := authorized.Group("/admin")
			{
				// 用户管理
				admin.GET("/users", handlers.AdminListUsers)
				admin.PUT("/users/:id", handlers.AdminUpdateUser)
				admin.DELETE("/users/:id", handlers.AdminDeleteUser)

				// 商品管理
				adminProducts := admin.Group("/products")
				{
					adminProducts.POST("", handlers.CreateProduct)       // 创建商品
					adminProducts.PUT("/:id", handlers.UpdateProduct)    // 更新商品
					adminProducts.DELETE("/:id", handlers.DeleteProduct) // 删除商品
					adminProducts.GET("", handlers.ListProducts)         // 管理员查看所有商品
					adminProducts.GET("/:id", handlers.GetProduct)       // 管理员查看商品详情
				}

				// 订单管理
				adminOrders := admin.Group("/orders")
				{
					adminOrders.GET("", handlers.AdminListOrders)                        // 管理员查看所有订单
					adminOrders.GET("/:id", handlers.AdminGetOrder)                      // 管理员查看订单详情
					adminOrders.PUT("/:id/status", handlers.UpdateOrderStatus)           // 更新订单状态
					adminOrders.POST("/:id/logistics", handlers.UpdateLogistics)         // 更新物流信息
					adminOrders.POST("/:id/logistics/trace", handlers.AddLogisticsTrace) // 添加物流跟踪记录
				}

				// 商品管理
				advertisements := admin.Group("/advertisements")
				{
					// 广告管理
					advertisements.POST("", handlers.CreateAdvertisement)
					advertisements.GET("", handlers.ListAdvertisements)
					advertisements.GET("/:id", handlers.GetAdvertisement)
					advertisements.PUT("/:id", handlers.UpdateAdvertisement)
					advertisements.DELETE("/:id", handlers.DeleteAdvertisement)
					advertisements.PUT("/:id/status", handlers.UpdateAdvertisementStatus)
				}

				// 评论管理
				reviews := admin.Group("/reviews")
				{
					reviews.GET("/products/:id", handlers.AdminReviewsListOfProduct) // 管理员查看商品评论
					reviews.DELETE("/reviews/:id", handlers.AdminDeleteReview)       // 管理员删除评论
				}
			}
		}
	}
}
