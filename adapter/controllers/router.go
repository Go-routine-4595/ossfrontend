package controllers

import (
	"encoding/json"
	"errors"
	"github.com/Go-routine-4995/ossfrontend/domain"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

const (
	createRouterKey = "router"
)

// for Paged request
type Links struct {
	Self     string `json:"self"`
	Previous string `json:"previous"`
	Next     string `json:"next"`
	First    string `json:"first"`
	Last     string `json:"last"`
}

// CreateRouters create a router the details of the router is in the body
// clients SHOULD NOT transmit PII (Personal Identification Information) parameters in the URL
// (as part of path or query string) because this information can be inadvertently exposed via client, network,
// and server logs and other mechanisms. (Microsoft API Guidelines)
func (a *ApiServer) CreateRouters(c *gin.Context) {

	var (
		routers []domain.Router
		err     error
		tenant  string
		res     []byte
		b       []byte
		form    *multipart.Form
		item    []*multipart.FileHeader
		f       multipart.File
	)

	if c.ContentType() == "multipart/form-data" {
		form, _ = c.MultipartForm()
		if form == nil {
			c.JSON(http.StatusOK, nil)
			return
		}
		item = form.File[routerfilenametag]
		for _, k := range item {
			f, err = k.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			b, err = io.ReadAll(f)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			if len(b) > maxMessageSize {
				c.JSON(http.StatusInternalServerError, errors.New("file size too large"))
				return
			}
			err = json.Unmarshal(b, &routers)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
			}
			//
		}
	}
	if c.ContentType() == "application/json" {
		err = c.BindJSON(&routers)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}
	}

	// this is coming form the authentication layer.
	tenant = c.Value("tenant").(string)

	res, err = a.next.CreateRouters(a.ctx, routers, tenant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, string(res))

}

// DeleteRouters create a router the details of the router is in the body
// clients SHOULD NOT transmit PII (Personal Identification Information) parameters in the URL
// (as part of path or query string) because this information can be inadvertently exposed via client, network,
// and server logs and other mechanisms. (Microsoft API Guidelines)
func (a *ApiServer) DeleteRouters(c *gin.Context) {

	var (
		routers []domain.Router
		b       []byte
		err     error
		tenant  string
		f       multipart.File
		form    *multipart.Form
	)

	form, _ = c.MultipartForm()
	if c.ContentType() == "multipart/form-data" {
		form, _ = c.MultipartForm()
		if form == nil {
			c.JSON(http.StatusOK, nil)
			return
		}
		item := form.File[routerfilenametag]
		for _, k := range item {
			f, err = k.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			b, err = io.ReadAll(f)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			if len(b) > maxMessageSize {
				c.JSON(http.StatusInternalServerError, errors.New("file size too large"))
				return
			}
			err = json.Unmarshal(b, &routers)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
			}
			//
		}

	}

	if c.ContentType() == "application/json" {
		err = c.BindJSON(&routers)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}
	}

	// this is coming form the authentication layer.
	tenant = c.Value("tenant").(string)

	err = a.next.DeleteRouters(a.ctx, routers, tenant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)

}

func (a *ApiServer) GetRouters(c *gin.Context) {

	var (
		page   Pagination
		b      []byte
		r      domain.Response
		err    error
		tenant string
		query  url.Values
		links  Links
	)

	tenant = c.Value("tenant").(string)

	query = c.Request.URL.Query()
	if query.Has("page") {
		page = GeneratePaginationFromRequest(query)

		b, err = json.Marshal(page)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		r, err = a.next.GetRoutersPage(a.ctx, b, tenant)

		links.First = "0"
		links.Last = strconv.Itoa(r.Last)
		links.Self = strconv.Itoa(page.Page)
		if page.Page == 0 {
			links.Previous = strconv.Itoa(page.Page)
		} else {
			links.Previous = strconv.Itoa(page.Page - 1)
		}
		if page.Page == r.Last {
			links.Next = strconv.Itoa(page.Page)
		} else {
			links.Next = strconv.Itoa(page.Page + 1)
		}

		c.JSON(http.StatusOK, gin.H{
			"records":  r.Routers,
			"metadata": links,
		})
		return
	}
	if query.Has("id") || query.Has("serial-id") {
		var (
			r domain.Router
		)
		r.RouterSerial = query.Get("id")
		r, err = a.next.GetRouters(a.ctx, r, tenant)
		if r.RouterSerial == "" {
			c.JSON(http.StatusOK, nil)
		} else {
			c.JSON(http.StatusOK, r)
		}
		return
	}
	c.JSON(http.StatusBadRequest, errors.New("url parameter error").Error())
}
