// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package replication_group

import (
	svcapitypes "github.com/aws-controllers-k8s/elasticache-controller/apis/v1alpha1"
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// CustomUpdateConditions sets conditions (terminal) on supplied replication group
// it examines supplied resource to determine conditions.
// It returns true if conditions are updated
func (rm *resourceManager) CustomUpdateConditions(
	ko *svcapitypes.ReplicationGroup,
	r *resource,
	err error,
) bool {
	rgStatus := r.ko.Status.Status
	if rgStatus == nil || *rgStatus != "create-failed" {
		return false
	}
	// Terminal condition
	var terminalCondition *ackv1alpha1.Condition = nil
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	} else {
		for _, condition := range ko.Status.Conditions {
			if condition.Type == ackv1alpha1.ConditionTypeTerminal {
				terminalCondition = condition
				break
			}
		}
		if terminalCondition != nil && terminalCondition.Status == corev1.ConditionTrue {
			// some other exception already put the resource in terminal condition
			return false
		}
	}
	if terminalCondition == nil {
		terminalCondition = &ackv1alpha1.Condition{
			Type: ackv1alpha1.ConditionTypeTerminal,
		}
		ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
	}
	terminalCondition.Status = corev1.ConditionTrue
	errorMessage := "Replication group status: create-failed"
	terminalCondition.Message = &errorMessage
	return true
}
