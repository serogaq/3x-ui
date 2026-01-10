package controller

import (
	"fmt"
	"strconv"

	"github.com/mhsanaei/3x-ui/v2/database/model"
	"github.com/mhsanaei/3x-ui/v2/web/service"
	"github.com/mhsanaei/3x-ui/v2/web/session"

	"github.com/gin-gonic/gin"
)

type ExternalInboundController struct {
	externalInboundService service.ExternalInboundService
}

type externalInboundClientRequest struct {
	Id       int    `json:"id" form:"id"`
	Settings string `json:"settings" form:"settings"`
}

type externalInboundResponse struct {
	Id      int                           `json:"id"`
	Enable  bool                          `json:"enable"`
	Name    string                        `json:"name"`
	Comment string                        `json:"comment"`
	Url     string                        `json:"url"`
	Clients []model.ExternalInboundClient `json:"clients"`
}

// NewExternalInboundController creates a new ExternalInboundController and sets up its routes.
func NewExternalInboundController(g *gin.RouterGroup) *ExternalInboundController {
	c := &ExternalInboundController{}
	c.initRouter(g)
	return c
}

func (c *ExternalInboundController) initRouter(g *gin.RouterGroup) {
	g.GET("/list", c.getExternalInbounds)
	g.POST("/add", c.addExternalInbound)
	g.POST("/update/:id", c.updateExternalInbound)
	g.POST("/del/:id", c.delExternalInbound)
	g.POST("/addClient", c.addExternalInboundClient)
	g.POST("/addClientAll", c.addExternalInboundClientAll)
	g.POST("/updateClient/:clientId", c.updateExternalInboundClient)
	g.POST("/:id/delClient/:clientId", c.delExternalInboundClient)
	g.POST("/delClientsBySub/:subId", c.delExternalInboundClientsBySub)
}

func (c *ExternalInboundController) getExternalInbounds(ctx *gin.Context) {
	user := session.GetLoginUser(ctx)
	externalInbounds, err := c.externalInboundService.GetExternalInbounds(user.Id)
	if err != nil {
		jsonMsg(ctx, I18nWeb(ctx, "pages.inbounds.toasts.obtain"), err)
		return
	}
	responses := make([]externalInboundResponse, 0, len(externalInbounds))
	for _, externalInbound := range externalInbounds {
		clients, err := c.externalInboundService.GetClients(externalInbound.Id)
		if err != nil {
			jsonMsg(ctx, I18nWeb(ctx, "pages.inbounds.toasts.obtain"), err)
			return
		}
		responses = append(responses, externalInboundResponse{
			Id:      externalInbound.Id,
			Enable:  externalInbound.Enable,
			Name:    externalInbound.Name,
			Comment: externalInbound.Comment,
			Url:     externalInbound.Url,
			Clients: clients,
		})
	}
	jsonObj(ctx, responses, nil)
}

func (c *ExternalInboundController) addExternalInbound(ctx *gin.Context) {
	externalInbound := &model.ExternalInbound{}
	if err := ctx.ShouldBind(externalInbound); err != nil {
		jsonMsg(ctx, I18nWeb(ctx, "pages.inbounds.toasts.inboundCreateSuccess"), err)
		return
	}
	user := session.GetLoginUser(ctx)
	externalInbound.UserId = user.Id
	if err := c.externalInboundService.AddExternalInbound(externalInbound); err != nil {
		jsonMsg(ctx, I18nWeb(ctx, "somethingWentWrong"), err)
		return
	}
	jsonMsgObj(ctx, I18nWeb(ctx, "pages.inbounds.toasts.inboundCreateSuccess"), externalInbound, nil)
}

func (c *ExternalInboundController) updateExternalInbound(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		jsonMsg(ctx, I18nWeb(ctx, "update"), err)
		return
	}
	payload := &model.ExternalInbound{}
	if err := ctx.ShouldBind(payload); err != nil {
		jsonMsg(ctx, I18nWeb(ctx, "update"), err)
		return
	}
	externalInbound, err := c.externalInboundService.GetExternalInbound(id)
	if err != nil {
		jsonMsg(ctx, I18nWeb(ctx, "update"), err)
		return
	}
	externalInbound.Enable = payload.Enable
	externalInbound.Name = payload.Name
	externalInbound.Comment = payload.Comment
	externalInbound.Url = payload.Url
	if err := c.externalInboundService.UpdateExternalInbound(externalInbound); err != nil {
		jsonMsg(ctx, I18nWeb(ctx, "somethingWentWrong"), err)
		return
	}
	jsonMsg(ctx, I18nWeb(ctx, "update"), nil)
}

