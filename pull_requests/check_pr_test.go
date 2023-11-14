package pull_requests

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-github/github"
)

func wrapTimeSince(mins int64) func(time.Time) time.Duration {
	return func(time.Time) time.Duration {
		return time.Duration(mins) * time.Minute
	}
}

func TestCheckInvalidChecks(t *testing.T) {
	tenMins := time.Duration(10 * time.Minute)
	inProgressTime := time.Now().Add(-9 * time.Minute)
	mockRetryInShort := tenMins - wrapTimeSince(9)(inProgressTime)

	type args struct {
		checks       *github.ListCheckRunsResults
		getTimeSince func(time.Time) time.Duration
	}
	tests := []struct {
		name string
		args args
		want []InvalidChecks
	}{
		{
			name: "check is completed and success",
			args: args{
				checks: &github.ListCheckRunsResults{
					Total: github.Int(1),
					CheckRuns: []*github.CheckRun{
						{
							Status:     github.String("completed"),
							Conclusion: github.String("success"),
						},
					},
				},
				getTimeSince: time.Since,
			},
			want: nil,
		},
		{
			name: "check is completed and skipped",
			args: args{
				checks: &github.ListCheckRunsResults{
					Total: github.Int(1),
					CheckRuns: []*github.CheckRun{
						{
							Status:     github.String("completed"),
							Conclusion: github.String("skipped"),
						},
					},
				},
				getTimeSince: time.Since,
			},
			want: nil,
		},
		{
			name: "check is completed and failed",
			args: args{
				checks: &github.ListCheckRunsResults{
					Total: github.Int(1),
					CheckRuns: []*github.CheckRun{
						{
							Status:     github.String("completed"),
							Conclusion: github.String("failure"),
							Name:       github.String("failed check"),
						},
					},
				},
				getTimeSince: time.Since,
			},
			want: []InvalidChecks{
				{
					Name:    "failed check",
					Message: "this check failed, check your pr and ammend",
					Status:  Failure,
					retryIn: 0,
				},
			},
		},
		{
			name: "check is completed and action required",
			args: args{
				checks: &github.ListCheckRunsResults{
					Total: github.Int(1),
					CheckRuns: []*github.CheckRun{
						{
							Status:     github.String("completed"),
							Conclusion: github.String("action_required"),
							Name:       github.String("action required check"),
						},
					},
				},
				getTimeSince: time.Since,
			},
			want: []InvalidChecks{
				{
					Name:    "action required check",
					Message: "this check failed because an action is required, check your pr and ammend",
					Status:  Failure,
					retryIn: 0,
				},
			},
		},
		{
			name: "check is completed and cancelled",
			args: args{
				checks: &github.ListCheckRunsResults{
					Total: github.Int(1),
					CheckRuns: []*github.CheckRun{
						{
							Status:     github.String("completed"),
							Conclusion: github.String("cancelled"),
							Name:       github.String("cancelled check"),
						},
					},
				},
				getTimeSince: time.Since,
			},
			want: []InvalidChecks{
				{
					Name:    "cancelled check",
					Message: "this check failed because somebody manually cancelled the check",
					Status:  Failure,
					retryIn: 0,
				},
			},
		},
		{
			name: "check is completed and timed out",
			args: args{
				checks: &github.ListCheckRunsResults{
					Total: github.Int(1),
					CheckRuns: []*github.CheckRun{
						{
							Status:     github.String("completed"),
							Conclusion: github.String("timed_out"),
							Name:       github.String("timed out check"),
						},
					},
				},
				getTimeSince: time.Since,
			},
			want: []InvalidChecks{
				{
					Name:    "timed out check",
					Message: "this check failed because it timed out",
					Status:  Failure,
					retryIn: 0,
				},
			},
		},
		{
			name: "check is completed and stale",
			args: args{
				checks: &github.ListCheckRunsResults{
					Total: github.Int(1),
					CheckRuns: []*github.CheckRun{
						{
							Status:     github.String("completed"),
							Conclusion: github.String("stale"),
							Name:       github.String("stale check"),
						},
					},
				},
				getTimeSince: time.Since,
			},
			want: []InvalidChecks{
				{
					Name:    "stale check",
					Message: "this check failed because it was stale",
					Status:  Failure,
					retryIn: 0,
				},
			},
		},
		{
			name: "default case",
			args: args{
				checks: &github.ListCheckRunsResults{
					Total: github.Int(1),
					CheckRuns: []*github.CheckRun{
						{
							Status:     github.String("completed"),
							Conclusion: github.String(""),
							Name:       github.String("default check"),
						},
					},
				},
				getTimeSince: time.Since,
			},
			want: []InvalidChecks{
				{
					Name:    "default check",
					Message: "unaccounted for state conclusion: ",
					Status:  Failure,
					retryIn: 0,
				},
			},
		},
		{
			name: "check in progress and LESS than 10 mins old",
			args: args{
				checks: &github.ListCheckRunsResults{
					Total: github.Int(1),
					CheckRuns: []*github.CheckRun{
						{
							Status:    github.String("in_progress"),
							StartedAt: &github.Timestamp{Time: inProgressTime},
							Name:      github.String("in progress short running check"),
						},
					},
				},
				getTimeSince: wrapTimeSince(9),
			},
			want: []InvalidChecks{
				{
					Name:    "in progress short running check",
					Message: "this check is in_progress and has just been started. check back again in " + mockRetryInShort.String(),
					Status:  Pending,
					retryIn: mockRetryInShort,
				},
			},
		},
		{
			name: "check in progress and MORE than 10 mins old",
			args: args{
				checks: &github.ListCheckRunsResults{
					Total: github.Int(1),
					CheckRuns: []*github.CheckRun{
						{
							Status:    github.String("in_progress"),
							StartedAt: &github.Timestamp{Time: inProgressTime},
							Name:      github.String("in progress long running check"),
						},
					},
				},
				getTimeSince: wrapTimeSince(20),
			},
			want: []InvalidChecks{
				{
					Name:    "in progress long running check",
					Message: "this check has been in_progress for at least 10 mins, looks like something has gone wrong?",
					Status:  Pending,
					retryIn: 0,
				},
			},
		},
		{
			name: "check queued and LESS than 10 mins old",
			args: args{
				checks: &github.ListCheckRunsResults{
					Total: github.Int(1),
					CheckRuns: []*github.CheckRun{
						{
							Status:    github.String("queued"),
							StartedAt: &github.Timestamp{Time: inProgressTime},
							Name:      github.String("queued short running check"),
						},
					},
				},
				getTimeSince: wrapTimeSince(9),
			},
			want: []InvalidChecks{
				{
					Name:    "queued short running check",
					Message: "this check has been queued for less than 10 minutes, check back again in " + mockRetryInShort.String(),
					Status:  Pending,
					retryIn: mockRetryInShort,
				},
			},
		},
		{
			name: "check queued and MORE than 10 mins old",
			args: args{
				checks: &github.ListCheckRunsResults{
					Total: github.Int(1),
					CheckRuns: []*github.CheckRun{
						{
							Status:    github.String("queued"),
							StartedAt: &github.Timestamp{Time: inProgressTime},
							Name:      github.String("queued long running check"),
						},
					},
				},
				getTimeSince: wrapTimeSince(20),
			},
			want: []InvalidChecks{
				{
					Name:    "queued long running check",
					Message: "this check has been queued for at least 10 mins, looks like something has gone wrong?",
					Status:  Pending,
					retryIn: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPRStatus(tt.args.checks, tt.args.getTimeSince); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckPRStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
