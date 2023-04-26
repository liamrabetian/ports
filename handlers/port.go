package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mohammadrabetian/ports/domain"
	"github.com/mohammadrabetian/ports/service"
	"github.com/mohammadrabetian/ports/util"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
)

var portSvc service.PortService

func InitPortHandlers(portService service.PortService) {
	portSvc = portService
}

func ProcessJSONFile(filePath string) error {
	return util.ProcessJSONFile(filePath, func(port *domain.Port) error {
		return portSvc.AddOrUpdatePort(port)
	})
}

//	@Summary		Get Port By ID API
//	@Description	Returns a Port object providing an ID

// @Tags		port
// @Accept		json
// @Produce	json
//
// @Param		id	path		string	true	"port id"
//
// @Success	200	{object}	string
// @Failure	400	{string}	httputil.HTTPError
// @Router		/api/v1/ports/{id} [get]
func GetPortByID(c *gin.Context) {

	/* NOTES: validations */
	portID := c.Param("id")

	port, err := portSvc.GetPortByID(portID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Port not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"port": port})
}

//	@Summary		List ports API
//	@Description	Returns a list of ports

// @Tags		port
// @Accept		json
// @Produce	json
//
// @Param		page	query		string	false	"page query"
//
// @Param		limit	query		string	false	"page limit"
//
// @Success	200		{object}	string
// @Failure	400		{string}	httputil.HTTPError
// @Router		/api/v1/ports [get]
func ListPorts(c *gin.Context) {

	/* NOTES: validations */
	pageStr := c.Request.URL.Query().Get("page")
	limitStr := c.Request.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = DefaultPage
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = DefaultLimit
	}

	ports, err := portSvc.ListPorts(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve the ports"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ports": ports})
}
