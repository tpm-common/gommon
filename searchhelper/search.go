package searchhelper

import (
	"errors"
	"strings"
)

type SearchQuery struct {
	Field       *string
	Query       string
	Inverse     *bool
	Insensitive *bool
}

func (q SearchQuery) Validate() error {
	if q.Field != nil && *q.Field == "" {
		return errors.New("search field cannot be empty")
	}

	if q.Query == "" {
		return errors.New("search query cannot be empty")
	}

	return nil
}

type SearchQueryOperator struct {
	SearchQuery
	Operator string
}

func (q SearchQueryOperator) Validate() error {
	if err := q.SearchQuery.Validate(); err != nil {
		return err
	}

	if strings.EqualFold(q.Operator, "AND") || strings.EqualFold(q.Operator, "OR") {
		return errors.New("search operator must be AND or OR")
	}

	return nil
}

type SearchQueries []SearchQueryOperator

func (qs SearchQueries) Validate() error {
	for _, q := range qs {
		if err := q.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type OrderBy struct {
	Field *string
	Sort  *string
	Empty *string
}

func (o OrderBy) Validate() error {
	if o.Field != nil && o.Sort != nil && o.Empty != nil {
		return errors.New("order by must be at least one filled")
	}

	if o.Field != nil && *o.Field == "" {
		return errors.New("order by field cannot be empty")
	}

	if o.Sort != nil {
		if strings.EqualFold(*o.Sort, "ASC") || strings.EqualFold(*o.Sort, "DESC") {
			return errors.New("order by sort must be ASC or DESC")
		}
	}

	if o.Empty != nil {
		if strings.EqualFold(*o.Empty, "FIRST") || strings.EqualFold(*o.Empty, "LAST") {
			return errors.New("order by nils must be FIRST or LAST")
		}
	}

	return nil
}

type OrderBys []OrderBy

func (os OrderBys) Validate() error {
	for _, o := range os {
		if err := o.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type Pagination struct {
	Limit int
	Page  int
}

func (p Pagination) Validate() error {
	if p.Limit < 1 {
		return errors.New("pagination limit must be at least 1")
	}

	if p.Page < 1 {
		return errors.New("pagination page must be at least 1")
	}

	return nil
}

type Search struct {
	SearchQueries SearchQueries
	OrderBys      OrderBys
	Pagination    *Pagination
}

func (s Search) Validate() error {
	if len(s.SearchQueries) != 0 {
		if err := s.SearchQueries.Validate(); err != nil {
			return err
		}
	}

	if len(s.OrderBys) != 0 {
		if err := s.OrderBys.Validate(); err != nil {
			return err
		}
	}

	if s.Pagination != nil {
		if err := s.Pagination.Validate(); err != nil {
			return err
		}
	}

	return nil
}
