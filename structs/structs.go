package structs

type Request struct {
	Expression string `json:"expression"`
}
type ResponseOK struct {
	Id int `json:"id"`
}
type ResponseERROR struct {
	Error string `json:"error"`
}
type ResponseResult struct {
	Id     int     `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}
type ResponseExpression struct {
	Expression ResponseResult `json:"expression"`
}
type AgentResponse struct {
	Arg1           float64 `json:"arg1"`
	Arg2           float64 `json:"arg2"`
	Operation      string  `json:"operation"`
	Operation_time int     `json:"operation_time"`
}

type AgentResult struct {
	Result float64 `json:"result"`
}

type User struct {
    Login    string `json:"login"`
    Password string `json:"password"`
}
