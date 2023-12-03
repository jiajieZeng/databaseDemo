package model

type RequestBody struct {  
	SQL string `json:"sql"`  
} 

type RedisRequestBody struct {
	ID string `json:"ID"`
	Command string `json:"Command"`
	ItemNam string `json:"name"`
	Score   string	`json:"score"`
}
