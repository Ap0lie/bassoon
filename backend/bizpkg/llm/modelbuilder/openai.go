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
	"context"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/openai"

	"github.com/coze-dev/coze-studio/backend/api/model/admin/config"
	"github.com/coze-dev/coze-studio/backend/api/model/app/bot_common"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
)

type openaiModelBuilder struct {
	cfg *config.Model
}

const (
	moonshotBaseURL          = "https://api.moonshot.cn/v1"
	moonshotKimiK26ModelID   = "kimi-k2.6"
	moonshotKimiK26MaxTokens = 32768
)

func newOpenaiModelBuilder(cfg *config.Model) Service {
	return &openaiModelBuilder{
		cfg: cfg,
	}
}

func (o *openaiModelBuilder) getDefaultConfig() *openai.ChatModelConfig {
	return &openai.ChatModelConfig{
		MaxCompletionTokens: ptr.Of(4096),
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type:       "text",
			JSONSchema: nil,
		},
	}
}

func (o *openaiModelBuilder) applyParamsToOpenaiConfig(conf *openai.ChatModelConfig, params *LLMParams) {
	if params == nil {
		return
	}

	if params.Temperature != nil {
		conf.Temperature = ptr.Of(*params.Temperature)
	}

	if params.MaxTokens != 0 {
		conf.MaxCompletionTokens = ptr.Of(params.MaxTokens)
	}

	if params.FrequencyPenalty != 0 {
		conf.FrequencyPenalty = ptr.Of(params.FrequencyPenalty)
	}

	if params.PresencePenalty != 0 {
		conf.PresencePenalty = ptr.Of(params.PresencePenalty)
	}

	conf.TopP = params.TopP

	if params.ResponseFormat == bot_common.ModelResponseFormat_JSON {
		conf.ResponseFormat = &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		}
	} else {
		conf.ResponseFormat = &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeText,
		}
	}
}

func isMoonshotKimiK26(base *config.BaseConnectionInfo) bool {
	return base != nil &&
		strings.TrimRight(base.BaseURL, "/") == moonshotBaseURL &&
		base.Model == moonshotKimiK26ModelID
}

func applyMoonshotKimiK26Params(conf *openai.ChatModelConfig, params *LLMParams) {
	maxTokens := moonshotKimiK26MaxTokens
	thinkingType := "enabled"
	if params != nil {
		if params.MaxTokens > 0 {
			maxTokens = params.MaxTokens
		}
		if params.EnableThinking != nil && !*params.EnableThinking {
			thinkingType = "disabled"
		}
	}
	if thinkingType == "enabled" && maxTokens < 16000 {
		maxTokens = 16000
	}

	// Moonshot k2.6 rejects the sampling parameters exposed by generic OpenAI
	// models. It also expects max_tokens instead of max_completion_tokens.
	conf.MaxCompletionTokens = nil
	conf.MaxTokens = ptr.Of(maxTokens)
	conf.Temperature = nil
	conf.TopP = nil
	conf.FrequencyPenalty = nil
	conf.PresencePenalty = nil
	conf.ExtraFields = map[string]any{
		"thinking": map[string]string{"type": thinkingType},
	}
}

func (o *openaiModelBuilder) Build(ctx context.Context, params *LLMParams) (ToolCallingChatModel, error) {
	base := o.cfg.Connection.BaseConnInfo

	conf := o.getDefaultConfig()
	conf.APIKey = base.APIKey
	conf.Model = base.Model

	if base.BaseURL != "" {
		conf.BaseURL = base.BaseURL
	}

	if o.cfg.Connection.Openai != nil {
		conf.APIVersion = o.cfg.Connection.Openai.APIVersion
		conf.ByAzure = o.cfg.Connection.Openai.ByAzure
	}

	o.applyParamsToOpenaiConfig(conf, params)
	if isMoonshotKimiK26(base) {
		applyMoonshotKimiK26Params(conf, params)
	}

	return openai.NewChatModel(ctx, conf)
}
