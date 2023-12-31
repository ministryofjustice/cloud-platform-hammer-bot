package pull_requests

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-github/v57/github"
)

// mock time.Now() with fixed time
func mockTimeNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestInProgressCheck(t *testing.T) {
	n := mockTimeNow(time.Now())
	tests := []struct {
		name           string
		check          *github.CheckRun
		prStatus       []InvalidChecks
		getTimeSince   func(time.Time) time.Duration
		want           []InvalidChecks
		timeSinceStart time.Duration
	}{
		{
			name: "check started more than 10 minutes ago and is in progress",
			check: &github.CheckRun{
				Name:       github.String("test check in progress"),
				StartedAt:  &github.Timestamp{Time: n().Add(-15 * time.Minute)},
				Conclusion: github.String("neutral"),
			},
			prStatus:     nil,
			getTimeSince: func(t time.Time) time.Duration { return time.Since(t) },
			want: []InvalidChecks{
				{
					Name:           "test check in progress",
					Message:        "this check has been in_progress for at least 10 mins, looks like something has gone wrong?",
					Status:         Pending,
					RetryInNanoSec: 0,
				},
			},
		},
		{
			name: "check started exactly 10 minutes ago and is in progress",
			check: &github.CheckRun{
				Name:       github.String("test check in progress"),
				StartedAt:  &github.Timestamp{Time: n().Add(-10 * time.Minute)},
				Conclusion: github.String("neutral"),
			},
			prStatus:     nil,
			getTimeSince: func(t time.Time) time.Duration { return time.Since(t) },
			want: []InvalidChecks{
				{
					Name:           "test check in progress",
					Message:        "this check has been in_progress for at least 10 mins, looks like something has gone wrong?",
					Status:         Pending,
					RetryInNanoSec: 0,
				},
			},
		},
		{
			name: "check started less than 10 minutes ago and is in progress",
			check: &github.CheckRun{
				Name:       github.String("test check in progress"),
				StartedAt:  &github.Timestamp{Time: n().Add(-5 * time.Minute)},
				Conclusion: github.String("neutral"),
			},
			prStatus:     nil,
			getTimeSince: func(t time.Time) time.Duration { return time.Since(t) },
			want: []InvalidChecks{
				{
					Name:           "test check in progress",
					Message:        "this check is in_progress and has just been started. check back again in " + (10*time.Minute - 5*time.Minute).String(),
					Status:         Pending,
					RetryInNanoSec: (10*time.Minute - 5*time.Minute),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InProgressCheck(tt.check, tt.prStatus, tt.getTimeSince); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InProgressCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueuedCheck(t *testing.T) {
	n := mockTimeNow(time.Now())
	tests := []struct {
		name           string
		check          *github.CheckRun
		prStatus       []InvalidChecks
		getTimeSince   func(time.Time) time.Duration
		want           []InvalidChecks
		timeSinceStart time.Duration
	}{
		{
			name: "check started more than 10 minutes ago and is queued",
			check: &github.CheckRun{
				Name:       github.String("test check queued"),
				StartedAt:  &github.Timestamp{Time: n().Add(-15 * time.Minute)},
				Conclusion: github.String("neutral"),
			},
			prStatus:     nil,
			getTimeSince: func(t time.Time) time.Duration { return time.Since(t) },
			want: []InvalidChecks{
				{
					Name:           "test check queued",
					Message:        "this check has been queued for at least 10 mins, looks like something has gone wrong?",
					Status:         Pending,
					RetryInNanoSec: 0,
				},
			},
		},
		{
			name: "check started exactly 10 minutes ago and is queued",
			check: &github.CheckRun{
				Name:       github.String("test check queued"),
				StartedAt:  &github.Timestamp{Time: n().Add(-10 * time.Minute)},
				Conclusion: github.String("neutral"),
			},
			prStatus:     nil,
			getTimeSince: func(t time.Time) time.Duration { return time.Since(t) },
			want: []InvalidChecks{
				{
					Name:           "test check queued",
					Message:        "this check has been queued for at least 10 mins, looks like something has gone wrong?",
					Status:         Pending,
					RetryInNanoSec: 0,
				},
			},
		},
		{
			name: "check started less than 10 minutes ago and is queued",
			check: &github.CheckRun{
				Name:       github.String("test check queued"),
				StartedAt:  &github.Timestamp{Time: n().Add(-5 * time.Minute)},
				Conclusion: github.String("neutral"),
			},
			prStatus:     nil,
			getTimeSince: func(t time.Time) time.Duration { return time.Since(t) },
			want: []InvalidChecks{
				{
					Name:           "test check queued",
					Message:        "this check has been queued for less than 10 minutes, check back again in " + (10*time.Minute - 5*time.Minute).String(),
					Status:         Pending,
					RetryInNanoSec: (10*time.Minute - 5*time.Minute),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := QueuedCheck(tt.check, tt.prStatus, tt.getTimeSince); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueuedCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
