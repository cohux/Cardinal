package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/vidar-team/Cardinal/src/conf"
	"github.com/vidar-team/Cardinal/src/frontend"
	"github.com/vidar-team/Cardinal/src/locales"
	"github.com/vidar-team/Cardinal/src/utils"
)

func (s *Service) initRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders: []string{"Authorization", "Content-type", "User-Agent"},
		AllowOrigins: []string{"*"},
	}))

	api := r.Group("/api")
	api.Use(locales.Middleware())

	// Frontend
	r.Use(static.Serve("/", frontend.FS()))

	// Cardinal basic info
	api.Any("/", func(c *gin.Context) {
		c.JSON(utils.MakeSuccessJSON("Cardinal"))
	})

	api.GET("/base", func(c *gin.Context) {
		c.JSON(utils.MakeSuccessJSON(gin.H{
			"Title": conf.Get().Title,
		}))
	})
	api.GET("/time", func(c *gin.Context) {
		c.JSON(s.getTime())
	})

	// Static files
	api.Static("/uploads", "./uploads")

	// Team login
	api.POST("/login", func(c *gin.Context) {
		c.JSON(s.TeamLogin(c))
	})
	// Team logout
	api.GET("/logout", func(c *gin.Context) {
		c.JSON(s.TeamLogout(c))
	})

	// Submit flag
	api.POST("/flag", func(c *gin.Context) {
		c.JSON(s.SubmitFlag(c))
	})

	// For team
	team := api.Group("/team")
	team.Use(s.TeamAuthRequired())
	{
		team.GET("/info", func(c *gin.Context) {
			c.JSON(s.GetTeamInfo(c))
		})
		team.GET("/gameboxes", func(c *gin.Context) {
			c.JSON(s.GetSelfGameBoxes(c))
		})
		team.GET("/rank", func(c *gin.Context) {
			c.JSON(utils.MakeSuccessJSON(gin.H{"Title": s.GetRankListTitle(), "Rank": s.GetRankList()}))
		})
		team.GET("/bulletins", func(c *gin.Context) {
			c.JSON(s.GetAllBulletins())
		})
	}

	// Manager login
	api.POST("/manager/login", func(c *gin.Context) {
		c.JSON(s.ManagerLogin(c))
	})

	// Manager logout
	api.GET("/manager/logout", func(c *gin.Context) {
		c.JSON(s.ManagerLogout(c))
	})

	// For manager
	check := api.Group("/manager").Use(s.AdminAuthRequired())
	manager := api.Group("/manager").Use(s.AdminAuthRequired(), s.ManagerRequired())
	{
		// Challenge
		manager.GET("/challenges", func(c *gin.Context) {
			c.JSON(s.GetAllChallenges())
		})
		manager.POST("/challenge", func(c *gin.Context) {
			c.JSON(s.NewChallenge(c))
		})
		manager.PUT("/challenge", func(c *gin.Context) {
			c.JSON(s.EditChallenge(c))
		})
		manager.DELETE("/challenge", func(c *gin.Context) {
			c.JSON(s.DeleteChallenge(c))
		})
		manager.POST("/challenge/visible", func(c *gin.Context) {
			c.JSON(s.SetVisible(c))
		})

		// GameBox
		manager.GET("/gameboxes", func(c *gin.Context) {
			c.JSON(s.GetGameBoxes(c))
		})
		manager.POST("/gameboxes", func(c *gin.Context) {
			c.JSON(s.NewGameBoxes(c))
		})
		manager.PUT("/gamebox", func(c *gin.Context) {
			c.JSON(s.EditGameBox(c))
		})
		manager.GET("/gameboxes/sshTest", func(c *gin.Context) {
			c.JSON(s.testSSH(c))
		})

		// Team
		manager.GET("/teams", func(c *gin.Context) {
			c.JSON(s.GetAllTeams())
		})
		manager.POST("/teams", func(c *gin.Context) {
			c.JSON(s.NewTeams(c))
		})
		manager.PUT("/team", func(c *gin.Context) {
			c.JSON(s.EditTeam(c))
		})
		manager.DELETE("/team", func(c *gin.Context) {
			c.JSON(s.DeleteTeam(c))
		})
		manager.POST("/team/resetPassword", func(c *gin.Context) {
			c.JSON(s.ResetTeamPassword(c))
		})

		// Manager
		manager.GET("/managers", func(c *gin.Context) {
			c.JSON(s.GetAllManager())
		})
		manager.POST("/manager", func(c *gin.Context) {
			c.JSON(s.NewManager(c))
		})
		manager.GET("/manager/token", func(c *gin.Context) {
			c.JSON(s.RefreshManagerToken(c))
		})
		manager.GET("/manager/changePassword", func(c *gin.Context) {
			c.JSON(s.ChangeManagerPassword(c))
		})
		manager.DELETE("/manager", func(c *gin.Context) {
			c.JSON(s.DeleteManager(c))
		})

		// Flag
		manager.GET("/flags", func(c *gin.Context) {
			c.JSON(s.GetFlags(c))
		})
		manager.POST("/flag/generate", func(c *gin.Context) {
			c.JSON(s.GenerateFlag(c))
		})
		manager.GET("/flag/export", func(c *gin.Context) {
			c.JSON(s.ExportFlag(c))
		})

		// Check
		check.POST("/checkDown", func(c *gin.Context) {
			c.JSON(s.CheckDown(c))
		})

		// Bulletin
		manager.GET("/bulletins", func(c *gin.Context) {
			c.JSON(s.GetAllBulletins())
		})
		manager.POST("/bulletin", func(c *gin.Context) {
			c.JSON(s.NewBulletin(c))
		})
		manager.PUT("/bulletin", func(c *gin.Context) {
			c.JSON(s.EditBulletin(c))
		})
		manager.DELETE("/bulletin", func(c *gin.Context) {
			c.JSON(s.DeleteBulletin(c))
		})

		// File
		manager.POST("/uploadPicture", func(c *gin.Context) {
			c.JSON(s.UploadPicture(c))
		})

		// Log
		manager.GET("/logs", func(c *gin.Context) {
			c.JSON(s.GetLogs(c))
		})
		manager.GET("/rank", func(c *gin.Context) {
			c.JSON(utils.MakeSuccessJSON(gin.H{"Title": s.GetRankListTitle(), "Rank": s.GetManagerRankList()}))
		})
		manager.GET("/panel", func(c *gin.Context) {
			c.JSON(s.Panel(c))
		})
	}

	// 404
	r.NoRoute(func(c *gin.Context) {
		c.JSON(utils.MakeErrJSON(404, 40400,
			locales.I18n.T(c.GetString("lang"), "general.not_found"),
		))
	})

	// 405
	r.NoMethod(func(c *gin.Context) {
		c.JSON(utils.MakeErrJSON(405, 40500,
			locales.I18n.T(c.GetString("lang"), "general.method_not_allow"),
		))
	})

	return r
}

