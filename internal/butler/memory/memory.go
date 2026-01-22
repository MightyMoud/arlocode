package memory

// This memory shape is made to match the ATIF format
// See: https://github.com/laude-institute/harbor/blob/main/docs/rfcs/0001-trajectory-format.md
// This is meant to follow standards and will be later stored in db for tracking and analysis

type MemoryEntry struct {
	StepID           int                    `json:"step_id"`
	Timestamp        string                 `json:"timestamp,omitempty"`
	Source           string                 `json:"source"`
	ModelName        string                 `json:"model_name,omitempty"`
	ReasoningEffort  interface{}            `json:"reasoning_effort,omitempty"`
	Message          string                 `json:"message"`
	ReasoningContent string                 `json:"reasoning_content,omitempty"`
	ToolCalls        []ToolCall             `json:"tool_calls,omitempty"`
	Observation      Observation            `json:"observation,omitempty"`
	Metrics          Metrics                `json:"metrics,omitempty"`
	Extra            map[string]interface{} `json:"extra,omitempty"`
}

type ToolCall struct {
	ToolCallID   string                 `json:"tool_call_id"`
	FunctionName string                 `json:"function_name"`
	Arguments    map[string]interface{} `json:"arguments"`
}

type Observation struct {
	Results []ObservationResult `json:"results,omitempty"`
}

type ObservationResult struct {
	SourceCallID          string                  `json:"source_call_id,omitempty"`
	Content               string                  `json:"content,omitempty"`
	SubagentTrajectoryRef []SubagentTrajectoryRef `json:"subagent_trajectory_ref,omitempty"`
}

type SubagentTrajectoryRef struct {
	SessionID      string                 `json:"session_id"`
	TrajectoryPath string                 `json:"trajectory_path,omitempty"`
	Extra          map[string]interface{} `json:"extra,omitempty"`
}

type Metrics struct {
	PromptTokens       int                    `json:"prompt_tokens,omitempty"`
	CompletionTokens   int                    `json:"completion_tokens,omitempty"`
	CachedTokens       int                    `json:"cached_tokens,omitempty"`
	CostUSD            float64                `json:"cost_usd,omitempty"`
	PromptTokenIDs     []int                  `json:"prompt_token_ids,omitempty"`
	CompletionTokenIDs []int                  `json:"completion_token_ids,omitempty"`
	Logprobs           []float64              `json:"logprobs,omitempty"`
	Extra              map[string]interface{} `json:"extra,omitempty"`
}
