package router

import (
	"adcms-backend/internal/config"
	"adcms-backend/internal/handler"
	"adcms-backend/internal/middleware"
	"adcms-backend/internal/repository"
	"adcms-backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func withPermission(code string, h gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware.PermissionMiddleware(code)(c)
		if c.IsAborted() {
			return
		}
		h(c)
	}
}

type CRUDHandlers struct {
	List        gin.HandlerFunc
	Create      gin.HandlerFunc
	Update      gin.HandlerFunc
	Delete      gin.HandlerFunc
	Status      gin.HandlerFunc
	BatchDelete gin.HandlerFunc
}

func registerCRUDRoutes(group *gin.RouterGroup, handlers CRUDHandlers) {
	if handlers.List != nil {
		group.POST("/list", handlers.List)
	}
	if handlers.Create != nil {
		group.POST("/create", handlers.Create)
	}
	if handlers.Update != nil {
		group.POST("/update", handlers.Update)
	}
	if handlers.Delete != nil {
		group.POST("/delete", handlers.Delete)
	}
	if handlers.Status != nil {
		group.POST("/status", handlers.Status)
	}
	if handlers.BatchDelete != nil {
		group.POST("/batch-delete", handlers.BatchDelete)
	}
}

// Setup 设置路由
func Setup(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// 初始化 Repository
	userRepo := repository.NewUserRepository(db)
	infoRepo := repository.NewInfoRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	infoCategoryRepo := repository.NewInfoCategoryRepository(db)
	siteGroupRepo := repository.NewSiteGroupRepository(db)
	formRepo := repository.NewFormRepository(db)
	systemLogRepo := repository.NewSystemLogRepository(db)

	// 初始化 Service
	authService := service.NewAuthService(userRepo, cfg)
	middleware.SetPermissionChecker(func(c *gin.Context, code string) bool {
		if isSuper, ok := c.Get("is_platform_super_admin"); ok {
			if v, ok2 := isSuper.(bool); ok2 && v {
				return true
			}
		}
		userIDVal, ok := c.Get("user_id")
		if !ok {
			return false
		}
		roleVal, ok := c.Get("role")
		if !ok {
			return false
		}
		userID, ok := userIDVal.(string)
		if !ok || userID == "" {
			return false
		}
		role, _ := roleVal.(string)
		codes, err := userRepo.GetAccessCodesByUser(userID, role)
		if err != nil {
			return false
		}
		for _, c := range codes {
			if c == code {
				return true
			}
		}
		return false
	})

	// 初始化 Handler
	authHandler := handler.NewAuthHandler(authService)
	siteConfigHandler := handler.NewSiteConfigHandler(userRepo)
	menuHandler := handler.NewMenuHandler(menuRepo)
	permissionHandler := handler.NewPermissionHandler(permissionRepo)
	infoHandler := handler.NewInfoHandler(infoRepo, userRepo, infoCategoryRepo)
	infoCategoryHandler := handler.NewInfoCategoryHandler(infoCategoryRepo, userRepo)
	userHandler := handler.NewUserHandler(userRepo)
	roleHandler := handler.NewRoleHandler(roleRepo)
	siteGroupHandler := handler.NewSiteGroupHandler(siteGroupRepo)
	formManageHandler := handler.NewFormManageHandler(formRepo)
	systemLogHandler := handler.NewSystemLogHandler(systemLogRepo)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "ADCMS Backend is running",
		})
	})

	// API 路由组
	api := r.Group("/api")

	// 认证路由（不需要 JWT 验证）
	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// 需要认证的路由
	authorized := api.Group("")
	authorized.Use(middleware.AuthMiddleware())
	authorized.Use(middleware.TenantMiddleware())
	authorized.Use(middleware.OperationLogMiddleware(db))
	{
		// 获取用户信息和权限码
		authorized.GET("/auth/codes", authHandler.GetAccessCodes)
		authorized.POST("/auth/impersonate", authHandler.Impersonate)

		// 菜单接口
		menu := authorized.Group("/menu")
		{
			menu.GET("/all", menuHandler.GetAllMenus) // 系统导航菜单
			registerCRUDRoutes(menu, CRUDHandlers{
				List:        withPermission("system:menu:list", menuHandler.GetMenuList),
				Create:      withPermission("system:menu:update", menuHandler.CreateMenu),
				Update:      withPermission("system:menu:update", menuHandler.UpdateMenu),
				Delete:      withPermission("system:menu:update", menuHandler.DeleteMenu),
				Status:      withPermission("system:menu:update", menuHandler.UpdateMenuStatus),
				BatchDelete: withPermission("system:menu:update", menuHandler.BatchDeleteMenu),
			})
			menu.POST("/show", withPermission("system:menu:update", menuHandler.UpdateMenuShow))
		}

		// 权限管理
		permissions := authorized.Group("/permission")
		{
			registerCRUDRoutes(permissions, CRUDHandlers{
				List:        withPermission("system:permission:list", permissionHandler.GetPermissionList),
				Create:      withPermission("system:permission:list", permissionHandler.CreatePermission),
				Update:      withPermission("system:permission:list", permissionHandler.UpdatePermission),
				Delete:      withPermission("system:permission:list", permissionHandler.DeletePermission),
				BatchDelete: withPermission("system:permission:list", permissionHandler.BatchDeletePermission),
			})
		}

		// 用户管理
		users := authorized.Group("/user")
		{
			registerCRUDRoutes(users, CRUDHandlers{
				List:        withPermission("system:user:list", userHandler.GetUserList),
				Create:      withPermission("system:user:create", userHandler.CreateUser),
				Update:      withPermission("system:user:update", userHandler.UpdateUserByBody),
				Delete:      withPermission("system:user:delete", userHandler.DeleteUserByBody),
				Status:      withPermission("system:user:update", userHandler.UpdateUserStatus),
				BatchDelete: withPermission("system:user:delete", userHandler.BatchDeleteUsers),
			})
			users.POST("/reset-password", withPermission("system:user:update", userHandler.ResetPassword))
			users.GET("/info", authHandler.GetUserInfo)
			users.GET("/:id", userHandler.GetUser)
			users.POST("/profile/update", userHandler.UpdateProfile)
			users.GET("/security/settings", userHandler.GetSecuritySettings)
			users.POST("/security/update", userHandler.UpdateSecuritySettings)
			users.POST("/password/update", userHandler.UpdatePassword)
			users.GET("/phone/settings", userHandler.GetPhoneSetting)
			users.POST("/phone/update", userHandler.UpdatePhoneSetting)
			users.GET("/question/settings", userHandler.GetQuestionSetting)
			users.POST("/question/update", userHandler.UpdateQuestionSetting)
			users.GET("/email/settings", userHandler.GetEmailSetting)
			users.POST("/email/update", userHandler.UpdateEmailSetting)
			users.GET("/google-auth/settings", userHandler.GetGoogleAuthSetting)
			users.POST("/google-auth/bind", userHandler.BindGoogleAuth)
			users.POST("/google-auth/unbind", userHandler.UnbindGoogleAuth)
			users.GET("/notification/settings", userHandler.GetNotificationSettings)
			users.POST("/notification/update", userHandler.UpdateNotificationSettings)
		}

		// 角色管理
		roles := authorized.Group("/role")
		{
			registerCRUDRoutes(roles, CRUDHandlers{
				List:        withPermission("system:role:list", roleHandler.GetRoleList),
				Create:      withPermission("system:role:create", roleHandler.CreateRole),
				Update:      withPermission("system:role:update", roleHandler.UpdateRole),
				Delete:      withPermission("system:role:update", roleHandler.DeleteRole),
				Status:      withPermission("system:role:update", roleHandler.UpdateRoleStatus),
				BatchDelete: withPermission("system:role:update", roleHandler.BatchDeleteRole),
			})
			roles.GET("/permission/:id", roleHandler.GetRolePermission)
			roles.POST("/permission", withPermission("system:role:permission", roleHandler.UpdateRolePermission))
		}

		// 信息管理
		infos := authorized.Group("/info")
		{
			registerCRUDRoutes(infos, CRUDHandlers{
				List:        infoHandler.GetInfoList,
				Create:      infoHandler.CreateInfo,
				Update:      infoHandler.UpdateInfoByBody,
				Delete:      infoHandler.DeleteInfoByBody,
				Status:      infoHandler.UpdateInfoStatus,
				BatchDelete: infoHandler.BatchDeleteInfo,
			})
			infoCategory := infos.Group("/category")
			{
				registerCRUDRoutes(infoCategory, CRUDHandlers{
					List:        infoCategoryHandler.GetCategoryList,
					Create:      infoCategoryHandler.CreateCategory,
					Update:      infoCategoryHandler.UpdateCategory,
					Delete:      infoCategoryHandler.DeleteCategory,
					Status:      infoCategoryHandler.UpdateCategoryStatus,
					BatchDelete: infoCategoryHandler.BatchDeleteCategory,
				})
			}
			infos.GET("/detail/:id", infoHandler.GetInfoDetailByPath)
		}

		// 站群管理
		siteGroups := authorized.Group("/site-group")
		{
			registerCRUDRoutes(siteGroups, CRUDHandlers{
				List:        siteGroupHandler.GetSiteGroupList,
				BatchDelete: siteGroupHandler.BatchDeleteSiteGroup,
			})
			siteGroups.GET("/admins", siteGroupHandler.GetAdminOptions)
			siteGroups.POST("/cities", siteGroupHandler.GetCityList)
			siteGroups.POST("", siteGroupHandler.CreateSiteGroup)
			siteGroups.PUT("/:id", siteGroupHandler.UpdateSiteGroup)
			siteGroups.DELETE("/:id", siteGroupHandler.DeleteSiteGroup)
		}

		// 网站配置
		siteConfig := authorized.Group("/site-config")
		{
			siteConfig.GET("", siteConfigHandler.GetSiteConfig)
			siteConfig.POST("", siteConfigHandler.UpdateSiteConfig)
			siteConfig.POST("/logo/upload", siteConfigHandler.UploadSiteLogo)
		}

		// 表单管理
		forms := authorized.Group("/form-manage")
		{
			forms.POST("/list", formManageHandler.GetFormList)
			forms.POST("/batch-delete", formManageHandler.BatchDeleteForm)
			forms.DELETE("/:id", formManageHandler.DeleteForm)
			forms.POST("/export", formManageHandler.ExportForm)
		}

		// 系统日志
		logs := authorized.Group("/system-logs")
		{
			logs.POST("/list", systemLogHandler.GetSystemLogList)
			logs.POST("/batch-delete", systemLogHandler.BatchDeleteSystemLog)
			logs.DELETE("/:id", systemLogHandler.DeleteSystemLog)
			logs.POST("/clear", systemLogHandler.ClearSystemLog)
		}
	}
}
