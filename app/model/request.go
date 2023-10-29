package model

type RequestBody struct {  
	SQL string `json:"sql"`  
} 

type RedisRequestBody struct {
	ID int64 `json:"ID"`
	Command string `json:"Command"`
}

type ZsetBody struct {
	ItemNam string `json:"name"`
	Score   float64	`json:"score"`
	Command string `json:"Command"`
	offset int64 `json:"offset"`
}