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
	"reflect"
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

// CreateRoutersFromFile create a router the details of the router is in the body
// clients SHOULD NOT transmit PII (Personal Identification Information) parameters in the URL
// (as part of path or query string) because this information can be inadvertently exposed via client, network,
// and server logs and other mechanisms. (Microsoft API Guidelines)
func (a *ApiServer) CreateRoutersFromFile(c *gin.Context) {

	var (
		routers []domain.Router
		f       multipart.File
		b       []byte
		err     error
		tenant  string
		res     []byte
		form    *multipart.Form
	)

	form, _ = c.MultipartForm()
	if c.ContentType() != "multipart/form-data" {
		c.JSON(http.StatusBadRequest, "expect multipart/form-data")
		return
	}

	if form == nil {
		c.JSON(http.StatusOK, nil)
		return
	}

	item := form.File[routerfilenametag]
	for _, k := range item {
		if k.Header.Get("Content-Type") != "application/json" {
			c.JSON(http.StatusBadRequest, "expect application/json")
		}
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

	// this is coming form the authentication layer.
	tenant = c.Value("tenant").(string)

	res, err = a.next.CreateRouters(a.ctx, routers, tenant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, string(res))

}

// CreateRouters create a router the details of the router is in the body
// clients SHOULD NOT transmit PII (Personal Identification Information) parameters in the URL
// (as part of path or query string) because this information can be inadvertently exposed via client, network,
// and server logs and other mechanisms. (Microsoft API Guidelines)
func (a *ApiServer) CreateRouters(c *gin.Context) {

	var (
		routers []domain.Router
		r       domain.Router
		err     error
		tenant  string
		res     []byte
		form    *multipart.Form
		items   []string
	)

	form, _ = c.MultipartForm()

	if form == nil {
		c.JSON(http.StatusOK, nil)
		return
	}

	if c.ContentType() != "multipart/form-data" {
		c.JSON(http.StatusBadRequest, "expect multipart/form-data")
		return
	}

	items = form.Value[createRouterKey]
	for _, v := range items {
		err = json.Unmarshal([]byte(v), &r)
		routers = append(routers, r)
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

// DeleteRoutersFromFile create a router the details of the router is in the body
// clients SHOULD NOT transmit PII (Personal Identification Information) parameters in the URL
// (as part of path or query string) because this information can be inadvertently exposed via client, network,
// and server logs and other mechanisms. (Microsoft API Guidelines)
func (a *ApiServer) DeleteRoutersFromFile(c *gin.Context) {

	var (
		routers []domain.Router
		f       multipart.File
		b       []byte
		err     error
		tenant  string
		form    *multipart.Form
	)

	form, _ = c.MultipartForm()
	if c.ContentType() != "multipart/form-data" {
		c.JSON(http.StatusBadRequest, "expect multipart/form-data")
		return
	}

	if form == nil {
		c.JSON(http.StatusOK, nil)
		return
	}

	item := form.File[routerfilenametag]
	for _, k := range item {
		if k.Header.Get("Content-Type") != "application/json" {
			c.JSON(http.StatusBadRequest, "expect application/json")
		}
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

	// this is coming form the authentication layer.
	tenant = c.Value("tenant").(string)

	err = a.next.DeleteRouters(a.ctx, routers, tenant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)

}

// DeleteRouters create a router the details of the router is in the body
// clients SHOULD NOT transmit PII (Personal Identification Information) parameters in the URL
// (as part of path or query string) because this information can be inadvertently exposed via client, network,
// and server logs and other mechanisms. (Microsoft API Guidelines)
func (a *ApiServer) DeleteRouters(c *gin.Context) {

	var (
		routers []domain.Router
		r       domain.Router
		err     error
		tenant  string
		form    *multipart.Form
	)

	form, _ = c.MultipartForm()

	if form == nil {
		c.JSON(http.StatusOK, nil)
		return
	}

	if c.ContentType() != "multipart/form-data" {
		c.JSON(http.StatusBadRequest, "expect multipart/form-data")
		return
	}

	st := reflect.TypeOf(r)
	field := st.Field(0)
	item := form.Value[field.Tag.Get("json")]
	for _, v := range item {
		r.RouterSerial = v
		routers = append(routers, r)
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
	if query.Has("id") {
		var (
			r domain.Router
		)
		r.RouterSerial = query.Get("id")
		r, err = a.next.GetRouters(a.ctx, r, tenant)
		if r.RouterSerial == "" {
			c.JSON(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, r)
		}
		return
	}
	c.JSON(http.StatusBadRequest, errors.New("url parameter error").Error())
}
