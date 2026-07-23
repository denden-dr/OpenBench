package dashboard

import (
	"context"
	"sync"
)

type Service interface {
	GetDashboardData(ctx context.Context) (DashboardResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetDashboardData(ctx context.Context) (DashboardResponse, error) {
	var (
		activeCount   int
		pendingCount  int
		salesToday    float64
		warrantyCount int
		recentTickets []RecentTicket
		errs          [5]error
		wg            sync.WaitGroup
	)

	wg.Add(5)

	go func() {
		defer wg.Done()
		activeCount, errs[0] = s.repo.GetActiveTicketsCount(ctx)
	}()

	go func() {
		defer wg.Done()
		pendingCount, errs[1] = s.repo.GetPendingDiagnosesCount(ctx)
	}()

	go func() {
		defer wg.Done()
		salesToday, errs[2] = s.repo.GetSalesToday(ctx)
	}()

	go func() {
		defer wg.Done()
		warrantyCount, errs[3] = s.repo.GetActiveWarrantiesCount(ctx)
	}()

	go func() {
		defer wg.Done()
		recentTickets, errs[4] = s.repo.GetRecentTickets(ctx)
	}()

	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return DashboardResponse{}, err
		}
	}

	if recentTickets == nil {
		recentTickets = []RecentTicket{}
	}

	return DashboardResponse{
		Metrics: DashboardMetrics{
			ActiveTickets:    activeCount,
			PendingDiagnoses: pendingCount,
			SalesToday:       salesToday,
			ActiveWarranties: warrantyCount,
		},
		RecentTickets: recentTickets,
	}, nil
}
