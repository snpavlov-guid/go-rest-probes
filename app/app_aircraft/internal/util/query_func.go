package util

import (
	"fmt"
	"strings"
	"github.com/snpavlov/app_aircraft/internal/model"
)
	
func AddInClause(srcquery string, values []string, field string, keyword string) (string) {
	svalues := Map(values, func(p string) string {
        return `'` + p + `'`
    })
	return fmt.Sprintf("%s %s (%s IN (%s)) ", srcquery, keyword, field, strings.Join(svalues, ",")) 
}

func AddOrderByClause(srcquery string, fields []model.OrderInfo) (string) {
	svalues := Map(fields, func(p model.OrderInfo) string {
		clause := p.Field
		if p.Desc {
			clause = clause + ` DESC` 
		}
        return clause
    })
	return fmt.Sprintf("%s ORDER BY %s) ", srcquery, strings.Join(svalues, ",")) 
}

func AddGroupClause(srcquery string, fields []string) (string) {
	return fmt.Sprintf("%s GROUP BY %s) ", srcquery, strings.Join(fields, ",")) 
}

func AddPaginationClause(srcquery string, pager model.PageInfo) (string, []interface{}) {
	var args []interface{}
	dstquery := srcquery
    paramCount := 0
    
    if pager.Limit != nil {
        paramCount++
        dstquery += fmt.Sprintf(" LIMIT $%d", paramCount)
        args = append(args, *pager.Limit)
    }
    
    if pager.Offset != nil {
        paramCount++
        dstquery += fmt.Sprintf(" OFFSET $%d", paramCount)
        args = append(args, *pager.Offset)
    }

	return dstquery, args
}