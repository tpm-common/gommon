package sqlhelper

import (
	"fmt"
	"reflect"
	"strings"
)

func SetInsert(ins [][2]any) (cols, nums string, args []any) {
	for _, val := range ins {
		if reflect.TypeOf(val[0]).Kind() == reflect.Pointer {
			if reflect.ValueOf(val[0]).IsNil() {
				continue
			}
		}

		if len(args) > 0 {
			cols += `, `
			nums += `, `
		}

		cols += fmt.Sprint(val[1])
		args = append(args, val[0])
		nums += fmt.Sprintf(`$%v`, len(args))
	}
	return
}

type WhereCondition struct {
	Expression string
	Value      interface{}
}

func (wc WhereCondition) Set(sql *string, args *[]any, addWhere bool) {
	if wc.Expression == "" {
		return
	}

	if addWhere {
		*sql += ` WHERE`
	}

	*sql += fmt.Sprintf(` %s`, wc.Expression)

	*args = append(*args, wc.Value)
	*sql += fmt.Sprintf(` $%v`, len(*args))
}

type WhereLike struct {
	QueryString string
	Column      string
	Inverse     bool
	Insensitive bool
}

func (wl WhereLike) Set(sql *string, args *[]any, addWhere bool) {
	if wl.QueryString == "" || wl.Column == "" {
		return
	}

	if addWhere {
		*sql += ` WHERE`
	}

	*sql += fmt.Sprintf(` %s`, wl.Column)

	if wl.Inverse {
		*sql += ` NOT`
	}

	if wl.Insensitive {
		*sql += ` ILIKE`
	} else {
		*sql += ` LIKE`
	}

	*args = append(*args, wl.QueryString)
	*sql += fmt.Sprintf(` $%d`, len(*args))
}

type WhereLikeOperator struct {
	WhereLike
	Operator string
}

type WhereLikes []WhereLikeOperator

func (wls WhereLikes) Set(sql *string, args *[]any, addWhere bool) {
	if len(wls) == 0 {
		return
	}

	if addWhere {
		*sql += ` WHERE`
	}

	for i, wl := range wls {
		wl.Set(sql, args, false)

		if i+1 == len(wls) {
			break
		}

		if wl.Operator != "" {
			*sql += fmt.Sprintf(` %s`, strings.ToUpper(wl.Operator))
		} else {
			*sql += ` AND`
		}
	}
}

func SetWhere(sql *string, args *[]any, query []any, addWhere bool) {
	if addWhere {
		*sql += ` WHERE`
	}

	for _, i := range query {
		switch q := i.(type) {

		case string:
			*sql += fmt.Sprintf(` %s`, q)

		case WhereCondition:
			q.Set(sql, args, false)

		case WhereLike:
			q.Set(sql, args, false)

		case WhereLikes:
			q.Set(sql, args, false)

		default:
			continue
		}
	}
}

type OrderBy struct {
	Exprs string
	Sort  string
	Nulls string
}

func (ob OrderBy) Set(sql *string, addOrderBy bool) {
	if ob.Exprs == "" {
		return
	}

	if addOrderBy {
		*sql += ` ORDER BY`
	}

	*sql += ob.Exprs

	ob.Sort = strings.ToUpper(ob.Sort)
	if ob.Sort == "ASC" || ob.Sort == "DESC" {
		*sql += fmt.Sprintf(` %s`, ob.Sort)
	}

	ob.Nulls = strings.ToUpper(ob.Nulls)
	if ob.Nulls == "FIRST" || ob.Nulls == "LAST" {
		*sql += fmt.Sprintf(` NULLS %s`, ob.Nulls)
	}
}

type OrderBys []OrderBy

func (obs OrderBys) Set(sql *string, addOrderBy bool) {
	if len(obs) == 0 {
		return
	}

	if addOrderBy {
		*sql += ` ORDER BY`
	}

	for i, ob := range obs {
		ob.Set(sql, false)

		if i >= 0 && i != len(obs)-1 {
			*sql += `, `
		}
	}
}

type LimitOffset struct {
	Limit  int
	Offset int
}

func (lo LimitOffset) Set(sql *string, args *[]any) {
	if lo.Limit < 1 {
		return
	}

	*args = append(*args, lo.Limit)
	*sql += fmt.Sprintf(` LIMIT $%d`, len(*args))

	if lo.Offset > 0 {
		*args = append(*args, lo.Offset)
		*sql += fmt.Sprintf(` OFFSET $%d`, len(*args))
	}
}
