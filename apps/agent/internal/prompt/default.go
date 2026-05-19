package prompt

const defaultSessionPrompt = "You are a helpful agent."

// DefaultSystemPrompt returns the application-level prompt used by the agent
// runtime when no external prompt source is configured.
func DefaultSystemPrompt() SystemPrompt {
	return SystemPrompt{
		SessionText: defaultSessionPrompt,
	}
}
