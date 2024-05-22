package commons

import "net/http"


func Pagination(page int, pageSize int) (startPage int, endPage int){
	
	startPage = page*pageSize
	endPage = ((page+1)*pageSize) - 1 

	return startPage, endPage
}

type ServiceInfo struct {
	Name string `json:"name"`
	Version string `json:"version"`
}


func IsHTTPRequestSuccess(statusCode int) bool {
	if(statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices){
		return true 
	}
	return false
}