package registry

type tagsResponse struct {
	Tags []string `json:"tags"`
}

// Tags returns a list of tags for the given repository.
func (srv *Server) GetTags(repository string) ([]string, error) {
	var response tagsResponse
	var err error

	url := srv.url("/v2/%s/tags/list", repository)
	tags := make([]string, 0)

	for {
		srv.Logger.Debug("registry.tags url=%s repository=%s", url, repository)
		url, err = srv.getPaginatedJSON(url, &response)

		switch err {
		case ErrNoMorePages:
			tags = append(tags, response.Tags...)
			return tags, nil
		case nil:
			tags = append(tags, response.Tags...)
			continue
		default:
			return nil, err
		}
	}
}
