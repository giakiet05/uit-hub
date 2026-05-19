package app

import (
	"github.com/giakiet05/uit-hub/apps/agent/internal/localtool"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// newToolRegistry registers the temporary local tools used by the current
// baseline agent.
func newToolRegistry() (*tool.Registry, error) {
	return tool.NewRegistry(
		localtool.NewEcho(),
		localtool.NewCalculator(),
		localtool.NewCurrentTime(),
		localtool.NewWriteFile("tmp/agent-files"),
		localtool.NewKnowledgeLookup(),
		localtool.NewStudentProfileLookup(),
	)
}
