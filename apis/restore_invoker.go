/*
Copyright The Stash Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package apis

import (
	"context"
	"fmt"
	"time"

	"stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	cs "stash.appscode.dev/apimachinery/client/clientset/versioned"
	stash_scheme "stash.appscode.dev/apimachinery/client/clientset/versioned/scheme"
	v1beta1_util "stash.appscode.dev/apimachinery/client/clientset/versioned/typed/stash/v1beta1/util"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/reference"
	kmapi "kmodules.xyz/client-go/api/v1"
	core_util "kmodules.xyz/client-go/core/v1"
	"kmodules.xyz/client-go/meta"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	EventSourceRestoreBatchController   = "RestoreBatch Controller"
	EventSourceRestoreSessionController = "RestoreSession Controller"
)

type RestoreTargetInfo struct {
	Task                  v1beta1.TaskRef
	Target                *v1beta1.RestoreTarget
	RuntimeSettings       ofst.RuntimeSettings
	TempDir               v1beta1.EmptyDirSettings
	InterimVolumeTemplate *ofst.PersistentVolumeClaim
	Hooks                 *v1beta1.RestoreHooks
}

type RestoreInvokerStatus struct {
	Phase           v1beta1.RestorePhase
	SessionDuration string
	Conditions      []kmapi.Condition
	TargetStatus    []v1beta1.RestoreMemberStatus
}

type RestoreInvoker struct {
	TypeMeta        metav1.TypeMeta
	ObjectMeta      metav1.ObjectMeta
	Labels          map[string]string
	Hash            string
	Driver          v1beta1.Snapshotter
	Repository      string
	TargetsInfo     []RestoreTargetInfo
	ExecutionOrder  v1beta1.ExecutionOrder
	Hooks           *v1beta1.RestoreHooks
	ObjectRef       *core.ObjectReference
	OwnerRef        *metav1.OwnerReference
	Status          RestoreInvokerStatus
	ObjectJson      []byte
	AddFinalizer    func() error
	RemoveFinalizer func() error
	HasCondition    func(*v1beta1.TargetRef, string) (bool, error)
	GetCondition    func(*v1beta1.TargetRef, string) (int, *kmapi.Condition, error)
	SetCondition    func(*v1beta1.TargetRef, kmapi.Condition) error
	IsConditionTrue func(*v1beta1.TargetRef, string) (bool, error)
	NextInOrder     func(v1beta1.TargetRef, []v1beta1.RestoreMemberStatus) bool

	UpdateTargetStatus func(v1beta1.RestoreMemberStatus) error
	UpdateRestorePhase func(v1beta1.RestorePhase, string) error
	CreateEvent        func(eventType, source, reason, message string) error
}

func ExtractRestoreInvokerInfo(kubeClient kubernetes.Interface, stashClient cs.Interface, invokerType, invokerName, namespace string) (RestoreInvoker, error) {
	var invoker RestoreInvoker
	switch invokerType {
	case v1beta1.ResourceKindRestoreBatch:
		// get RestoreBatch
		restoreBatch, err := stashClient.StashV1beta1().RestoreBatches(namespace).Get(context.TODO(), invokerName, metav1.GetOptions{})
		if err != nil {
			return invoker, err
		}
		invoker.TypeMeta = metav1.TypeMeta{
			Kind:       v1beta1.ResourceKindRestoreBatch,
			APIVersion: v1beta1.SchemeGroupVersion.String(),
		}
		invoker.ObjectMeta = restoreBatch.ObjectMeta
		invoker.Labels = restoreBatch.OffshootLabels()
		invoker.Hash = restoreBatch.GetSpecHash()
		invoker.Driver = restoreBatch.Spec.Driver
		invoker.Repository = restoreBatch.Spec.Repository.Name
		invoker.Hooks = restoreBatch.Spec.Hooks
		invoker.ExecutionOrder = restoreBatch.Spec.ExecutionOrder
		invoker.OwnerRef = metav1.NewControllerRef(restoreBatch, v1beta1.SchemeGroupVersion.WithKind(v1beta1.ResourceKindRestoreBatch))
		invoker.ObjectRef, err = reference.GetReference(stash_scheme.Scheme, restoreBatch)
		if err != nil {
			return invoker, err
		}

		invoker.ObjectJson, err = meta.MarshalToJson(restoreBatch, v1beta1.SchemeGroupVersion)
		if err != nil {
			return invoker, err
		}

		for _, member := range restoreBatch.Spec.Members {
			invoker.TargetsInfo = append(invoker.TargetsInfo, RestoreTargetInfo{
				Task:                  member.Task,
				Target:                member.Target,
				RuntimeSettings:       member.RuntimeSettings,
				TempDir:               member.TempDir,
				InterimVolumeTemplate: member.InterimVolumeTemplate,
				Hooks:                 member.Hooks,
			})
		}

		invoker.Status.Phase = restoreBatch.Status.Phase
		invoker.Status.Conditions = restoreBatch.Status.Conditions
		invoker.Status.SessionDuration = restoreBatch.Status.SessionDuration
		invoker.Status.TargetStatus = restoreBatch.Status.Members

		invoker.AddFinalizer = func() error {
			_, _, err := v1beta1_util.PatchRestoreBatch(context.TODO(), stashClient.StashV1beta1(), restoreBatch, func(in *v1beta1.RestoreBatch) *v1beta1.RestoreBatch {
				in.ObjectMeta = core_util.AddFinalizer(in.ObjectMeta, v1beta1.StashKey)
				return in
			}, metav1.PatchOptions{})
			return err
		}
		invoker.RemoveFinalizer = func() error {
			_, _, err := v1beta1_util.PatchRestoreBatch(context.TODO(), stashClient.StashV1beta1(), restoreBatch, func(in *v1beta1.RestoreBatch) *v1beta1.RestoreBatch {
				in.ObjectMeta = core_util.RemoveFinalizer(in.ObjectMeta, v1beta1.StashKey)
				return in
			}, metav1.PatchOptions{})
			return err
		}
		invoker.HasCondition = func(target *v1beta1.TargetRef, condType string) (bool, error) {
			restoreBatch, err := stashClient.StashV1beta1().RestoreBatches(namespace).Get(context.TODO(), invokerName, metav1.GetOptions{})
			if err != nil {
				return false, err
			}
			if target != nil {
				return hasRestoreMemberCondition(restoreBatch.Status.Members, *target, condType), nil
			}
			return kmapi.HasCondition(restoreBatch.Status.Conditions, condType), nil
		}
		invoker.GetCondition = func(target *v1beta1.TargetRef, condType string) (int, *kmapi.Condition, error) {
			restoreBatch, err := stashClient.StashV1beta1().RestoreBatches(namespace).Get(context.TODO(), invokerName, metav1.GetOptions{})
			if err != nil {
				return -1, nil, err
			}
			if target != nil {
				idx, cond := getRestoreMemberCondition(restoreBatch.Status.Members, *target, condType)
				return idx, cond, nil
			}
			idx, cond := kmapi.GetCondition(restoreBatch.Status.Conditions, condType)
			return idx, cond, nil

		}
		invoker.SetCondition = func(target *v1beta1.TargetRef, condition kmapi.Condition) error {
			_, err = v1beta1_util.UpdateRestoreBatchStatus(context.TODO(), stashClient.StashV1beta1(), restoreBatch.ObjectMeta, func(in *v1beta1.RestoreBatchStatus) *v1beta1.RestoreBatchStatus {
				if target != nil {
					in.Members = setRestoreMemberCondition(in.Members, *target, condition)
				} else {
					in.Conditions = kmapi.SetCondition(in.Conditions, condition)
				}
				return in
			}, metav1.UpdateOptions{})
			return err
		}
		invoker.IsConditionTrue = func(target *v1beta1.TargetRef, condType string) (bool, error) {
			restoreBatch, err := stashClient.StashV1beta1().RestoreBatches(namespace).Get(context.TODO(), invokerName, metav1.GetOptions{})
			if err != nil {
				return false, err
			}
			if target != nil {
				return isRestoreMemberConditionTrue(restoreBatch.Status.Members, *target, condType), nil
			}
			return kmapi.IsConditionTrue(restoreBatch.Status.Conditions, condType), nil
		}

		invoker.NextInOrder = func(ref v1beta1.TargetRef, targets []v1beta1.RestoreMemberStatus) bool {
			for i := range targets {
				if TargetMatched(ref, targets[i].Ref) && targets[i].Phase == "" {
					break
				}
				if targets[i].Phase != v1beta1.TargetRestoreSucceeded {
					return false
				}
			}
			return true
		}
		invoker.UpdateTargetStatus = func(memberStatus v1beta1.RestoreMemberStatus) error {
			_, err = v1beta1_util.UpdateRestoreBatchStatus(
				context.TODO(),
				stashClient.StashV1beta1(),
				invoker.ObjectMeta,
				func(in *v1beta1.RestoreBatchStatus) *v1beta1.RestoreBatchStatus {
					in.Members = upsertRestoreMemberStatus(in.Members, memberStatus)
					return in
				},
				metav1.UpdateOptions{},
			)
			return err
		}
		invoker.UpdateRestorePhase = func(phase v1beta1.RestorePhase, sessionDuration string) error {
			_, err = v1beta1_util.UpdateRestoreBatchStatus(
				context.TODO(),
				stashClient.StashV1beta1(),
				invoker.ObjectMeta,
				func(in *v1beta1.RestoreBatchStatus) *v1beta1.RestoreBatchStatus {
					if phase != "" {
						in.Phase = phase
					}
					if sessionDuration != "" {
						in.SessionDuration = sessionDuration
					}
					return in
				},
				metav1.UpdateOptions{},
			)
			return err
		}
		invoker.CreateEvent = func(eventType, source, reason, message string) error {
			t := metav1.Time{Time: time.Now()}

			if source == "" {
				source = EventSourceRestoreBatchController
			}
			_, err := kubeClient.CoreV1().Events(invoker.ObjectMeta.Namespace).Create(context.TODO(), &core.Event{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%v.%x", invoker.ObjectRef.Name, t.UnixNano()),
					Namespace: invoker.ObjectRef.Namespace,
				},
				InvolvedObject: *invoker.ObjectRef,
				Reason:         reason,
				Message:        message,
				FirstTimestamp: t,
				LastTimestamp:  t,
				Count:          1,
				Type:           eventType,
				Source:         core.EventSource{Component: source},
			}, metav1.CreateOptions{})
			return err
		}
	case v1beta1.ResourceKindRestoreSession:
		// get RestoreSession
		restoreSession, err := stashClient.StashV1beta1().RestoreSessions(namespace).Get(context.TODO(), invokerName, metav1.GetOptions{})
		if err != nil {
			return invoker, err
		}
		invoker.TypeMeta = metav1.TypeMeta{
			Kind:       v1beta1.ResourceKindRestoreSession,
			APIVersion: v1beta1.SchemeGroupVersion.String(),
		}
		invoker.ObjectMeta = restoreSession.ObjectMeta
		invoker.Labels = restoreSession.OffshootLabels()
		invoker.Hash = restoreSession.GetSpecHash()
		invoker.Driver = restoreSession.Spec.Driver
		invoker.Repository = restoreSession.Spec.Repository.Name
		invoker.OwnerRef = metav1.NewControllerRef(restoreSession, v1beta1.SchemeGroupVersion.WithKind(v1beta1.ResourceKindRestoreSession))
		invoker.ObjectRef, err = reference.GetReference(stash_scheme.Scheme, restoreSession)
		if err != nil {
			return invoker, err
		}

		invoker.ObjectJson, err = meta.MarshalToJson(restoreSession, v1beta1.SchemeGroupVersion)
		if err != nil {
			return invoker, err
		}

		invoker.TargetsInfo = append(invoker.TargetsInfo, RestoreTargetInfo{
			Task:                  restoreSession.Spec.Task,
			Target:                restoreSession.Spec.Target,
			RuntimeSettings:       restoreSession.Spec.RuntimeSettings,
			TempDir:               restoreSession.Spec.TempDir,
			InterimVolumeTemplate: restoreSession.Spec.InterimVolumeTemplate,
			Hooks:                 restoreSession.Spec.Hooks,
		})

		invoker.Status.Phase = restoreSession.Status.Phase
		invoker.Status.Conditions = restoreSession.Status.Conditions
		invoker.Status.SessionDuration = restoreSession.Status.SessionDuration
		if restoreSession.Spec.Target != nil {
			invoker.Status.TargetStatus = append(invoker.Status.TargetStatus, v1beta1.RestoreMemberStatus{
				Ref:        restoreSession.Spec.Target.Ref,
				Conditions: restoreSession.Status.Conditions,
				TotalHosts: restoreSession.Status.TotalHosts,
				Phase:      v1beta1.RestoreTargetPhase(restoreSession.Status.Phase),
				Stats:      restoreSession.Status.Stats,
			})
		}

		invoker.AddFinalizer = func() error {
			_, _, err := v1beta1_util.PatchRestoreSession(context.TODO(), stashClient.StashV1beta1(), restoreSession, func(in *v1beta1.RestoreSession) *v1beta1.RestoreSession {
				in.ObjectMeta = core_util.AddFinalizer(in.ObjectMeta, v1beta1.StashKey)
				return in
			}, metav1.PatchOptions{})
			return err
		}
		invoker.RemoveFinalizer = func() error {
			_, _, err := v1beta1_util.PatchRestoreSession(context.TODO(), stashClient.StashV1beta1(), restoreSession, func(in *v1beta1.RestoreSession) *v1beta1.RestoreSession {
				in.ObjectMeta = core_util.RemoveFinalizer(in.ObjectMeta, v1beta1.StashKey)
				return in
			}, metav1.PatchOptions{})
			return err
		}
		invoker.HasCondition = func(target *v1beta1.TargetRef, condType string) (bool, error) {
			restoreSession, err := stashClient.StashV1beta1().RestoreSessions(namespace).Get(context.TODO(), invokerName, metav1.GetOptions{})
			if err != nil {
				return false, err
			}
			return kmapi.HasCondition(restoreSession.Status.Conditions, condType), nil
		}
		invoker.GetCondition = func(target *v1beta1.TargetRef, condType string) (int, *kmapi.Condition, error) {
			restoreSession, err := stashClient.StashV1beta1().RestoreSessions(namespace).Get(context.TODO(), invokerName, metav1.GetOptions{})
			if err != nil {
				return -1, nil, err
			}
			idx, cond := kmapi.GetCondition(restoreSession.Status.Conditions, condType)
			return idx, cond, nil
		}
		invoker.SetCondition = func(target *v1beta1.TargetRef, condition kmapi.Condition) error {
			_, err = v1beta1_util.UpdateRestoreSessionStatus(context.TODO(), stashClient.StashV1beta1(), restoreSession.ObjectMeta, func(in *v1beta1.RestoreSessionStatus) *v1beta1.RestoreSessionStatus {
				in.Conditions = kmapi.SetCondition(in.Conditions, condition)
				return in
			}, metav1.UpdateOptions{})
			return err
		}
		invoker.IsConditionTrue = func(target *v1beta1.TargetRef, condType string) (bool, error) {
			restoreSession, err := stashClient.StashV1beta1().RestoreSessions(namespace).Get(context.TODO(), invokerName, metav1.GetOptions{})
			if err != nil {
				return false, err
			}
			return kmapi.IsConditionTrue(restoreSession.Status.Conditions, condType), nil
		}

		invoker.NextInOrder = func(ref v1beta1.TargetRef, targets []v1beta1.RestoreMemberStatus) bool {
			for i := range targets {
				if TargetMatched(ref, targets[i].Ref) && targets[i].Phase == "" {
					break
				}
				if targets[i].Phase != v1beta1.TargetRestoreSucceeded {
					return false
				}
			}
			return true
		}

		invoker.UpdateTargetStatus = func(memberStatus v1beta1.RestoreMemberStatus) error {
			_, err = v1beta1_util.UpdateRestoreSessionStatus(
				context.TODO(),
				stashClient.StashV1beta1(),
				invoker.ObjectMeta,
				func(in *v1beta1.RestoreSessionStatus) *v1beta1.RestoreSessionStatus {
					in.TotalHosts = memberStatus.TotalHosts
					in.Stats = upsertRestoreHostStatus(in.Stats, memberStatus.Stats)
					in.Conditions = upsertConditions(in.Conditions, memberStatus.Conditions)
					return in
				},
				metav1.UpdateOptions{},
			)
			return err
		}
		invoker.UpdateRestorePhase = func(phase v1beta1.RestorePhase, sessionDuration string) error {
			_, err = v1beta1_util.UpdateRestoreSessionStatus(
				context.TODO(),
				stashClient.StashV1beta1(),
				invoker.ObjectMeta,
				func(in *v1beta1.RestoreSessionStatus) *v1beta1.RestoreSessionStatus {
					if phase != "" {
						in.Phase = phase
					}
					if sessionDuration != "" {
						in.SessionDuration = sessionDuration
					}
					return in
				},
				metav1.UpdateOptions{},
			)
			return err
		}
		invoker.CreateEvent = func(eventType, source, reason, message string) error {
			t := metav1.Time{Time: time.Now()}
			if source == "" {
				source = EventSourceRestoreSessionController
			}
			_, err := kubeClient.CoreV1().Events(invoker.ObjectMeta.Namespace).Create(context.TODO(), &core.Event{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%v.%x", invoker.ObjectRef.Name, t.UnixNano()),
					Namespace: invoker.ObjectRef.Namespace,
				},
				InvolvedObject: *invoker.ObjectRef,
				Reason:         reason,
				Message:        message,
				FirstTimestamp: t,
				LastTimestamp:  t,
				Count:          1,
				Type:           eventType,
				Source:         core.EventSource{Component: source},
			}, metav1.CreateOptions{})
			return err
		}
	default:
		return invoker, fmt.Errorf("failed to extract invoker info. Reason: unknown invoker")
	}
	return invoker, nil
}

func hasRestoreMemberCondition(status []v1beta1.RestoreMemberStatus, target v1beta1.TargetRef, condType string) bool {
	// If the target is present in the list, then return the respective value
	for i := range status {
		if TargetMatched(status[i].Ref, target) {
			return kmapi.HasCondition(status[i].Conditions, condType)
		}
	}
	// Member is not present in the list, so the condition is not there too
	return false
}

func getRestoreMemberCondition(status []v1beta1.RestoreMemberStatus, target v1beta1.TargetRef, condType string) (int, *kmapi.Condition) {
	// If the target is present in the list, then return the respective condition
	for i := range status {
		if TargetMatched(status[i].Ref, target) {
			return kmapi.GetCondition(status[i].Conditions, condType)
		}
	}
	// Member is not present in the list
	return -1, nil
}

func setRestoreMemberCondition(status []v1beta1.RestoreMemberStatus, target v1beta1.TargetRef, newCondition kmapi.Condition) []v1beta1.RestoreMemberStatus {
	// If the target is already exist in the list, update its condition
	for i := range status {
		if TargetMatched(status[i].Ref, target) {
			status[i].Conditions = kmapi.SetCondition(status[i].Conditions, newCondition)
			return status
		}
	}
	// The target does not exist in the list. So, add a new entry.
	memberStatus := v1beta1.RestoreMemberStatus{
		Ref:        target,
		Conditions: kmapi.SetCondition(nil, newCondition),
	}
	return append(status, memberStatus)
}

func isRestoreMemberConditionTrue(status []v1beta1.RestoreMemberStatus, target v1beta1.TargetRef, condType string) bool {
	// If the target is present in the list, then return the respective value
	for i := range status {
		if TargetMatched(status[i].Ref, target) {
			return kmapi.IsConditionTrue(status[i].Conditions, condType)
		}
	}
	// Member is not present in the list, so the condition is false
	return false
}

func upsertRestoreMemberStatus(cur []v1beta1.RestoreMemberStatus, new v1beta1.RestoreMemberStatus) []v1beta1.RestoreMemberStatus {
	// if the member status already exist, then update it
	for i := range cur {
		if TargetMatched(cur[i].Ref, new.Ref) {
			cur[i].Ref = new.Ref
			cur[i].Conditions = upsertConditions(cur[i].Conditions, new.Conditions)
			if new.TotalHosts != nil {
				cur[i].TotalHosts = new.TotalHosts
			}
			cur[i].Stats = upsertRestoreHostStatus(cur[i].Stats, new.Stats)
		}
	}
	// the member status does not exist. so, add new entry.
	cur = append(cur, new)
	return cur
}

func upsertConditions(cur []kmapi.Condition, new []kmapi.Condition) []kmapi.Condition {
	for i := range new {
		cur = kmapi.SetCondition(cur, new[i])
	}
	return cur
}

func upsertRestoreHostStatus(cur []v1beta1.HostRestoreStats, new []v1beta1.HostRestoreStats) []v1beta1.HostRestoreStats {
	for i := range new {
		index, hostEntryExist := hostEntryIndex(cur, new[i])
		if hostEntryExist {
			cur[index] = new[i]
		} else {
			cur = append(cur, new[i])
		}
	}
	return cur
}

func hostEntryIndex(entries []v1beta1.HostRestoreStats, target v1beta1.HostRestoreStats) (int, bool) {
	for i := range entries {
		if entries[i].Hostname == target.Hostname {
			return i, true
		}
	}
	return -1, false
}