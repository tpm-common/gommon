package searchhelper

import (
	"fmt"
	"strings"

	"github.com/tpm-common/gommon/sqlhelper"
)

func (q *SearchQuery) SQLHelper(params SearchQueryParams) *sqlhelper.WhereLike {
	if q == nil {
		return nil
	}

	wl := sqlhelper.WhereLike{}

	if q.Field != nil && *q.Field != "" {
		if len(params.WhitelistedField) != 0 {
			wFound := false
			for _, w := range params.WhitelistedField {
				if *q.Field == w {
					wFound = true
					break
				}
			}
			if !wFound {
				return nil
			}
		}

		wl.Column = *q.Field
	} else {
		wl.Column = params.DefaultField
	}

	if q.Query == "" {
		return nil
	}

	qString := q.Query

	var (
		numWildcard int
		p1, e1, e2  string
	)

	c1 := strings.Count(qString, "*")
	if eString := `\*`; strings.Contains(qString, eString) {
		c1 = c1 - strings.Count(qString, eString)
		p1 = fmt.Sprintf("-%s-", replaceString(5))
		qString = strings.ReplaceAll(qString, eString, p1)
	}
	numWildcard += c1

	c2 := strings.Count(qString, "_")
	if eString := `\_`; strings.Contains(qString, eString) {
		c2 = c2 - strings.Count(qString, eString)
		e1 = fmt.Sprintf("-%s-", replaceString(5))
		qString = strings.ReplaceAll(qString, eString, e1)
	}
	numWildcard += c2

	if numWildcard > params.MaxWildcards {
		return nil
	}

	if eString := `\\`; strings.Contains(qString, eString) {
		e2 = fmt.Sprintf("-%s-", replaceString(5))
		qString = strings.ReplaceAll(qString, eString, e2)
	}

	qString = strings.ReplaceAll(qString, `\`, `\\`)
	qString = strings.ReplaceAll(qString, "%", `\%`)

	qString = strings.ReplaceAll(qString, "*", "%")
	if p1 != "" {
		qString = strings.ReplaceAll(qString, p1, `*`)
	}
	if strings.HasSuffix(qString, `\\*`) {
		qString = strings.TrimSuffix(qString, "*") + "%"
	}

	if e1 != "" {
		qString = strings.ReplaceAll(qString, e1, `\_`)
	}
	if e2 != "" {
		qString = strings.ReplaceAll(qString, e2, `\\`)
	}

	wl.QueryString = qString

	if q.Inverse != nil {
		wl.Inverse = *q.Inverse
	}

	if q.Insensitive != nil {
		wl.Insensitive = *q.Insensitive
	}

	return &wl
}

func (qs SearchQueries) SQLHelper(params SearchQueryParams) sqlhelper.WhereLikes {
	var wls sqlhelper.WhereLikes

	for _, q := range qs {
		wl := q.SQLHelper(params)
		if wl == nil {
			continue
		}

		wlo := sqlhelper.WhereLikeOperator{
			WhereLike: *wl,
			Operator:  q.Operator,
		}

		wls = append(wls, wlo)
	}

	return wls
}

func (o *OrderBy) SQLHelper(params OrderByParams) *sqlhelper.OrderBy {
	if o == nil {
		return nil
	}

	ob := sqlhelper.OrderBy{}

	switch {
	case o.Field != nil && *o.Field != "":
		if len(params.WhitelistedField) != 0 {
			wFound := false
			for _, w := range params.WhitelistedField {
				if ob.Exprs == w {
					wFound = true
					break
				}
			}
			if !wFound {
				return nil
			}
		}
		ob.Exprs = *o.Field
	case (o.Sort != nil && *o.Sort != "") || (o.Empty != nil && *o.Empty != ""):
		ob.Exprs = params.DefaultField
	}

	if o.Sort != nil && *o.Sort != "" {
		ob.Sort = *o.Sort
	}
	if o.Empty != nil && *o.Empty != "" {
		ob.Nulls = *o.Empty
	}

	return &ob
}

func (os OrderBys) SQLHelper(params OrderByParams) sqlhelper.OrderBys {
	var obs sqlhelper.OrderBys

	for _, o := range os {
		ob := o.SQLHelper(params)
		if ob == nil {
			continue
		}

		for _, c := range obs {
			if ob.Exprs == c.Exprs {
				continue
			}
		}

		obs = append(obs, *ob)
	}

	return obs
}

func (p *Pagination) SQLHelper() *sqlhelper.LimitOffset {
	if p == nil {
		return nil
	}

	return &sqlhelper.LimitOffset{
		Limit:  p.Limit * p.Page,
		Offset: p.Limit * (p.Page - 1),
	}
}

func (s *Search) SQLHelper(params SearchParams) *sqlhelper.SearchParams {
	if s == nil {
		return nil
	}

	return &sqlhelper.SearchParams{
		WhereLikes:  s.SearchQueries.SQLHelper(params.SearchQuery),
		OrderBys:    s.OrderBys.SQLHelper(params.OrderBy),
		LimitOffset: s.Pagination.SQLHelper(),
	}
}
