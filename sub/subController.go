package sub

import (
	"encoding/base64"
	"net"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type SUBController struct {
	subTitle       string
	subAnnounce    string
	subPath        string
	subJsonPath    string
	subEncrypt     bool
	updateInterval string

	subService     *SubService
	subJsonService *SubJsonService
}

func NewSUBController(
	g *gin.RouterGroup,
	subPath string,
	jsonPath string,
	encrypt bool,
	showInfo bool,
	rModel string,
	update string,
	jsonFragment string,
	jsonNoise string,
	jsonMux string,
	jsonRules string,
	subTitle string,
	subAnnounce string,
) *SUBController {
	sub := NewSubService(showInfo, rModel)
	a := &SUBController{
		subTitle:       subTitle,
		subAnnounce:    subAnnounce,
		subPath:        subPath,
		subJsonPath:    jsonPath,
		subEncrypt:     encrypt,
		updateInterval: update,

		subService:     sub,
		subJsonService: NewSubJsonService(jsonFragment, jsonNoise, jsonMux, jsonRules, sub),
	}
	a.initRouter(g)
	return a
}

func (a *SUBController) initRouter(g *gin.RouterGroup) {
	gLink := g.Group(a.subPath)
	gJson := g.Group(a.subJsonPath)

	gLink.GET(":subid", a.subs)

	gJson.GET(":subid", a.subJsons)
}

func (a *SUBController) subs(c *gin.Context) {
	subId := c.Param("subid")
	var host string
	if h, err := getHostFromXFH(c.GetHeader("X-Forwarded-Host")); err == nil {
		host = h
	}
	if host == "" {
		host = c.GetHeader("X-Real-IP")
	}
	if host == "" {
		var err error
		host, _, err = net.SplitHostPort(c.Request.Host)
		if err != nil {
			host = c.Request.Host
		}
	}
	var supportUrl string
	supportUrl = os.Getenv("XUI_SUB_SUPPORT_URL")
	if supportUrl == "" {
		supportUrl = os.Getenv("XUI_SUB_DOMAIN")
	}
	var profileWebPageUrl string
	profileWebPageUrl = os.Getenv("XUI_SUB_PROFILE_WEB_PAGE_URL")
	if profileWebPageUrl == "" {
		profileWebPageUrl = os.Getenv("XUI_SUB_DOMAIN")
	}
	announceText := a.subAnnounce
	subs, header, err := a.subService.GetSubs(subId, host)
	if err != nil || len(subs) == 0 {
		c.String(400, "Error!")
	} else {
		result := ""
		for _, sub := range subs {
			result += sub + "\n"
		}

		// Add headers
		c.Writer.Header().Set("Subscription-Userinfo", header)
		c.Writer.Header().Set("Profile-Update-Interval", a.updateInterval)
		c.Writer.Header().Set("Profile-Title", "base64:"+base64.StdEncoding.EncodeToString([]byte(a.subTitle)))
		c.Writer.Header().Set("Support-Url", supportUrl)
		c.Writer.Header().Set("Profile-Web-Page-Url", profileWebPageUrl)
		if announceText != "" {
			c.Writer.Header().Set("Announce", announceText)
		}

		if a.subEncrypt {
			c.String(200, base64.StdEncoding.EncodeToString([]byte(result)))
		} else {
			c.String(200, result)
		}
	}
}

func (a *SUBController) subJsons(c *gin.Context) {
	subId := c.Param("subid")
	var host string
	if h, err := getHostFromXFH(c.GetHeader("X-Forwarded-Host")); err == nil {
		host = h
	}
	if host == "" {
		host = c.GetHeader("X-Real-IP")
	}
	if host == "" {
		var err error
		host, _, err = net.SplitHostPort(c.Request.Host)
		if err != nil {
			host = c.Request.Host
		}
	}
	var supportUrl string
	supportUrl = os.Getenv("XUI_SUB_SUPPORT_URL")
	if supportUrl == "" {
		supportUrl = os.Getenv("XUI_SUB_DOMAIN")
	}
	var profileWebPageUrl string
	profileWebPageUrl = os.Getenv("XUI_SUB_PROFILE_WEB_PAGE_URL")
	if profileWebPageUrl == "" {
		profileWebPageUrl = os.Getenv("XUI_SUB_DOMAIN")
	}
	announceText := a.subAnnounce
	jsonSub, header, err := a.subJsonService.GetJson(subId, host)
	if err != nil || len(jsonSub) == 0 {
		c.String(400, "Error!")
	} else {

		// Add headers
		c.Writer.Header().Set("Subscription-Userinfo", header)
		c.Writer.Header().Set("Profile-Update-Interval", a.updateInterval)
		c.Writer.Header().Set("Profile-Title", "base64:"+base64.StdEncoding.EncodeToString([]byte(a.subTitle)))
		c.Writer.Header().Set("Support-Url", supportUrl)
		c.Writer.Header().Set("Profile-Web-Page-Url", profileWebPageUrl)
		if announceText != "" {
			c.Writer.Header().Set("Announce", announceText)
		}

		c.String(200, jsonSub)
	}
}

func getHostFromXFH(s string) (string, error) {
	if strings.Contains(s, ":") {
		realHost, _, err := net.SplitHostPort(s)
		if err != nil {
			return "", err
		}
		return realHost, nil
	}
	return s, nil
}
