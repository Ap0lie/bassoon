/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package llm

import (
	"encoding/json"
	"testing"

	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
)

func TestRequiresThinkingDisabledForTools(t *testing.T) {
	tests := []struct {
		name    string
		fcParam string
		want    bool
	}{
		{name: "no function call config", fcParam: `{}`, want: false},
		{name: "workflow tool", fcParam: `{"workflowFCParam":{"workflowList":[{"workflow_id":"1"}]}}`, want: true},
		{name: "plugin tool", fcParam: `{"pluginFCParam":{"pluginList":[{"plugin_id":"1","api_id":"2"}]}}`, want: true},
		{name: "knowledge tool", fcParam: `{"knowledgeFCParam":{"knowledgeList":[{"id":"3"}]}}`, want: true},
	}

	if (&Config{}).requiresThinkingDisabledForTools() {
		t.Fatal("nil function-call configuration must not disable thinking")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fcParam vo.FCParam
			if err := json.Unmarshal([]byte(tt.fcParam), &fcParam); err != nil {
				t.Fatalf("unmarshal function-call config: %v", err)
			}

			if got := (&Config{FCParam: &fcParam}).requiresThinkingDisabledForTools(); got != tt.want {
				t.Fatalf("requiresThinkingDisabledForTools() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModelBuilderParamsDisablesThinkingForToolCalls(t *testing.T) {
	var fcParam vo.FCParam
	if err := json.Unmarshal([]byte(`{"pluginFCParam":{"pluginList":[{"plugin_id":"1","api_id":"2"}]}}`), &fcParam); err != nil {
		t.Fatalf("unmarshal function-call config: %v", err)
	}

	c := &Config{
		LLMParams: &vo.LLMParams{},
		FCParam:   &fcParam,
	}
	params := c.modelBuilderParams(c.LLMParams)
	if params.EnableThinking == nil || *params.EnableThinking {
		t.Fatalf("expected tool call to explicitly disable thinking, got %#v", params.EnableThinking)
	}
}
