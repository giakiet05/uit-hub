package conversation

import (
	"fmt"

	"github.com/pkoukk/tiktoken-go"
)

const DefaultMicrocompactMinChars = 12000

var tokenizer *tiktoken.Tiktoken

func init() {
	var err error
	// Use cl100k_base which is a very good generic BPE encoding (used by GPT-3.5/GPT-4).
	tokenizer, err = tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		// Fallback to a dumb estimator if tiktoken fails to load
		tokenizer = nil
	}
}

// EstimateTokens calculates the token size of the conversation array.
func EstimateTokens(messages []Message) int {
	total := 0
	for _, msg := range messages {
		text := Text(msg)
		if tokenizer != nil {
			total += len(tokenizer.Encode(text, nil, nil))
		} else {
			// Fallback heuristic: 4 chars ~ 1 token
			total += len(text) / 4
		}
	}
	return total
}

// CompactionOptions configures the lightweight conversation compaction layers.
type CompactionOptions struct {
	DisableSnip                       bool
	DisableMicrocompact               bool
	MicrocompactMinChars              int
	MicrocompactKeepRecentToolResults int
}

// CompactContext manages the size of the conversation context array to ensure
// it stays within the token thresholds to avoid LLM limit errors.
func CompactContext(messages []Message, maxContextTokens int, options CompactionOptions) ([]Message, string) {
	if maxContextTokens <= 0 {
		return messages, ""
	}

	if options.MicrocompactMinChars <= 0 {
		options.MicrocompactMinChars = DefaultMicrocompactMinChars
	}

	// Calculate thresholds
	// Microcompact fires at 85%
	microThreshold := int(float64(maxContextTokens) * 0.85)
	// Snip fires at 95%
	snipThreshold := int(float64(maxContextTokens) * 0.95)

	currentTokens := EstimateTokens(messages)
	info := ""

	// Layer 1: Snip Compact (Remove old messages)
	if !options.DisableSnip && currentTokens > snipThreshold {
		preSnip := currentTokens
		messages = SnipCompact(messages, snipThreshold)
		currentTokens = EstimateTokens(messages)
		if currentTokens < preSnip {
			info += fmt.Sprintf("[SnipCompact] Removed old messages (%d -> %d tokens).\n", preSnip, currentTokens)
		}
	}

	// Layer 2: Microcompact (Hide Tool Result Data)
	if !options.DisableMicrocompact && currentTokens > microThreshold {
		preMicrocompact := currentTokens
		messages = Microcompact(
			messages,
			microThreshold,
			options.MicrocompactMinChars,
			options.MicrocompactKeepRecentToolResults,
		)
		currentTokens = EstimateTokens(messages)
		if currentTokens < preMicrocompact {
			info += fmt.Sprintf("[Microcompact] Truncated tool results (%d -> %d tokens).\n", preMicrocompact, currentTokens)
		}
	}

	return messages, info
}

// Microcompact iterates from oldest to newest messages and truncates the
// contents of ToolResultMessages that are taking up too much space.
func Microcompact(messages []Message, targetTokens int, minChars int, keepRecentToolResults int) []Message {
	currentTokens := EstimateTokens(messages)
	if currentTokens <= targetTokens {
		return messages
	}
	if minChars <= 0 {
		minChars = DefaultMicrocompactMinChars
	}

	protectedToolResults := recentToolResultIndexes(messages, keepRecentToolResults)
	for i, msg := range messages {
		if _, protected := protectedToolResults[i]; protected {
			continue
		}
		if trm, ok := msg.(ToolResultMessage); ok {
			text := Text(trm)
			if len([]rune(text)) < minChars {
				continue // Skip small results
			}

			// We hide the output but keep the ID
			hiddenMsg := fmt.Sprintf(
				"[Compact: Tool result hidden to save context. Tool call %s executed successfully. Original chars: %d.]",
				trm.ToolCallID,
				len([]rune(text)),
			)
			messages[i] = ToolResultMessage{
				ToolCallID: trm.ToolCallID,
				Content:    NewTextContent(hiddenMsg),
			}

			currentTokens = EstimateTokens(messages)
			if currentTokens <= targetTokens {
				break
			}
		}
	}

	return messages
}

func recentToolResultIndexes(messages []Message, keep int) map[int]struct{} {
	protected := make(map[int]struct{})
	if keep <= 0 {
		return protected
	}

	for index := len(messages) - 1; index >= 0 && keep > 0; index-- {
		if _, ok := messages[index].(ToolResultMessage); !ok {
			continue
		}
		protected[index] = struct{}{}
		keep--
	}
	return protected
}

// SnipCompact removes old user-turn blocks from the beginning of the array
// until the token count drops below target.
func SnipCompact(messages []Message, targetTokens int) []Message {
	currentTokens := EstimateTokens(messages)
	if currentTokens <= targetTokens {
		return messages
	}

	removed := false
	for currentTokens > targetTokens {
		blocks := userTurnBlocks(messages)
		if len(blocks) <= 1 {
			break
		}

		oldest := blocks[0]
		messages = append(messages[:oldest.start], messages[oldest.end:]...)
		removed = true
		currentTokens = EstimateTokens(messages)
	}

	if removed {
		// Prepend a notice message at the beginning
		notice := UserMessage{
			Content: NewTextContent("[System: Older messages have been removed due to context limits]"),
		}
		messages = append([]Message{notice}, messages...)
	}

	return messages
}

type messageBlock struct {
	start int
	end   int
}

func userTurnBlocks(messages []Message) []messageBlock {
	if len(messages) == 0 {
		return nil
	}

	userIndexes := make([]int, 0)
	for index, message := range messages {
		if _, ok := message.(UserMessage); ok {
			userIndexes = append(userIndexes, index)
		}
	}
	if len(userIndexes) == 0 {
		return []messageBlock{{start: 0, end: len(messages)}}
	}

	blocks := make([]messageBlock, 0, len(userIndexes))
	for index, userIndex := range userIndexes {
		start := userIndex
		if index == 0 {
			start = 0
		}

		end := len(messages)
		if index+1 < len(userIndexes) {
			end = userIndexes[index+1]
		}
		blocks = append(blocks, messageBlock{
			start: start,
			end:   end,
		})
	}
	return blocks
}
