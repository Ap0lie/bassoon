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

package intentdetector

import (
	"testing"

	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
)

func TestModelBuilderParamsDisablesThinking(t *testing.T) {
	params := (&Config{LLMParams: &vo.LLMParams{}}).modelBuilderParams()

	if params.EnableThinking == nil || *params.EnableThinking {
		t.Fatalf("expected intent detection to disable thinking, got %#v", params.EnableThinking)
	}
}
