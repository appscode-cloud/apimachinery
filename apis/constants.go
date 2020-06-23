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

import "time"

const (
	StashDockerRegistry = "STASH_DOCKER_REGISTRY"
	StashDockerImage    = "STASH_DOCKER_IMAGE"
	StashImageTag       = "STASH_IMAGE_TAG"
	ImageStash          = "stash"

	Namespace      = "NAMESPACE"
	BackupSession  = "BACKUP_SESSION"
	RestoreSession = "RESTORE_SESSION"

	RepositoryName       = "REPOSITORY_NAME"
	RepositoryProvider   = "REPOSITORY_PROVIDER"
	RepositorySecretName = "REPOSITORY_SECRET_NAME"
	RepositoryBucket     = "REPOSITORY_BUCKET"
	RepositoryPrefix     = "REPOSITORY_PREFIX"
	RepositoryEndpoint   = "REPOSITORY_ENDPOINT"
	RepositoryRegion     = "REPOSITORY_REGION"

	Hostname       = "HOSTNAME"
	SourceHostname = "SOURCE_HOSTNAME"
	InterimDataDir = "INTERIM_DATA_DIR"

	TargetName       = "TARGET_NAME"
	TargetAPIVersion = "TARGET_API_VERSION"
	TargetKind       = "TARGET_KIND"
	TargetResource   = "TARGET_RESOURCE"
	TargetNamespace  = "TARGET_NAMESPACE"
	TargetMountPath  = "TARGET_MOUNT_PATH"
	TargetPaths      = "TARGET_PATHS"

	TargetAppVersion  = "TARGET_APP_VERSION"
	TargetAppType     = "TARGET_APP_TYPE"
	TargetAppGroup    = "TARGET_APP_GROUP"
	TargetAppResource = "TARGET_APP_RESOURCE"
	TargetAppReplicas = "TARGET_APP_REPLICAS"

	InvokerKind = "INVOKER_KIND"
	InvokerName = "INVOKER_NAME"

	ExcludePatterns = "EXCLUDE_PATTERNS"
	IncludePatterns = "INCLUDE_PATTERNS"

	RestorePaths     = "RESTORE_PATHS"
	RestoreSnapshots = "RESTORE_SNAPSHOTS"

	RetentionKeepLast    = "RETENTION_KEEP_LAST"
	RetentionKeepHourly  = "RETENTION_KEEP_HOURLY"
	RetentionKeepDaily   = "RETENTION_KEEP_DAILY"
	RetentionKeepWeekly  = "RETENTION_KEEP_WEEKLY"
	RetentionKeepMonthly = "RETENTION_KEEP_MONTHLY"
	RetentionKeepYearly  = "RETENTION_KEEP_YEARLY"
	RetentionKeepTags    = "RETENTION_KEEP_TAGS"
	RetentionPrune       = "RETENTION_PRUNE"
	RetentionDryRun      = "RETENTION_DRY_RUN"

	// default true
	// false when TmpDir.DisableCaching is true in backupConfig/restoreSession
	EnableCache    = "ENABLE_CACHE"
	MaxConnections = "MAX_CONNECTIONS"

	PushgatewayURL    = "PROMETHEUS_PUSHGATEWAY_URL"
	PrometheusJobName = "PROMETHEUS_JOB_NAME"

	StashDefaultVolume          = "stash-volume"
	StashDefaultMountPath       = "/stash-data"
	StashInterimVolume          = "stash-interim-volume"
	StashInterimVolumeMountPath = "/stash-interim-volume"
	StashInterimDataDir         = "/stash-interim-volume/data"

	// backup or restore hooks
	HookType        = "HOOK_TYPE"
	PreBackupHook   = "preBackup"
	PostBackupHook  = "postBackup"
	PreRestoreHook  = "preRestore"
	PostRestoreHook = "postRestore"
	PreTaskHook     = "pre-task-hook"
	PostTaskHook    = "post-task-hook"
)

