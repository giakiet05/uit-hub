package localtool

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// KnowledgeLookup returns fake private knowledge that the model should not know
// without calling a tool.
type KnowledgeLookup struct {
	records []knowledgeRecord
}

// knowledgeRecord is one searchable fake knowledge entry.
type knowledgeRecord struct {
	Key     string
	Answer  string
	Aliases []string
}

// NewKnowledgeLookup creates the fake knowledge lookup tool with built-in test
// records.
func NewKnowledgeLookup() *KnowledgeLookup {
	return &KnowledgeLookup{
		records: []knowledgeRecord{
			{
				Key:    "uit_hub_secret_code",
				Answer: "UIT-HUB-2026-ALPHA",
				Aliases: []string{
					"uit hub secret code",
					"secret code",
					"internal code",
					"private code",
					"ma bi mat noi bo uit hub",
					"ma bi mat uit hub",
					"ma noi bo uit hub",
					"ma bi mat",
					"ma noi bo",
				},
			},
			{
				Key:    "course_registration_rule",
				Answer: "Students may register for at most 24 credits in one semester unless an advisor override is approved.",
				Aliases: []string{
					"course registration rule",
					"credit limit",
					"maximum credits",
					"quy dinh dang ky mon hoc",
					"gioi han tin chi",
				},
			},
			{
				Key:    "academic_warning_threshold",
				Answer: "A student is placed under academic warning when semester GPA is below 2.0 or accumulated GPA is below 2.0.",
				Aliases: []string{
					"academic warning threshold",
					"academic warning",
					"gpa warning",
					"canh bao hoc vu",
					"nguong canh bao hoc vu",
				},
			},
		},
	}
}

// Definition describes the fake knowledge lookup schema.
func (t *KnowledgeLookup) Definition() tool.Definition {
	return tool.Definition{
		Name:        "knowledge_lookup",
		Description: "Look up private UIT Hub knowledge that is only available through this tool.",
		InputSchema: map[string]any{
			"type":                 "object",
			"additionalProperties": false,
			"properties": map[string]any{
				"query": map[string]any{
					"type":        "string",
					"description": "Search query for private knowledge.",
				},
			},
			"required": []string{"query"},
		},
	}
}

// Execute searches the private knowledge records and returns a JSON result.
func (t *KnowledgeLookup) Execute(ctx context.Context, call tool.Call) (tool.Result, error) {
	if err := ctx.Err(); err != nil {
		return tool.Result{}, err
	}

	query, err := stringArg(call.Arguments, "query")
	if err != nil {
		return tool.Result{}, err
	}

	key, answer, found := t.lookup(query)
	output, err := json.Marshal(map[string]any{
		"found":  found,
		"key":    key,
		"answer": answer,
	})
	if err != nil {
		return tool.Result{}, fmt.Errorf("marshal result: %w", err)
	}

	return tool.Result{
		CallID:  call.ID,
		Name:    call.Name,
		Content: string(output),
	}, nil
}

// lookup searches exact aliases first, then falls back to longer token matches.
func (t *KnowledgeLookup) lookup(query string) (string, string, bool) {
	normalizedQuery := normalizeLookupText(query)
	for _, record := range t.records {
		candidates := append([]string{record.Key, record.Answer}, record.Aliases...)
		for _, candidate := range candidates {
			if strings.Contains(normalizedQuery, normalizeLookupText(candidate)) {
				return record.Key, record.Answer, true
			}
		}
	}

	for _, record := range t.records {
		candidates := append([]string{record.Key, record.Answer}, record.Aliases...)
		for _, token := range strings.Fields(normalizeLookupText(strings.Join(candidates, " "))) {
			if len(token) >= 5 && strings.Contains(normalizedQuery, token) {
				return record.Key, record.Answer, true
			}
		}
	}
	return "", "", false
}

// normalizeLookupText lowercases, removes simple Vietnamese accents, and
// normalizes separators for fuzzy matching.
func normalizeLookupText(text string) string {
	text = strings.ToLower(text)
	replacer := strings.NewReplacer(
		"_", " ",
		"-", " ",
		".", " ",
		",", " ",
		":", " ",
		"á", "a",
		"à", "a",
		"ả", "a",
		"ã", "a",
		"ạ", "a",
		"ă", "a",
		"ắ", "a",
		"ằ", "a",
		"ẳ", "a",
		"ẵ", "a",
		"ặ", "a",
		"â", "a",
		"ấ", "a",
		"ầ", "a",
		"ẩ", "a",
		"ẫ", "a",
		"ậ", "a",
		"đ", "d",
		"é", "e",
		"è", "e",
		"ẻ", "e",
		"ẽ", "e",
		"ẹ", "e",
		"ê", "e",
		"ế", "e",
		"ề", "e",
		"ể", "e",
		"ễ", "e",
		"ệ", "e",
		"í", "i",
		"ì", "i",
		"ỉ", "i",
		"ĩ", "i",
		"ị", "i",
		"ó", "o",
		"ò", "o",
		"ỏ", "o",
		"õ", "o",
		"ọ", "o",
		"ô", "o",
		"ố", "o",
		"ồ", "o",
		"ổ", "o",
		"ỗ", "o",
		"ộ", "o",
		"ơ", "o",
		"ớ", "o",
		"ờ", "o",
		"ở", "o",
		"ỡ", "o",
		"ợ", "o",
		"ú", "u",
		"ù", "u",
		"ủ", "u",
		"ũ", "u",
		"ụ", "u",
		"ư", "u",
		"ứ", "u",
		"ừ", "u",
		"ử", "u",
		"ữ", "u",
		"ự", "u",
		"ý", "y",
		"ỳ", "y",
		"ỷ", "y",
		"ỹ", "y",
		"ỵ", "y",
	)
	return strings.Join(strings.Fields(replacer.Replace(text)), " ")
}
