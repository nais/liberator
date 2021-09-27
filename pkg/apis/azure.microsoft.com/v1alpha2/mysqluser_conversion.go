// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package azure_microsoft_com_v1alpha2

import "sigs.k8s.io/controller-runtime/pkg/conversion"

var _ conversion.Hub = &MySQLUser{}

func (*MySQLUser) Hub() {}
