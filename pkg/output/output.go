package output

import "github.com/dsreek/sreekdAgent/pkg/model"

type Config struct {
	Format string
}

func Render(issues []model.Issue, cfg Config) {
	if len(issues) == 0 {
		println("No issues found")
		return
	}

	println("Issues:")
	for _, issue := range issues {
		println(issue.ID, "-", issue.Message)
	}
}
