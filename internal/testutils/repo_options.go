package testutils

type RepoOption func(*GitTestRepo)

func WithDefaultBranch(branch string) RepoOption {
	return func(r *GitTestRepo) {
		r.Config.DefaultBranch = branch
	}
}

func WithUser(name, email string) RepoOption {
	return func(r *GitTestRepo) {
		r.Config.UserName = name
		r.Config.UserEmail = email
	}
}

func WithRemoteName(name string) RepoOption {
	return func(r *GitTestRepo) {
		r.Config.RemoteName = name
	}
}
