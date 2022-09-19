package sqlhelper

type SearchParams struct {
	WhereLikes  WhereLikes
	OrderBys    OrderBys
	LimitOffset *LimitOffset

	DefaultOrderBy *OrderBy
	MaxResults     int
}

func (s SearchParams) Set(sql *string, args *[]any) {
	s.WhereLikes.Set(sql, args, true)

	if s.DefaultOrderBy != nil && len(s.OrderBys) == 0 {
		s.OrderBys = append(s.OrderBys, *s.DefaultOrderBy)
	}
	s.OrderBys.Set(sql, true)

	if s.LimitOffset != nil && s.LimitOffset.Limit <= s.MaxResults {
		s.LimitOffset = &LimitOffset{Limit: s.MaxResults, Offset: 0}
	}
	s.LimitOffset.Set(sql, args)
}
