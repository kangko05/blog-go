package post

type Service struct {
	repo Repository
}

func NewService(repo Repository) (*Service, error) {
	if repo == nil {
		memRepo, err := connectSqlite()
		if err != nil {
			return nil, err
		}

		repo = memRepo
	}

	return &Service{repo: repo}, nil
}

func (ps *Service) CreatePost(title, content string) (*Post, error) {
	return createPost(ps.repo, title, content)
}

func (ps *Service) GetPost(id int) (*Post, error) {
	return getPost(ps.repo, id)
}

func (ps *Service) UpdatePost(id int, title, content string) error {
	return updatePost(ps.repo, id, title, content)
}

func (ps *Service) DeletePost(id int) error {
	return deletePost(ps.repo, id)
}

func (ps *Service) ListAllPosts() ([]*Post, error) {
	return listAllPosts(ps.repo)
}
