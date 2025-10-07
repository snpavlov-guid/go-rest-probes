package model

type PageInfo struct {
    Limit  *int `form:"size"`
    Offset *int `form:"offset"`
}

type OrderInfo struct {
    Field  string
    Desc bool
}