// TeamAuthRequired is the team permission check middleware.
func (s *Service) TeamAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(utils.MakeErrJSON(403, 40300,
				locales.I18n.T(c.GetString("lang"), "general.no_auth"),
			))
			c.Abort()
			return
		}

		var tokenData Token
		s.Mysql.Where(&Token{Token: token}).Find(&tokenData)
		if tokenData.ID == 0 {
			c.JSON(utils.MakeErrJSON(401, 40100,
				locales.I18n.T(c.GetString("lang"), "general.no_auth"),
			))
			c.Abort()
			return
		}

		c.Set("teamID", tokenData.TeamID)
		c.Next()
	}
}

// AdminAuthRequired is the admin permission check middleware.
func (s *Service) AdminAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(utils.MakeErrJSON(403, 40300,
				locales.I18n.T(c.GetString("lang"), "general.no_auth"),
			))
			c.Abort()
			return
		}

		var managerData Manager
		s.Mysql.Where(&Manager{Token: token}).Find(&managerData)
		if managerData.ID == 0 {
			c.JSON(utils.MakeErrJSON(401, 40100,
				locales.I18n.T(c.GetString("lang"), "general.no_auth"),
			))
			c.Abort()
			return
		}

		c.Set("managerData", managerData)
		c.Set("isCheck", managerData.IsCheck)
		c.Next()
	}
}

// ManagerRequired make sure the account is the manager.
func (s *Service) ManagerRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool("isCheck") {
			c.JSON(utils.MakeErrJSON(401, 40100,
				locales.I18n.T(c.GetString("lang"), "manager.manager_required"),
			))
			c.Abort()
			return
		}
		c.Next()
	}
}
