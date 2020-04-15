// Copyright 2020 Coinbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package asserter

import (
	"context"
	"errors"

	"github.com/coinbase/rosetta-sdk-go/types"
)

var (
	// ErrAsserterNotInitialized is returned when some call in the asserter
	// package requires the asserter to be initialized first.
	ErrAsserterNotInitialized = errors.New("asserter not initialized")
)

// Asserter contains all logic to perform static
// validation on Rosetta Server responses.
type Asserter struct {
	operationTypes     []string
	operationStatusMap map[string]bool
	errorTypeMap       map[int32]*types.Error
	genesisIndex       int64
}

// NewWithResponses constructs a new Asserter
// from a NetworkStatusResponse and
// NetworkOptionsResponse.
func NewWithResponses(
	ctx context.Context,
	networkStatus *types.NetworkStatusResponse,
	networkOptions *types.NetworkOptionsResponse,
) (*Asserter, error) {
	if err := NetworkStatusResponse(networkStatus); err != nil {
		return nil, err
	}

	if err := NetworkOptionsResponse(networkOptions); err != nil {
		return nil, err
	}

	return NewWithOptions(
		ctx,
		networkStatus.GenesisBlockIdentifier,
		networkOptions.Allow.OperationTypes,
		networkOptions.Allow.OperationStatuses,
		networkOptions.Allow.Errors,
	), nil
}

// NewWithOptions constructs a new Asserter using the provided
// arguments instead of using a NetworkStatusResponse and a
// NetworkOptionsResponse.
func NewWithOptions(
	ctx context.Context,
	genesisBlockIdentifier *types.BlockIdentifier,
	operationTypes []string,
	operationStatuses []*types.OperationStatus,
	errors []*types.Error,
) *Asserter {
	asserter := &Asserter{
		operationTypes: operationTypes,
		genesisIndex:   genesisBlockIdentifier.Index,
	}

	asserter.operationStatusMap = map[string]bool{}
	for _, status := range operationStatuses {
		asserter.operationStatusMap[status.Status] = status.Successful
	}

	asserter.errorTypeMap = map[int32]*types.Error{}
	for _, err := range errors {
		asserter.errorTypeMap[err.Code] = err
	}

	return asserter
}