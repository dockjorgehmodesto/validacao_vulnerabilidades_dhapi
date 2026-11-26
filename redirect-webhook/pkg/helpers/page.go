package helpers

type PageResponse struct {
	TotalRecord int `json:"totalRecords" description:"Total amount of records within the given query"`
	TotalPage   int `json:"totalPages"   description:"Total amount of pages"`
	Offset      int `json:"offset"       description:"Records skipped to display the current page"`
	Limit       int `json:"limit"        description:"Amount of records displayed per page"`
	Page        int `json:"page"         description:"The position of the current page"`
	PrevPage    int `json:"prevPage"     description:"The position of the previous page"`
	NextPage    int `json:"nextPage"     description:"The position of the next page"`
}
