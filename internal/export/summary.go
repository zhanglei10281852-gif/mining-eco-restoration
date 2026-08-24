package export

import (
	"encoding/json"
	"github.com/zhanglei10281852-gif/mining-eco-restoration/internal/domain"
	"sort"
	"time"
)

type ProjectSummary struct {
	ProjectID, ProjectName, Region        string
	TaskCount, AcceptedCount, SampleCount int
	GeneratedAt                           time.Time
}

func Build(p domain.Project, tasks []domain.RemediationTask, samples []domain.Sample) ProjectSummary {
	a := 0
	for _, t := range tasks {
		if t.Status == "accepted" {
			a++
		}
	}
	return ProjectSummary{ProjectID: p.ID, ProjectName: p.Name, Region: p.Region, TaskCount: len(tasks), AcceptedCount: a, SampleCount: len(samples), GeneratedAt: time.Now()}
}
func SortByStatus(tasks []domain.RemediationTask) []domain.RemediationTask {
	out := append([]domain.RemediationTask(nil), tasks...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].Status < out[j].Status })
	return out
}
func JSON(v any) (string, error) { b, e := json.Marshal(v); return string(b), e }
