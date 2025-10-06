package model

type PageInfo struct {
    Limit  *int `form:"size"`
    Offset *int `form:"offset"`
}