func (c *ExternalInboundController) delExternalInbound(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		jsonMsg(ctx, I18nWeb(ctx, "delete"), err)
		return
	}
	if err := c.externalInboundService.DeleteExternalInbound(id); err != nil {
		jsonMsg(ctx, I18nWeb(ctx, "delete"), err)
		return
	}
	jsonMsg(ctx, I18nWeb(ctx, "delete"), nil)
}

func (c *ExternalInboundController) addExternalInboundClient(ctx *gin.Context) {
	payload := &externalInboundClientRequest{}
	if err := ctx.ShouldBind(payload); err != nil {
		jsonMsg(ctx, I18nWeb(ctx, "pages.inbounds.toasts.inboundClientAddSuccess"), err)
		return
	}
	if err := c.externalInboundService.AddExternalInboundClient(payload.Id, payload.Settings); err != nil {
		jsonMsg(ctx, I18nWeb(ctx, "pages.inbounds.toasts.inboundClientAddSuccess"), err)
		return
	}
	jsonMsg(ctx, I18nWeb(ctx, "pages.inbounds.toasts.inboundClientAddSuccess"), nil)
}

func (c *ExternalInboundController) addExternalInboundClientAll(ctx *gin.Context) {
	payload := &externalInboundClientRequest{}
	if err := ctx.ShouldBind(payload); err != nil {
		jsonMsg(ctx, I18nWeb(ctx, "pages.inbounds.toasts.inboundClientAddSuccess"), err)
		return
	}
	user := session.GetLoginUser(ctx)
	if err := c.externalInboundService.AddExternalInboundClientToAll(user.Id, payload.Settings); err != nil {
		jsonMsg(ctx, I18nWeb(ctx, "pages.inbounds.toasts.inboundClientAddSuccess"), err)
		return
	}
	jsonMsg(ctx, I18nWeb(ctx, "pages.inbounds.toasts.inboundClientAddSuccess"), nil)
}

func (c *ExternalInboundController) updateExternalInboundClient(ctx *gin.Context) {
	clientId := ctx.Param("clientId")
	payload := &externalInboundClientRequest{}
	if err := ctx.ShouldBind(payload); err != nil {
		jsonMsg(ctx, I18nWeb(ctx, "pages.inbounds.toasts.inboundClientUpdateSuccess"), err)
		return
	}
	if err := c.externalInboundService.UpdateExternalInboundClient(payload.Id, clientId, payload.Settings); err != nil {
		jsonMsg(ctx, I18nWeb(ctx, "pages.inbounds.toasts.inboundClientUpdateSuccess"), err)
		return
	}
	jsonMsg(ctx, I18nWeb(ctx, "pages.inbounds.toasts.inboundClientUpdateSuccess"), nil)
}

func (c *ExternalInboundController) delExternalInboundClient(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		jsonMsg(ctx, I18nWeb(ctx, "pages.inbounds.toasts.inboundClientDeleteSuccess"), err)
		return
	}
	clientId := ctx.Param("clientId")
	if err := c.externalInboundService.DeleteExternalInboundClient(id, clientId); err != nil {
		jsonMsg(ctx, I18nWeb(ctx, "pages.inbounds.toasts.inboundClientDeleteSuccess"), err)
		return
	}
	jsonMsg(ctx, I18nWeb(ctx, "pages.inbounds.toasts.inboundClientDeleteSuccess"), nil)
}

func (c *ExternalInboundController) delExternalInboundClientsBySub(ctx *gin.Context) {
	subId := ctx.Param("subId")
	if subId == "" {
		jsonMsg(ctx, I18nWeb(ctx, "pages.inbounds.toasts.inboundClientDeleteSuccess"), fmt.Errorf("missing subscription"))
		return
	}
	user := session.GetLoginUser(ctx)
	if err := c.externalInboundService.DeleteExternalInboundClientsBySubId(user.Id, subId); err != nil {
		jsonMsg(ctx, I18nWeb(ctx, "pages.inbounds.toasts.inboundClientDeleteSuccess"), err)
		return
	}
	jsonMsg(ctx, I18nWeb(ctx, "pages.inbounds.toasts.inboundClientDeleteSuccess"), nil)
}
