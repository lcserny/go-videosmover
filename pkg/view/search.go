package view

func Search() (tmplName string, tmplData interface{}) {
	return "search.gohtml", nil
}

func SearchResults() (tmplName string, tmplData interface{}) {
	return "search_results.gohtml", nil
}
