package pull_requests

import (
	"reflect"
	"testing"

	"github.com/google/go-github/github"
)

func TestCheckPRStatus(t *testing.T) {
	// mpatch.PatchMethod(time.Now, func() time.Time {
	// 	return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	// })
	// mpatch.PatchMethod(time.Since, func() time.Duration {
	// 	return time.Duration(0 * time.Minute)
	// })
	// // tenMins := time.Duration(10 * time.Minute)
	// inProgressTime := time.Now().Add(-9 * time.Minute)
	type args struct {
		checks *github.ListCheckRunsResults
	}
	tests := []struct {
		name string
		args args
		want []PRStatus
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
			},
			want: []PRStatus{
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
			},
			want: []PRStatus{
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
			},
			want: []PRStatus{
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
			},
			want: []PRStatus{
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
			},
			want: []PRStatus{
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
			},
			want: []PRStatus{
				{
					Name:    "default check",
					Message: "unaccounted for state conclusion: ",
					Status:  Failure,
					retryIn: 0,
				},
			},
		},
		// {
		// 	name: "check in progress and less than 10 mins old",
		// 	args: args{
		// 		checks: &github.ListCheckRunsResults{
		// 			Total: github.Int(1),
		// 			CheckRuns: []*github.CheckRun{
		// 				{
		// 					Status:    github.String("in_progress"),
		// 					StartedAt: &github.Timestamp{Time: inProgressTime},
		// 					Name:      github.String("in progress check"),
		// 				},
		// 			},
		// 		},
		// 	},
		// 	want: []PRStatus{
		// 		{
		// 			Name:    "in progress check",
		// 			Message: "this check has only just been started check back again in x mins",
		// 			Status:  Pending,
		// 			retryIn: time.Since(time.Now()),
		// 		},
		// 	},
		// },
		// {
		// 	name: "check in progress and more than 10 mins old",
		// },
		// {
		// 	name: "check still in a queued state",
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPRStatus(tt.args.checks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckPRStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
