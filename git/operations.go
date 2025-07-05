package git

func DefaultGitOptions() *GitOptions {
	return &GitOptions{}
}

type GitOptions struct {
	IdentityFile string
}
