package mote

type Manifest struct {
	Name                  string          `json:"name"`
	Workspace             string          `json:"workspace,omitempty"`
	Capabilities          []string        `json:"capabilities,omitempty"`
	Models                []ModelConfig   `json:"models,omitempty"`
	MaxIterations         *int            `json:"max_iterations,omitempty"`
	MaxPerTool            *int            `json:"max_per_tool,omitempty"`
	MaxRuntimeSeconds     *int            `json:"max_runtime_seconds,omitempty"`
	CommandTimeoutSeconds *int            `json:"command_timeout_seconds,omitempty"`
	MaxOutputBytes        *int            `json:"max_output_bytes,omitempty"`
	OutputFile            *string         `json:"output_file,omitempty"`
	Shell                 *ShellConfig    `json:"shell,omitempty"`
	Decision              *DecisionConfig `json:"decision,omitempty"`
}

type ModelConfig struct {
	Provider       string `json:"provider"`
	Model          string `json:"model"`
	Endpoint       string `json:"endpoint"`
	AuthEnv        string `json:"auth_env,omitempty"`
	TimeoutSeconds *int   `json:"timeout_seconds,omitempty"`
}

type ShellConfig struct {
	AllowedPrograms []string `json:"allowed_programs,omitempty"`
	UnsafeShell     bool     `json:"unsafe_shell,omitempty"`
}

type DecisionConfig struct {
	Provider       string `json:"provider"`
	Endpoint       string `json:"endpoint"`
	Model          string `json:"model,omitempty"`
	AuthEnv        string `json:"auth_env,omitempty"`
	Fallback       string `json:"fallback,omitempty"`
	TimeoutSeconds *int   `json:"timeout_seconds,omitempty"`
}
