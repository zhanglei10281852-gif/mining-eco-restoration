package export_test

import (
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/export"
	"strings"
	"testing"
)

func TestBuildSummary(t *testing.T) {
	p := domain.Project{ID: "p", Name: "Mine", Region: "R"}
	tasks := []domain.RemediationTask{{Status: "accepted"}, {Status: "submitted"}}
	samples := []domain.Sample{{ID: "s"}}
	s := export.Build(p, tasks, samples)
	if s.TaskCount != 2 || s.AcceptedCount != 1 || s.SampleCount != 1 {
		t.Fatal(s)
	}
}
func TestSortAndJSON(t *testing.T) {
	in := []domain.RemediationTask{{ID: "1", Status: "submitted"}, {ID: "2", Status: "accepted"}, {ID: "3", Status: "planned"}}
	out := export.SortByStatus(in)
	if out[0].Status != "accepted" || in[0].ID != "1" {
		t.Fatal(out)
	}
	text, e := export.JSON(out)
	if e != nil || !strings.Contains(text, "accepted") {
		t.Fatal(text, e)
	}
}
