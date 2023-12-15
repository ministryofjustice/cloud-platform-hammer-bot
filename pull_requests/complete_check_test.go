package pull_requests

import (
	"reflect"
	"testing"

	"github.com/google/go-github/v57/github"
)

func TestCompletedCheck(t *testing.T) {
	tests := []struct {
		name     string
		check    *github.CheckRun
		prStatus []InvalidChecks
		want     []InvalidChecks
	}{
		{
			name: "success",
			check: &github.CheckRun{
				Conclusion: github.String("success"),
			},
			prStatus: nil,
			want:     nil,
		},
		{
			name: "skipped",
			check: &github.CheckRun{
				Conclusion: github.String("skipped"),
			},
			prStatus: nil,
			want:     nil,
		},
		{
			name: "failure",
			check: &github.CheckRun{
				Conclusion: github.String("failure"),
				Name:       github.String("failed check"),
			},
			prStatus: []InvalidChecks{},
			want: []InvalidChecks{
				{
					Name:           "failed check",
					Message:        "this check failed, check your pr and ammend",
					Status:         Failure,
					RetryInNanoSec: 0,
				},
			},
		},
		{
			name: "action_required",
			check: &github.CheckRun{
				Conclusion: github.String("action_required"),
				Name:       github.String("action required check"),
			},
			prStatus: []InvalidChecks{},
			want: []InvalidChecks{
				{
					Name:           "action required check",
					Message:        "this check failed because an action is required, check your pr and ammend",
					Status:         Failure,
					RetryInNanoSec: 0,
				},
			},
		},
		{
			name: "cancelled",
			check: &github.CheckRun{
				Conclusion: github.String("cancelled"),
				Name:       github.String("cancelled check"),
			},
			prStatus: []InvalidChecks{},
			want: []InvalidChecks{
				{
					Name:           "cancelled check",
					Message:        "this check failed because somebody manually cancelled the check",
					Status:         Failure,
					RetryInNanoSec: 0,
				},
			},
		},
		{
			name: "timed_out",
			check: &github.CheckRun{
				Conclusion: github.String("timed_out"),
				Name:       github.String("timed out check"),
			},
			prStatus: []InvalidChecks{},
			want: []InvalidChecks{
				{
					Name:           "timed out check",
					Message:        "this check failed because it timed out",
					Status:         Failure,
					RetryInNanoSec: 0,
				},
			},
		},
		{
			name: "stale",
			check: &github.CheckRun{
				Conclusion: github.String("stale"),
				Name:       github.String("stale check"),
			},
			prStatus: []InvalidChecks{},
			want: []InvalidChecks{
				{
					Name:           "stale check",
					Message:        "this check failed because it was stale",
					Status:         Failure,
					RetryInNanoSec: 0,
				},
			},
		},
		{
			name: "unaccounted for state conclusion",
			check: &github.CheckRun{
				Conclusion: github.String("unknown"),
				Name:       github.String("unknown check"),
			},
			prStatus: []InvalidChecks{},
			want: []InvalidChecks{
				{
					Name:           "unknown check",
					Message:        "unaccounted for state conclusion: unknown",
					Status:         Failure,
					RetryInNanoSec: 0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompletedCheck(tt.check, tt.prStatus); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompletedCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
