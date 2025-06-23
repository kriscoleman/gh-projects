// Copyright 2025 Kris Coleman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package github

import "time"

type Project struct {
	ID     string
	Title  string
	Number int
	Owner  struct {
		Login string
	}
}

type Iteration struct {
	ID        string
	Title     string
	StartDate time.Time
	Duration  int
	Field     struct {
		ID   string
		Name string
	}
}

type Issue struct {
	ID         string
	Number     int
	Title      string
	State      string
	Repository struct {
		Name  string
		Owner struct {
			Login string
		}
	}
	ProjectItems struct {
		Nodes []struct {
			ID         string
			FieldValues struct {
				Nodes []FieldValue
			}
		}
	}
}

type FieldValue struct {
	TypeName string `json:"__typename"`
	Field    struct {
		ID   string
		Name string
	}
	Title    string
	Duration int
	ID       string
}