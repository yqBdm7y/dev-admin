package dadmin

import (
	d "github.com/yqBdm7y/devtool"
)

func Err(err error) (i d.LibraryApi) {
	api := d.Api[d.LibraryApi]{}.Get()
	api.Response.Error = err.Error()
	return api
}

func Success(data interface{}) (i d.LibraryApi) {
	api := d.Api[d.LibraryApi]{}.Get()
	api.Response.Data = data
	return api
}

func Pagination(page, page_size, total int, data interface{}) (m interface{}) {
	pg := d.Pagination[d.LibraryPagination]{}.Get()
	pg.Page = page
	pg.PageSize = page_size
	pg.Total = int(total)
	pg.DataList = data
	api := d.Api[d.LibraryApi]{}.Get()
	return api.Pagination(pg)
}