const (
	PrefixStashBackup         = "stash-backup"
	PrefixStashRestore        = "stash-restore"
	PrefixStashVolumeSnapshot = "stash-vs"

	StashContainer       = "stash"
	StashInitContainer   = "stash-init"
	LocalVolumeName      = "stash-local"
	ScratchDirVolumeName = "stash-scratchdir"
	TmpDirVolumeName     = "tmp-dir"
	TmpDirMountPath      = "/tmp"
	PodinfoVolumeName    = "stash-podinfo"

	RecoveryJobPrefix   = "stash-recovery-"
	ScaledownCronPrefix = "stash-scaledown-cron-"
	CheckJobPrefix      = "stash-check-"

	AnnotationRestic     = "restic"
	AnnotationRecovery   = "recovery"
	AnnotationOperation  = "operation"
	AnnotationOldReplica = "old-replica"

	OperationRecovery = "recovery"
	OperationCheck    = "check"

	AppLabelStash        = "stash"
	AppLabelStashV1Beta1 = "stash-v1beta1"
	OperationScaleDown   = "scale-down"

	RepositoryFinalizer = "stash"
	SnapshotIDLength    = 8

	ModelSidecar        = "sidecar"
	ModelCronJob        = "cronjob"
	LabelApp            = "app"
	LabelInvokerType    = StashKey + "/invoker-type"
	LabelInvokerName    = StashKey + "/invoker-name"
	StashSecretVolume   = "stash-secret-volume"
	StashSecretMountDir = "/etc/stash/repository/secret"

	KeyPodName    = "POD_NAME"
	KeyNodeName   = "NODE_NAME"
	KeyPodOrdinal = "POD_ORDINAL"

	RetryInterval    = 50 * time.Millisecond
	ReadinessTimeout = 2 * time.Minute
)

const (
	CallerWebhook       = "webhook"
	CallerController    = "controller"
	PushgatewayLocalURL = "http://localhost:56789"
	DefaultHost         = "host-0"
)

// ==================== Prometheus metrics related constants ============
const (
	PromJobStashBackup  = "stash-backup"
	PromJobStashRestore = "stash-restore"
)

// ==================== RBAC related constants ==========================
const (
	KindRole        = "Role"
	KindClusterRole = "ClusterRole"

	StashBackupJobClusterRole            = "stash-backup-job"
	StashRestoreJobClusterRole           = "stash-restore-job"
	StashCronJobClusterRole              = "stash-cron-job"
	StashSidecarClusterRole              = "stash-sidecar"
	StashRestoreInitContainerClusterRole = "stash-restore-init-container"

	StashVolumeSnapshotterClusterRole      = "stash-vs-job"
	StashVolumeSnapshotRestorerClusterRole = "stash-vs-restorer-job"
	StashStorageClassReaderClusterRole     = "stash-sc-reader"
)

// ================== Condition related constants ===========================
// Condition Types Related Constants
const (
	// RepositoryFound indicates whether the respective Repository object was found or not.
	RepositoryFound = "RepositoryFound"
	// BackendSecretFound indicates whether the respective backend secret was found or not.
	BackendSecretFound = "BackendSecretFound"

	// BackupTargetFound indicates whether the backup target was found
	BackupTargetFound = "BackupTargetFound"
	// StashSidecarInjected indicates whether stash sidecar was injected into the targeted workload
	// This condition is applicable only for sidecar model
	StashSidecarInjected = "StashSidecarInjected"
	// CronJobCreated indicates whether the backup triggering CronJob was created
	CronJobCreated = "CronJobCreated"

	// RestoreTargetFound indicates whether the restore target was found
	RestoreTargetFound = "RestoreTargetFound"
	// StashInitContainerInjected indicates whether stash init-container was injected into the targeted workload
	// This condition is applicable only for sidecar model
	StashInitContainerInjected = "StashInitContainerInjected"
	// RestoreJobCreated indicates whether the restore job was created
	RestoreJobCreated = "RestoreJobCreated"

	// PreBackupHookSucceeded indicates whether PreBackupHook was executed successfully or not
	PreBackupHookSucceeded = "PreBackupHookSucceeded"
	// PostBackupHookSucceeded indicates whether PostBackupHook was executed successfully or not
	PostBackupHookSucceeded = "PostBackupHookSucceeded"
	// PreRestoreHookSucceeded indicates whether PreRestoreHook was executed successfully or not
	PreRestoreHookSucceeded = "PreRestoreHookSucceeded"
	// PostRestoreHookSucceeded indicates whether PostRestoreHook was executed successfully or not
	PostRestoreHookSucceeded = "PostRestoreHookSucceeded"
)

