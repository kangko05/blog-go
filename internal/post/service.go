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

func (ps *Service) CreatePost(cat Category, title, content string, tags []string) (*Post, error) {
	return createPost(ps.repo, cat, title, content, tags)
}

func (ps *Service) GetPost(id int) (*Post, error) {
	return getPost(ps.repo, id)
}

func (ps *Service) UpdatePost(id int, title, content string, tags []string) error {
	return updatePost(ps.repo, id, title, content, tags)
}

func (ps *Service) DeletePost(id int) error {
	return deletePost(ps.repo, id)
}

func (ps *Service) ListAllPosts() ([]*Post, error) {
	return listAllPosts(ps.repo)
}

// NOTE: if dataset gets larger, below should be handled in database
func (ps *Service) ListCategory(cat Category) ([]*Post, error) {
	posts, err := ps.ListAllPosts()
	if err != nil {
		return nil, err
	}

	if cat == ALL {
		return posts, nil
	}

	var result []*Post

	for _, post := range posts {
		if post.Category == cat {
			result = append(result, post)
		}
	}

	return result, nil
}

func (ps *Service) ListNotes() ([]*Post, error) {
	return ps.ListCategory(NOTES)
}

func (ps *Service) ListProjects() ([]*Post, error) {
	return ps.ListCategory(PROJECTS)
}
