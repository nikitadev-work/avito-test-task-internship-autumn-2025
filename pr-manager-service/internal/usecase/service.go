package usecase

type Service struct {
	teams  TeamRepositoryInterface
	users  UserRepositoryInterface
	prs    PullRequestRepositoryInterface
	logger LoggerInterface
}

func NewService(
	teams TeamRepositoryInterface,
	users UserRepositoryInterface,
	prs PullRequestRepositoryInterface,
	logger LoggerInterface,
) *Service {
	return &Service{
		teams:  teams,
		users:  users,
		prs:    prs,
		logger: logger,
	}
}
