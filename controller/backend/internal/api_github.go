package internal

import (
	"net/http"
	"slices"
	"strings"

	"github.com/Alonza0314/it-system/controller/backend/constant"
	"github.com/Alonza0314/it-system/controller/backend/model"

	"github.com/free-ran-ue/util"
	"github.com/gin-gonic/gin"
)

func (b *backend) getGithubRoutes() util.Routes {
	return util.Routes{
		{
			Name:        "Get Github PRs",
			Method:      http.MethodGet,
			Pattern:     "",
			HandlerFunc: b.handleGetGithubPRs,
		},
		{
			Name:        "Suggest dependency PRs",
			Method:      http.MethodPost,
			Pattern:     "/dependency-suggestions",
			HandlerFunc: b.handleDependencySuggestions,
		},
	}
}

func (b *backend) handleGetGithubPRs(c *gin.Context) {
	b.GitLog.Infof("Get Github PRs request from %s, user: %s", c.ClientIP(), c.GetHeader("user"))

	nf := c.Query("nf")
	if nf == "" {
		c.JSON(http.StatusBadRequest, model.ResponseGetGithubPRs{
			Message: "NF parameter is required",
		})
		return
	}

	if exist := slices.Contains(constant.NF_LIST, nf) || slices.Contains(constant.LIBRARY_LIST, nf); !exist {
		c.JSON(http.StatusBadRequest, model.ResponseGetGithubPRs{
			Message: "Invalid NF parameter, must be one of: " + strings.Join(append(constant.NF_LIST, constant.LIBRARY_LIST...), ", "),
		})
		return
	}

	if nf == constant.UPF {
		nf = constant.GO_UPF
	}

	response, errDetail := b.Processor.GetGithubPRs(nf)
	if errDetail != nil {
		c.JSON(errDetail.HttpStatus, model.ResponseGetGithubPRs{
			Message: errDetail.Detail,
		})
		return
	}

	b.GitLog.Infof("Get Github PRs successful for %s", c.ClientIP())
	c.JSON(http.StatusOK, response)
}

func (b *backend) handleDependencySuggestions(c *gin.Context) {
	b.GitLog.Infof("Dependency suggestions request from %s, user: %s", c.ClientIP(), c.GetHeader("user"))

	var req model.RequestDependencySuggestions
	if err := c.ShouldBindJSON(&req); err != nil {
		b.GitLog.Warnf("Invalid dependency suggestions request from %s: %v", c.ClientIP(), err)
		c.JSON(http.StatusBadRequest, model.ResponseDependencySuggestions{
			Message: "Invalid request",
		})
		return
	}

	for _, nfPr := range req.NFPrList {
		nf := nfPr.NfName
		if nf == constant.GO_UPF {
			nf = constant.UPF
		}
		if exist := slices.Contains(constant.NF_LIST, nf); !exist {
			c.JSON(http.StatusBadRequest, model.ResponseDependencySuggestions{
				Message: "Invalid NF parameter, must be one of: " + strings.Join(constant.NF_LIST, ", "),
			})
			return
		}
	}

	response, errDetail := b.Processor.SuggestLibraryPRs(&req)
	if errDetail != nil {
		c.JSON(errDetail.HttpStatus, model.ResponseDependencySuggestions{
			Message: errDetail.Detail,
		})
		return
	}

	b.GitLog.Infof("Dependency suggestions successful for %s", c.ClientIP())
	c.JSON(http.StatusOK, response)
}
