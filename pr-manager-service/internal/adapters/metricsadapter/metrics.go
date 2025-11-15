package metricsadapter

import (
	"pr-manager-service/internal/usecase"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	teamCreated     prometheus.Counter
	userActivated   prometheus.Counter
	userDeactivated prometheus.Counter
	prCreated       prometheus.Counter
	prMerged        prometheus.Counter
	prReassigned    prometheus.Counter
}

var _ usecase.MetricsInterface = (*Metrics)(nil)

func NewMetrics(serviceName string) *Metrics {
	constNamespace := "pr_manager"

	commonLabels := prometheus.Labels{
		"service": serviceName,
	}

	return &Metrics{
		teamCreated: promauto.NewCounter(prometheus.CounterOpts{
			Namespace:   constNamespace,
			Name:        "teams_created_total",
			Help:        "Total number of created teams",
			ConstLabels: commonLabels,
		}),
		userActivated: promauto.NewCounter(prometheus.CounterOpts{
			Namespace:   constNamespace,
			Name:        "users_activated_total",
			Help:        "Total number of user activations",
			ConstLabels: commonLabels,
		}),
		userDeactivated: promauto.NewCounter(prometheus.CounterOpts{
			Namespace:   constNamespace,
			Name:        "users_deactivated_total",
			Help:        "Total number of user deactivations",
			ConstLabels: commonLabels,
		}),
		prCreated: promauto.NewCounter(prometheus.CounterOpts{
			Namespace:   constNamespace,
			Name:        "pull_requests_created_total",
			Help:        "Total number of created pull requests",
			ConstLabels: commonLabels,
		}),
		prMerged: promauto.NewCounter(prometheus.CounterOpts{
			Namespace:   constNamespace,
			Name:        "pull_requests_merged_total",
			Help:        "Total number of merged pull requests",
			ConstLabels: commonLabels,
		}),
		prReassigned: promauto.NewCounter(prometheus.CounterOpts{
			Namespace:   constNamespace,
			Name:        "pull_requests_reassigned_total",
			Help:        "Total number of reviewer reassignments",
			ConstLabels: commonLabels,
		}),
	}
}

func (m *Metrics) IncTeamCreated() {
	m.teamCreated.Inc()
}

func (m *Metrics) IncUserActivated() {
	m.userActivated.Inc()
}

func (m *Metrics) IncUserDeactivated() {
	m.userDeactivated.Inc()
}

func (m *Metrics) IncPullRequestCreated() {
	m.prCreated.Inc()
}

func (m *Metrics) IncPullRequestMerged() {
	m.prMerged.Inc()
}

func (m *Metrics) IncPullRequestReassigned() {
	m.prReassigned.Inc()
}
