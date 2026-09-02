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

package modelbuilder

import (
	"testing"

	"github.com/cloudwego/eino-ext/components/model/openai"

	"github.com/coze-dev/coze-studio/backend/api/model/admin/config"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
)

func TestApplyMoonshotKimiK26Params(t *testing.T) {
	conf := &openai.ChatModelConfig{
		MaxCompletionTokens: ptr.Of(4096),
		Temperature:         ptr.Of(float32(0.8)),
		TopP:                ptr.Of(float32(0.7)),
		FrequencyPenalty:    ptr.Of(float32(0.1)),
		PresencePenalty:     ptr.Of(float32(0.1)),
	}

	applyMoonshotKimiK26Params(conf, &LLMParams{
		MaxTokens:      4096,
		EnableThinking: ptr.Of(false),
	})

	if conf.MaxCompletionTokens != nil {
		t.Fatal("Moonshot k2.6 must not send max_completion_tokens")
	}
	if conf.MaxTokens == nil || *conf.MaxTokens != 4096 {
		t.Fatalf("unexpected disabled-thinking max_tokens: %#v", conf.MaxTokens)
	}
	if conf.Temperature != nil || conf.TopP != nil || conf.FrequencyPenalty != nil || conf.PresencePenalty != nil {
		t.Fatal("Moonshot k2.6 must not send restricted sampling parameters")
	}
	thinking, ok := conf.ExtraFields["thinking"].(map[string]string)
	if !ok || thinking["type"] != "disabled" {
		t.Fatalf("unexpected thinking config: %#v", conf.ExtraFields)
	}
}

func TestApplyMoonshotKimiK26ParamsUsesSafeThinkingMinimum(t *testing.T) {
	conf := &openai.ChatModelConfig{}

	applyMoonshotKimiK26Params(conf, &LLMParams{MaxTokens: 1})

	if conf.MaxTokens == nil || *conf.MaxTokens != 16000 {
		t.Fatalf("thinking calls need at least 16000 max_tokens, got %#v", conf.MaxTokens)
	}
	thinking := conf.ExtraFields["thinking"].(map[string]string)
	if thinking["type"] != "enabled" {
		t.Fatalf("unexpected thinking config: %#v", thinking)
	}
}

func TestApplyMoonshotKimiK26ParamsDefaultsTo32768AndEnabledThinking(t *testing.T) {
	conf := &openai.ChatModelConfig{}

	applyMoonshotKimiK26Params(conf, nil)

	if conf.MaxTokens == nil || *conf.MaxTokens != moonshotKimiK26MaxTokens {
		t.Fatalf("default max_tokens = %#v, want %d", conf.MaxTokens, moonshotKimiK26MaxTokens)
	}
	if len(conf.ExtraFields) != 1 {
		t.Fatalf("unexpected extra fields: %#v", conf.ExtraFields)
	}
	thinking, ok := conf.ExtraFields["thinking"].(map[string]string)
	if !ok || len(thinking) != 1 || thinking["type"] != "enabled" {
		t.Fatalf("unexpected thinking config: %#v", conf.ExtraFields)
	}
}

func TestNonMoonshotKimiK26ParamsRemainGeneric(t *testing.T) {
	base := &config.BaseConnectionInfo{
		BaseURL: moonshotBaseURL,
		Model:   "another-model",
	}
	if isMoonshotKimiK26(base) {
		t.Fatal("unexpected K2.6 match for a different model")
	}

	builder := &openaiModelBuilder{}
	conf := builder.getDefaultConfig()
	builder.applyParamsToOpenaiConfig(conf, &LLMParams{
		MaxTokens:        512,
		Temperature:      ptr.Of(float32(0.4)),
		TopP:             ptr.Of(float32(0.6)),
		FrequencyPenalty: 0.2,
		PresencePenalty:  0.1,
	})

	if conf.MaxCompletionTokens == nil || *conf.MaxCompletionTokens != 512 {
		t.Fatalf("generic max_completion_tokens = %#v, want 512", conf.MaxCompletionTokens)
	}
	if conf.MaxTokens != nil || conf.Temperature == nil || conf.TopP == nil ||
		conf.FrequencyPenalty == nil || conf.PresencePenalty == nil {
		t.Fatalf("generic OpenAI parameters were unexpectedly cleared: %#v", conf)
	}
	if conf.ExtraFields != nil {
		t.Fatalf("generic OpenAI model unexpectedly has Moonshot fields: %#v", conf.ExtraFields)
	}
}

func TestIsMoonshotKimiK26(t *testing.T) {
	if !isMoonshotKimiK26(&config.BaseConnectionInfo{
		BaseURL: moonshotBaseURL + "/",
		Model:   moonshotKimiK26ModelID,
	}) {
		t.Fatal("expected Moonshot k2.6 configuration to match")
	}
	if isMoonshotKimiK26(&config.BaseConnectionInfo{
		BaseURL: moonshotBaseURL,
		Model:   "gpt-4o",
	}) {
		t.Fatal("non-Kimi OpenAI configuration must not match")
	}
}
