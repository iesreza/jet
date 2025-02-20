// Copyright 2016 José Santos <henrique_1609@me.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jettest

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/iesreza/jet/v7"
)

func RunWithSet(t *testing.T, set *jet.Set, variables jet.VarMap, context interface{}, testName, testExpected string) {
	tt, err := set.GetTemplate(testName)
	if err != nil {
		t.Errorf("Error parsing templates for test %s: %v", testName, err)
		return
	}
	RunWithTemplate(t, tt, variables, context, testExpected)
}

func RunWithTemplate(t *testing.T, tt *jet.Template, variables jet.VarMap, context interface{}, testExpected string) {
	if testing.RunTests(func(pat, str string) (bool, error) {
		return true, nil
	}, []testing.InternalTest{
		{
			Name: fmt.Sprintf("\tJetTest(%s)", tt.Name),
			F: func(t *testing.T) {
				var buf bytes.Buffer
				err := tt.Execute(&buf, variables, context)
				if err != nil {
					t.Errorf("Eval error: %q executing %s", err.Error(), tt.Name)
					return
				}
				result := strings.Replace(buf.String(), "\r\n", "\n", -1)
				if result != testExpected {
					t.Errorf("Result error expected %q got %q on %s", testExpected, result, tt.Name)
				}
			},
		},
	}) == false {
		t.Fail()
	}
}
