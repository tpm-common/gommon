package searchhelper

type SearchQueryParams struct {
	DefaultField     string
	WhitelistedField []string
	MaxWildcards     int
}

type OrderByParams struct {
	DefaultField     string
	WhitelistedField []string
}

type SearchParams struct {
	SearchQuery SearchQueryParams
	OrderBy     OrderByParams
}
