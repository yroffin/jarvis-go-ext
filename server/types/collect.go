/**
 * Copyright 2017 Yannick Roffin
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package types

import (
	"gopkg.in/mgo.v2/bson"
)

// PostCollectResource post query
type PostCollectResource struct {
	Find    bson.M   `json:"find,omitempty"`
	OrderBy []string `json:"orderby,omitempty"`
}

// GetCollectResource query
type GetCollectResource struct {
	Collections []GetCollectResourceDetail `json:"collections,omitempty"`
}

// GetCollectResourceDetail detail
type GetCollectResourceDetail struct {
	Name   string `json:"name,omitempty"`
	Entity bson.M `json:"entity,omitempty"`
}

// CollectResource : connector struct
type CollectResource struct {
	Collections []string `json:"collections,omitempty"`
	Data        []bson.M `json:"data,omitempty"`
}
