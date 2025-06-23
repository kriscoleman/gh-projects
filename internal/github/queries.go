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

const GetProjectQuery = `
query($owner: String!, $number: Int!) {
  user(login: $owner) {
    projectV2(number: $number) {
      id
      title
      number
    }
  }
  organization(login: $owner) {
    projectV2(number: $number) {
      id
      title
      number
    }
  }
}
`

const GetProjectFieldsQuery = `
query($projectId: ID!) {
  node(id: $projectId) {
    ... on ProjectV2 {
      fields(first: 100) {
        nodes {
          ... on ProjectV2Field {
            id
            name
            dataType
          }
          ... on ProjectV2IterationField {
            id
            name
            dataType
            configuration {
              iterations {
                id
                title
                startDate
                duration
              }
              completedIterations {
                id
                title
                startDate
                duration
              }
            }
          }
        }
      }
    }
  }
}
`

const GetIterationItemsQuery = `
query($projectId: ID!, $after: String) {
  node(id: $projectId) {
    ... on ProjectV2 {
      items(first: 100, after: $after) {
        pageInfo {
          hasNextPage
          endCursor
        }
        nodes {
          id
          content {
            ... on Issue {
              id
              number
              title
              state
              repository {
                name
                owner {
                  login
                }
              }
            }
          }
          fieldValues(first: 20) {
            nodes {
              ... on ProjectV2ItemFieldIterationValue {
                __typename
                field {
                  ... on ProjectV2IterationField {
                    id
                    name
                  }
                }
                iterationId
                title
              }
              ... on ProjectV2ItemFieldSingleSelectValue {
                __typename
                field {
                  ... on ProjectV2SingleSelectField {
                    id
                    name
                  }
                }
                name
              }
            }
          }
        }
      }
    }
  }
}
`

const UpdateItemIterationMutation = `
mutation($projectId: ID!, $itemId: ID!, $fieldId: ID!, $iterationId: String!) {
  updateProjectV2ItemFieldValue(input: {
    projectId: $projectId
    itemId: $itemId
    fieldId: $fieldId
    value: {
      iterationId: $iterationId
    }
  }) {
    projectV2Item {
      id
    }
  }
}
`