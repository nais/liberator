// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package azure_microsoft_com_api

// +kubebuilder:validation:Enum={"CreateOrUpdate","Delete"}
type PollingURLKind string

const (
	PollingURLKindCreateOrUpdate = PollingURLKind("CreateOrUpdate")
	PollingURLKindDelete         = PollingURLKind("Delete")
)
