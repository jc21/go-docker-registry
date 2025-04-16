package registry

type repositoriesResponse struct {
	Repositories []string `json:"repositories"`
}

// Tags returns a list of tags for the given repository.
func (srv *Server) GetCatalog() ([]string, error) {
	var response repositoriesResponse
	var err error

	url := srv.url("/v2/_catalog")
	repos := make([]string, 0)

	for {
		srv.Logger.Debug("registry.catalog url=%s", url)
		url, err = srv.getPaginatedJSON(url, &response)
		switch err {
		case ErrNoMorePages:
			repos = append(repos, response.Repositories...)
			return repos, nil
		case nil:
			repos = append(repos, response.Repositories...)
			continue
		default:
			return nil, err
		}
	}
}