// Condition Reason Related Constants
const (
	// RepositoryAvailable indicates that the condition transitioned to this state because the Repository was available
	RepositoryAvailable = "RepositoryAvailable"
	// RepositoryNotAvailable indicates that the condition transitioned to this state because the Repository was not available
	RepositoryNotAvailable = "RepositoryNotAvailable"
	// UnableToCheckRepositoryAvailability indicates that the condition transitioned to this state because operator was unable
	// to check the Repository availability
	UnableToCheckRepositoryAvailability = "UnableToCheckRepositoryAvailability"

	// BackendSecretAvailable indicates that the condition transitioned to this state because the backend Secret was available
	BackendSecretAvailable = "BackendSecretAvailable"
	// BackendSecretNotAvailable indicates that the condition transitioned to this state because the backend Secret was not available
	BackendSecretNotAvailable = "BackendSecretNotAvailable"
	// UnableToCheckBackendSecretAvailability indicates that the condition transitioned to this state because operator was unable
	// to check the backend Secret availability
	UnableToCheckBackendSecretAvailability = "UnableToCheckBackendSecretAvailability"

	// TargetAvailable indicates that the condition transitioned to this state because the target was available
	TargetAvailable = "TargetAvailable"
	// TargetNotAvailable indicates that the condition transitioned to this state because the target was not available
	TargetNotAvailable = "TargetNotAvailable"
	// UnableToCheckTargetAvailability indicates that the condition transitioned to this state because operator was unable
	// to check the target availability
	UnableToCheckTargetAvailability = "UnableToCheckTargetAvailability"

	// SidecarInjectionSucceeded indicates that the condition transitioned to this state because sidecar was injected
	// successfully into the targeted workload
	SidecarInjectionSucceeded = "SidecarInjectionSucceeded"
	// SidecarInjectionFailed indicates that the condition transitioned to this state because operator was unable
	// to inject sidecar into the targeted workload
	SidecarInjectionFailed = "SidecarInjectionFailed"

	// InitContainerInjectionSucceeded indicates that the condition transitioned to this state because stash init-container
	// was injected successfully into the targeted workload
	InitContainerInjectionSucceeded = "InitContainerInjectionSucceeded"
	// InitContainerInjectionFailed indicates that the condition transitioned to this state because operator was unable
	// to inject stash init-container into the targeted workload
	InitContainerInjectionFailed = "InitContainerInjectionFailed"

	// CronJobCreationSucceeded indicates that the condition transitioned to this state because backup triggering CronJob was created successfully
	CronJobCreationSucceeded = "CronJobCreationSucceeded"
	// CronJobCreationFailed indicates that the condition transitioned to this state because operator was unable to create backup triggering CronJob
	CronJobCreationFailed = "CronJobCreationFailed"

	// RestoreJobCreationSucceeded indicates that the condition transitioned to this state because restore job was created successfully
	RestoreJobCreationSucceeded = "RestoreJobCreationSucceeded"
	// RestoreJobCreationFailed indicates that the condition transitioned to this state because operator was unable to create restore job
	RestoreJobCreationFailed = "RestoreJobCreationFailed"

	// PreBackupHookExecutedSuccessfully indicates that the condition transitioned to this state because the PreBackupHook was executed successfully
	PreBackupHookExecutedSuccessfully = "PreBackupHookExecutedSuccessfully"
	// PreBackupHookExecutionFailed indicates that the condition transitioned to this state because the Stash was unable to execute PreBackupHook
	PreBackupHookExecutionFailed = "PreBackupHookExecutionFailed"

	// PostBackupHookExecutedSuccessfully indicates that the condition transitioned to this state because the PostBackupHook was executed successfully
	PostBackupHookExecutedSuccessfully = "PostBackupHookExecutedSuccessfully"
	// PostBackupHookExecutionFailed indicates that the condition transitioned to this state because the Stash was unable to execute PostBackupHook
	PostBackupHookExecutionFailed = "PostBackupHookExecutionFailed"

	// PreRestoreHookExecutedSuccessfully indicates that the condition transitioned to this state because the PreRestoreHook was executed successfully
	PreRestoreHookExecutedSuccessfully = "PreRestoreHookExecutedSuccessfully"
	// PreRestoreHookExecutionFailed indicates that the condition transitioned to this state because the Stash was unable to execute PreRestoreHook
	PreRestoreHookExecutionFailed = "PreRestoreHookExecutionFailed"

	// PostRestoreHookExecutedSuccessfully indicates that the condition transitioned to this state because the PostRestoreHook was executed successfully
	PostRestoreHookExecutedSuccessfully = "PostRestoreHookExecutedSuccessfully"
	// PostRestoreHookExecutionFailed indicates that the condition transitioned to this state because the Stash was unable to execute PostRestoreHook
	PostRestoreHookExecutionFailed = "PostRestoreHookExecutionFailed"
)
