package usecase

// Service contains business logic for teams, users and pull requests
type Service struct {
	teams   TeamRepositoryInterface
	users   UserRepositoryInterface
	prs     PullRequestRepositoryInterface
	logger  LoggerInterface
	metrics MetricsInterface
}

func NewService(
	teams TeamRepositoryInterface,
	users UserRepositoryInterface,
	prs PullRequestRepositoryInterface,
	logger LoggerInterface,
	metrics MetricsInterface,
) *Service {
	return &Service{
		teams:   teams,
		users:   users,
		prs:     prs,
		logger:  logger,
		metrics: metrics,
	}
}
