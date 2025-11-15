package usecase

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
