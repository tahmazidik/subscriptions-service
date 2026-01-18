package service

import (
	"context"
	"time"
)

func (s *Service) Total(ctx context.Context, userID, serviceName string, periodStart, periodEnd time.Time) (int, error) {
	subs, err := s.repo.ListForPeriod(ctx, userID, serviceName, periodStart, periodEnd)
	if err != nil {
		return 0, err
	}

	total := 0
	for _, sub := range subs {
		subEnd := periodEnd
		if sub.EndDate != nil {
			subEnd = *sub.EndDate
		}

		overlapStart := maxMonth(periodStart, sub.StartDate)
		overlapEnd := minMonth(periodEnd, subEnd)

		if monthIndex(overlapStart) > monthIndex(overlapEnd) {
			continue
		}

		months := monthIndex(overlapEnd) - monthIndex(overlapStart) + 1
		total += months * sub.Price
	}

	return total, nil
}

func monthIndex(t time.Time) int {
	return t.Year()*12 + int(t.Month())
}

func maxMonth(a, b time.Time) time.Time {
	if monthIndex(a) >= monthIndex(b) {
		return a
	}
	return b
}

func minMonth(a, b time.Time) time.Time {
	if monthIndex(a) <= monthIndex(b) {
		return a
	}
	return b
}
