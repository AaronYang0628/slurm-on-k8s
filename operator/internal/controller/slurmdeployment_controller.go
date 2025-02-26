/*
Copyright 2025.

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

package controller

import (
	"context"
	"log"
	"strings"
	"time"

	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/api/errors"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/cli"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"helm.sh/helm/v3/pkg/chart/loader"

	slurmv1 "github.com/AaronYang0628/slurm-on-k8s/api/v1"
)

// SlurmDeploymentReconciler reconciles a SlurmDeployment object
type SlurmDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=slurm.ay.dev,resources=slurmdeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=slurm.ay.dev,resources=slurmdeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=slurm.ay.dev,resources=slurmdeployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SlurmDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
func GetValueWithDefault[T any](ptr *T, defaultValue T) T {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}

func buildChartValues(r *slurmv1.SlurmDeployment) map[string]interface{} {
	if r.Spec.Values.CommonAnnotations == nil {
		r.Spec.Values.CommonAnnotations = map[string]string{}
	}
	if r.Spec.Values.CommonLabels == nil {
		r.Spec.Values.CommonLabels = map[string]string{}
	}
	if r.Spec.Values.Mariadb.Auth == nil {
		r.Spec.Values.Mariadb.Auth = &slurmv1.MariaDBAuthSpec{
			Username:     "slurm",
			Password:     "password-for-slurm",
			DatabaseName: "slurm_acct_db",
		}
	}

	values := map[string]interface{}{
		"nameOverride":      r.Spec.Values.NameOverride,
		"fullnameOverride":  r.Spec.Values.FullnameOverride,
		"commonAnnotations": r.Spec.Values.CommonAnnotations,
		"commonLabels":      r.Spec.Values.CommonLabels,
		"mariadb": map[string]interface{}{
			"enabled": r.Spec.Values.Mariadb.Enabled,
			"port":    r.Spec.Values.Mariadb.Port,
			"auth": map[string]interface{}{
				"username": r.Spec.Values.Mariadb.Auth.Username,
				"password": r.Spec.Values.Mariadb.Auth.Password,
				"database": r.Spec.Values.Mariadb.Auth.DatabaseName,
			},
			"primary": map[string]interface{}{
				"persistence": map[string]interface{}{
					"enabled":      r.Spec.Values.Mariadb.Primary.Persistence.Enabled,
					"storageClass": r.Spec.Values.Mariadb.Primary.Persistence.StorageClass,
					"size":         r.Spec.Values.Mariadb.Primary.Persistence.Size,
				},
			},
		},
		"auth": map[string]interface{}{
			"ssh": map[string]interface{}{
				"secret": map[string]interface{}{
					"name": "slurm-ssh-keys",
					"keys": map[string]interface{}{
						"public":         "id_rsa.pub",
						"private":        "id_rsa",
						"authorizedKeys": "authorized_keys",
					},
				},
				"configmap": map[string]interface{}{
					"name":          r.Spec.Values.Auth.SSH.ConfigMap.Name,
					"prefabPubKeys": r.Spec.Values.Auth.SSH.ConfigMap.PrefabPubKeys,
				},
			},
		},
		"persistence": map[string]interface{}{
			"shared": map[string]interface{}{
				"enabled":       r.Spec.Values.Persistence.Shared.Enabled,
				"name":          r.Spec.Values.Persistence.Shared.Name,
				"existingClaim": r.Spec.Values.Persistence.Shared.ExistingClaim,
				"accessModes":   r.Spec.Values.Persistence.Shared.AccessModes,
				"storageClass":  r.Spec.Values.Persistence.Shared.StorageClass,
				"size":          r.Spec.Values.Persistence.Shared.Size,
			},
		},
		"resourcesPreset": r.Spec.Values.ResourcesPreset,
		"munged": map[string]interface{}{
			"name":         "munged",
			"enabled":      true,
			"commonLabels": map[string]string{},
			"image": map[string]interface{}{
				"registry":    r.Spec.Values.Munged.Image.Registry,
				"repository":  r.Spec.Values.Munged.Image.Repository,
				"tag":         r.Spec.Values.Munged.Image.Tag,
				"pullPolicy":  r.Spec.Values.Munged.Image.PullPolicy,
				"pullSecrets": r.Spec.Values.Munged.Image.PullSecrets,
			},
			"diagnosticMode": map[string]interface{}{
				"enabled": r.Spec.Values.Munged.DiagnosticMode.Enabled,
				"command": r.Spec.Values.Munged.DiagnosticMode.Command,
				"args":    r.Spec.Values.Munged.DiagnosticMode.Args,
			},
			"extraVolumes":      r.Spec.Values.Munged.ExtraVolumes,
			"extraVolumeMounts": r.Spec.Values.Munged.ExtraVolumeMounts,
		},
		"slurmctld": map[string]interface{}{
			"name":         "slurmctld",
			"commonLabels": map[string]string{},
			"replicaCount": r.Spec.Values.Slurmctld.ReplicaCount,
			"image": map[string]interface{}{
				"registry":    r.Spec.Values.Slurmctld.Image.Registry,
				"repository":  r.Spec.Values.Slurmctld.Image.Repository,
				"tag":         r.Spec.Values.Slurmctld.Image.Tag,
				"pullPolicy":  r.Spec.Values.Slurmctld.Image.PullPolicy,
				"pullSecrets": r.Spec.Values.Slurmctld.Image.PullSecrets,
			},
			"diagnosticMode": map[string]interface{}{
				"enabled": r.Spec.Values.Slurmctld.DiagnosticMode.Enabled,
				"command": r.Spec.Values.Slurmctld.DiagnosticMode.Command,
				"args":    r.Spec.Values.Slurmctld.DiagnosticMode.Args,
			},
			"automountServiceAccountToken": false,
			"podLabels":                    map[string]string{},
			"affinity":                     map[string]string{},
			"podAnnatations":               map[string]string{},
			"podAffinityPreset":            "",
			"podAntiAffinityPreset":        "soft",
			"nodeAffinityPreset": map[string]interface{}{
				"type":   "",
				"key":    "",
				"values": []string{},
			},
			"hostNetwork":               false,
			"dnsPolicy":                 "",
			"dnsConfig":                 map[string]string{},
			"hostIPC":                   false,
			"priorityClassName":         "",
			"nodeSelector":              map[string]string{},
			"tolerations":               []string{},
			"schedulerName":             "",
			"topologySpreadConstraints": []string{},
			"podSecurityContext": map[string]interface{}{
				"enabled":             true,
				"fsGroup":             0,
				"fsGroupChangePolicy": "Always",
				"supplementalGroups":  []string{},
			},
			"terminationGracePeriodSeconds": "",
			"hostAliases":                   []string{},
			"extraEnvVars":                  []string{},
			"extraEnvVarsCM":                "",
			"extraEnvVarsSecret":            "",
			"revisionHistoryLimit":          10,
			"updateStrategy": map[string]interface{}{
				"type":          "RollingUpdate",
				"rollingUpdate": map[string]string{},
			},
			"lifecycleHooks": map[string]string{},
			"resources": map[string]interface{}{
				"requests": map[string]string{
					"cpu":               r.Spec.Values.Slurmctld.Resources.Requests.CPU,
					"memory":            r.Spec.Values.Slurmctld.Resources.Requests.Memory,
					"ephemeral-storage": r.Spec.Values.Slurmctld.Resources.Requests.EphemeralStorage,
				},
				"limit": map[string]string{
					"cpu":               r.Spec.Values.Slurmctld.Resources.Limits.CPU,
					"memory":            r.Spec.Values.Slurmctld.Resources.Limits.Memory,
					"ephemeral-storage": r.Spec.Values.Slurmctld.Resources.Limits.EphemeralStorage,
				},
			},
			"extraVolumes":      r.Spec.Values.Slurmctld.ExtraVolumes,
			"extraVolumeMounts": r.Spec.Values.Slurmctld.ExtraVolumeMounts,
			"livenessProbe": map[string]interface{}{
				"enabled":             false,
				"initialDelaySeconds": 30,
				"timeoutSeconds":      5,
				"periodSeconds":       10,
				"successThreshold":    1,
				"failureThreshold":    6,
			},
			"startupProbe": map[string]interface{}{
				"enabled":             false,
				"initialDelaySeconds": 30,
				"timeoutSeconds":      5,
				"periodSeconds":       10,
				"successThreshold":    1,
				"failureThreshold":    6,
			},
			"readinessProbe": map[string]interface{}{
				"enabled":             false,
				"initialDelaySeconds": 30,
				"timeoutSeconds":      5,
				"periodSeconds":       10,
				"successThreshold":    1,
				"failureThreshold":    6,
			},
			"service": map[string]interface{}{
				"name": "slurmctld-headless",
				"ssh": map[string]interface{}{
					"type":       "ClusterIP",
					"port":       22,
					"targetPort": 22,
				},
				"slurmctld": map[string]interface{}{
					"type":       "ClusterIP",
					"port":       6817,
					"targetPort": 6817,
				},
			},
		},
		"slurmd": map[string]interface{}{
			"name":         "slurmd",
			"commonLabels": map[string]string{},
			"replicaCount": r.Spec.Values.Slurmd.ReplicaCount,
			"image": map[string]interface{}{
				"registry":    r.Spec.Values.Slurmd.Image.Registry,
				"repository":  r.Spec.Values.Slurmd.Image.Repository,
				"tag":         r.Spec.Values.Slurmd.Image.Tag,
				"pullPolicy":  r.Spec.Values.Slurmd.Image.PullPolicy,
				"pullSecrets": r.Spec.Values.Slurmd.Image.PullSecrets,
			},
			"diagnosticMode": map[string]interface{}{
				"enabled": r.Spec.Values.Slurmd.DiagnosticMode.Enabled,
				"command": r.Spec.Values.Slurmd.DiagnosticMode.Command,
				"args":    r.Spec.Values.Slurmd.DiagnosticMode.Args,
			},
			"automountServiceAccountToken": false,
			"podLabels":                    map[string]string{},
			"affinity":                     map[string]string{},
			"podAnnatations":               map[string]string{},
			"podAffinityPreset":            "",
			"podAntiAffinityPreset":        "soft",
			"nodeAffinityPreset": map[string]interface{}{
				"type":   "",
				"key":    "",
				"values": []string{},
			},
			"hostNetwork":               false,
			"dnsPolicy":                 "",
			"dnsConfig":                 map[string]string{},
			"hostIPC":                   false,
			"priorityClassName":         "",
			"nodeSelector":              map[string]string{},
			"tolerations":               []string{},
			"schedulerName":             "",
			"topologySpreadConstraints": []string{},
			"podSecurityContext": map[string]interface{}{
				"enabled":             true,
				"fsGroup":             0,
				"fsGroupChangePolicy": "Always",
				"supplementalGroups":  []string{},
			},
			"terminationGracePeriodSeconds": "",
			"hostAliases":                   []string{},
			"extraEnvVars":                  []string{},
			"extraEnvVarsCM":                "",
			"extraEnvVarsSecret":            "",
			"revisionHistoryLimit":          10,
			"updateStrategy": map[string]interface{}{
				"type":          "RollingUpdate",
				"rollingUpdate": map[string]string{},
			},
			"lifecycleHooks": map[string]string{},
			"resources": map[string]interface{}{
				"requests": map[string]string{
					"cpu":               r.Spec.Values.Slurmd.Resources.Requests.CPU,
					"memory":            r.Spec.Values.Slurmd.Resources.Requests.Memory,
					"ephemeral-storage": r.Spec.Values.Slurmd.Resources.Requests.EphemeralStorage,
				},
				"limits": map[string]string{
					"cpu":               r.Spec.Values.Slurmd.Resources.Limits.CPU,
					"memory":            r.Spec.Values.Slurmd.Resources.Limits.Memory,
					"ephemeral-storage": r.Spec.Values.Slurmd.Resources.Limits.EphemeralStorage,
				},
			},
			"extraVolumes":      r.Spec.Values.Slurmd.ExtraVolumes,
			"extraVolumeMounts": r.Spec.Values.Slurmd.ExtraVolumeMounts,
			"livenessProbe": map[string]interface{}{
				"enabled":             false,
				"initialDelaySeconds": 30,
				"timeoutSeconds":      5,
				"periodSeconds":       10,
				"successThreshold":    1,
				"failureThreshold":    6,
			},
			"startupProbe": map[string]interface{}{
				"enabled":             false,
				"initialDelaySeconds": 30,
				"timeoutSeconds":      5,
				"periodSeconds":       10,
				"successThreshold":    1,
				"failureThreshold":    6,
			},
			"readinessProbe": map[string]interface{}{
				"enabled":             false,
				"initialDelaySeconds": 30,
				"timeoutSeconds":      5,
				"periodSeconds":       10,
				"successThreshold":    1,
				"failureThreshold":    6,
			},
			"service": map[string]interface{}{
				"name": "slurmd-headless",
				"ssh": map[string]interface{}{
					"type":       "ClusterIP",
					"port":       22,
					"targetPort": 22,
				},
				"slurmd": map[string]interface{}{
					"type":       "ClusterIP",
					"port":       6818,
					"targetPort": 6818,
				},
			},
		},
		"slurmdbd": map[string]interface{}{
			"name":         "slurmdbd",
			"commonLabels": map[string]string{},
			"replicaCount": 1,
			"image": map[string]interface{}{
				"registry":    r.Spec.Values.Slurmdbd.Image.Registry,
				"repository":  r.Spec.Values.Slurmdbd.Image.Repository,
				"tag":         r.Spec.Values.Slurmdbd.Image.Tag,
				"pullPolicy":  r.Spec.Values.Slurmdbd.Image.PullPolicy,
				"pullSecrets": r.Spec.Values.Slurmdbd.Image.PullSecrets,
			},
			"diagnosticMode": map[string]interface{}{
				"enabled": r.Spec.Values.Slurmdbd.DiagnosticMode.Enabled,
				"command": r.Spec.Values.Slurmdbd.DiagnosticMode.Command,
				"args":    r.Spec.Values.Slurmdbd.DiagnosticMode.Args,
			},
			"automountServiceAccountToken": false,
			"podLabels":                    map[string]string{},
			"affinity":                     map[string]string{},
			"podAnnatations":               map[string]string{},
			"podAffinityPreset":            "",
			"podAntiAffinityPreset":        "soft",
			"nodeAffinityPreset": map[string]interface{}{
				"type":   "",
				"key":    "",
				"values": []string{},
			},
			"hostNetwork":               false,
			"dnsPolicy":                 "",
			"dnsConfig":                 map[string]string{},
			"hostIPC":                   false,
			"priorityClassName":         "",
			"nodeSelector":              map[string]string{},
			"tolerations":               []string{},
			"schedulerName":             "",
			"topologySpreadConstraints": []string{},
			"podSecurityContext": map[string]interface{}{
				"enabled":             true,
				"fsGroup":             0,
				"fsGroupChangePolicy": "Always",
				"supplementalGroups":  []string{},
			},
			"terminationGracePeriodSeconds": "",
			"hostAliases":                   []string{},
			"extraEnvVars":                  []string{},
			"extraEnvVarsCM":                "",
			"extraEnvVarsSecret":            "",
			"revisionHistoryLimit":          10,
			"updateStrategy": map[string]interface{}{
				"type":          "RollingUpdate",
				"rollingUpdate": map[string]string{},
			},
			"lifecycleHooks":    map[string]string{},
			"extraVolumes":      r.Spec.Values.Slurmdbd.ExtraVolumes,
			"extraVolumeMounts": r.Spec.Values.Slurmdbd.ExtraVolumeMounts,
			"livenessProbe": map[string]interface{}{
				"enabled":             false,
				"initialDelaySeconds": 30,
				"timeoutSeconds":      5,
				"periodSeconds":       10,
				"successThreshold":    1,
				"failureThreshold":    6,
			},
			"startupProbe": map[string]interface{}{
				"enabled":             false,
				"initialDelaySeconds": 30,
				"timeoutSeconds":      5,
				"periodSeconds":       10,
				"successThreshold":    1,
				"failureThreshold":    6,
			},
			"readinessProbe": map[string]interface{}{
				"enabled":             false,
				"initialDelaySeconds": 30,
				"timeoutSeconds":      5,
				"periodSeconds":       10,
				"successThreshold":    1,
				"failureThreshold":    6,
			},
			"service": map[string]interface{}{
				"name": "slurmdbd-headless",
				"ssh": map[string]interface{}{
					"type":       "ClusterIP",
					"port":       22,
					"targetPort": 22,
				},
				"slurmdbd": map[string]interface{}{
					"type":       "ClusterIP",
					"port":       6819,
					"targetPort": 6819,
				},
			},
		},
		"login": map[string]interface{}{
			"name":         "login",
			"commonLabels": map[string]string{},
			"replicaCount": 1,
			"image": map[string]interface{}{
				"registry":    r.Spec.Values.SlurmLogin.Image.Registry,
				"repository":  r.Spec.Values.SlurmLogin.Image.Repository,
				"tag":         r.Spec.Values.SlurmLogin.Image.Tag,
				"pullPolicy":  r.Spec.Values.SlurmLogin.Image.PullPolicy,
				"pullSecrets": r.Spec.Values.SlurmLogin.Image.PullSecrets,
			},
			"diagnosticMode": map[string]interface{}{
				"enabled": r.Spec.Values.SlurmLogin.DiagnosticMode.Enabled,
				"command": r.Spec.Values.SlurmLogin.DiagnosticMode.Command,
				"args":    r.Spec.Values.SlurmLogin.DiagnosticMode.Args,
			},
			"automountServiceAccountToken": false,
			"podLabels":                    map[string]string{},
			"affinity":                     map[string]string{},
			"podAnnatations":               map[string]string{},
			"podAffinityPreset":            "",
			"podAntiAffinityPreset":        "soft",
			"nodeAffinityPreset": map[string]interface{}{
				"type":   "",
				"key":    "",
				"values": []string{},
			},
			"hostNetwork":               false,
			"dnsPolicy":                 "",
			"dnsConfig":                 map[string]string{},
			"hostIPC":                   false,
			"priorityClassName":         "",
			"nodeSelector":              map[string]string{},
			"tolerations":               []string{},
			"schedulerName":             "",
			"topologySpreadConstraints": []string{},
			"podSecurityContext": map[string]interface{}{
				"enabled":             true,
				"fsGroup":             0,
				"fsGroupChangePolicy": "Always",
				"supplementalGroups":  []string{},
			},
			"terminationGracePeriodSeconds": "",
			"hostAliases":                   []string{},
			"extraEnvVars":                  []string{},
			"extraEnvVarsCM":                "",
			"extraEnvVarsSecret":            "",
			"revisionHistoryLimit":          10,
			"updateStrategy": map[string]interface{}{
				"type":          "RollingUpdate",
				"rollingUpdate": map[string]string{},
			},
			"lifecycleHooks": map[string]string{},
			"resources": map[string]interface{}{
				"requests": map[string]string{
					"cpu":               r.Spec.Values.SlurmLogin.Resources.Requests.CPU,
					"memory":            r.Spec.Values.SlurmLogin.Resources.Requests.Memory,
					"ephemeral-storage": r.Spec.Values.SlurmLogin.Resources.Requests.EphemeralStorage,
				},
				"limits": map[string]string{
					"cpu":               r.Spec.Values.SlurmLogin.Resources.Limits.CPU,
					"memory":            r.Spec.Values.SlurmLogin.Resources.Limits.Memory,
					"ephemeral-storage": r.Spec.Values.SlurmLogin.Resources.Limits.EphemeralStorage,
				},
			},
			"extraVolumes":      r.Spec.Values.SlurmLogin.ExtraVolumes,
			"extraVolumeMounts": r.Spec.Values.SlurmLogin.ExtraVolumeMounts,
			"livenessProbe": map[string]interface{}{
				"enabled":             false,
				"initialDelaySeconds": 30,
				"timeoutSeconds":      5,
				"periodSeconds":       10,
				"successThreshold":    1,
				"failureThreshold":    6,
			},
			"startupProbe": map[string]interface{}{
				"enabled":             false,
				"initialDelaySeconds": 30,
				"timeoutSeconds":      5,
				"periodSeconds":       10,
				"successThreshold":    1,
				"failureThreshold":    6,
			},
			"readinessProbe": map[string]interface{}{
				"enabled":             false,
				"initialDelaySeconds": 30,
				"timeoutSeconds":      5,
				"periodSeconds":       10,
				"successThreshold":    1,
				"failureThreshold":    6,
			},
			"service": map[string]interface{}{
				"name": "login",
				"ssh": map[string]interface{}{
					"type":       "ClusterIP",
					"port":       22,
					"targetPort": 22,
				},
			},
		},
		"serviceAccount": map[string]interface{}{
			"automount":   true,
			"annotations": map[string]string{},
			"name":        "slurm",
			"role": map[string]string{
				"name": "slurm",
			},
			"roleBinding": map[string]string{
				"name": "slurm",
			},
		},
		"configuration": map[string]interface{}{
			"cgroup": map[string]interface{}{
				"name": "cgroup-conf",
				"value": `ConstrainCores=yes
ConstrainDevices=yes
ConstrainRAMSpace=yes
ConstrainSwapSpace=no`,
			},
			"slurmConf": `ClusterName=slurm-cluster
SlurmctldHost={{ include "common.names.fullname" . }}-{{ .Values.slurmctld.name }}-0
MpiDefault=pmi2
ProctrackType=proctrack/cgroup
ReturnToService=1
SlurmctldPidFile=/var/run/slurmctld.pid
SlurmctldPort={{ .Values.slurmctld.service.slurmctld.port }}
SlurmdPidFile=/var/run/slurmd.pid
SlurmdPort={{ .Values.slurmd.service.slurmd.port }}
SlurmdSpoolDir=/var/spool/slurmd
SlurmUser=slurm
StateSaveLocation=/var/spool/slurmctld
TaskPlugin=task/affinity,task/cgroup
InactiveLimit=0
KillWait=30
MinJobAge=300
SlurmctldTimeout=120
SlurmdTimeout=300
Waittime=0
SchedulerType=sched/backfill
SelectType=select/cons_tres
AccountingStorageHost={{ include "common.names.fullname" . }}-{{ .Values.slurmdbd.name }}-0
AccountingStoragePort={{ .Values.slurmdbd.service.slurmdbd.port }}
AccountingStorageType=accounting_storage/slurmdbd
AccountingStoreFlags=job_comment
JobAcctGatherType=jobacct_gather/linux
JobAcctGatherFrequency=30
SlurmctldDebug=info
SlurmctldLogFile=/var/log/slurm/slurmctld.log
SlurmdDebug=info
SlurmdLogFile=/var/log/slurm/slurmd.log
NodeName={{ include "common.names.fullname" . }}-slurmd-[0-999] CPUs={{ .Values.slurmd.resources.requests.cpu }} CoresPerSocket=6 ThreadsPerCore=1 RealMemory=1024 Procs=1 State=UNKNOWN
PartitionName=compute Nodes=ALL Default=YES MaxTime=INFINITE State=UP`,
			"slurmdbdConf": `AuthType=auth/munge
AuthInfo=/var/run/munge/munge.socket.2
SlurmUser=slurm
DebugLevel=verbose
LogFile=/var/log/slurm/slurmdbd.log
PidFile=/var/run/slurmdbd.pid
DbdHost={{ include "common.names.fullname" . }}-{{ .Values.slurmdbd.name }}-0
DbdPort={{ .Values.slurmdbd.service.slurmdbd.port }}
StorageType=accounting_storage/mysql
StorageHost={{ include "common.names.fullname" . }}-mariadb
StoragePort={{ .Values.mariadb.port }}
StoragePass={{ .Values.mariadb.auth.password }}
StorageUser={{ .Values.mariadb.auth.username }}
StorageLoc={{ .Values.mariadb.auth.database }}`,
		},
	}
	return values
}

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.0/pkg/reconcile
func (r *SlurmDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// get CR SlurmDeployment instance
	release := &slurmv1.SlurmDeployment{}
	if err := r.Get(ctx, req.NamespacedName, release); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Printf("Find SlurmDeployment %s", release.Name)

	// Check if namespace exists, if not, create it
	namespace := release.Spec.Chart.Namespace
	if namespace != "" {
		ns := &corev1.Namespace{}
		err := r.Get(ctx, client.ObjectKey{Name: namespace}, ns)
		if errors.IsNotFound(err) {
			// Namespace does not exist, create it
			ns = &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			if err := r.Create(ctx, ns); err != nil {
				log.Printf("Failed to create namespace %s: %v", namespace, err)
				return ctrl.Result{}, err
			}
			log.Printf("Namespace %s created successfully", namespace)
		} else if err != nil {
			log.Printf("Failed to get namespace %s: %v", namespace, err)
			return ctrl.Result{}, err
		}
	}

	// init Slurm Helm Chart settings
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), release.Spec.Chart.Namespace, "secret", log.Printf); err != nil {
		log.Printf("Failed to initialize Helm configuration: %v", err)
		return ctrl.Result{}, err
	}

	// build values yaml content for Slurm Chart
	values := buildChartValues(release)

	// Check release if exists
	histClient := action.NewHistory(actionConfig)
	if _, err := histClient.Run(release.Name); err == nil {
		// upgrade release
		upgradeClient := action.NewUpgrade(actionConfig)
		upgradeClient.Namespace = release.Spec.Chart.Namespace

		_, err = upgradeClient.Run(release.Name, getChart(release), values)
		return handleResult(err)
	}

	// install a new release
	installClient := action.NewInstall(actionConfig)
	installClient.ReleaseName = release.Name
	installClient.Namespace = release.Spec.Chart.Namespace
	_, err := installClient.Run(getChart(release), values)
	return handleResult(err)
}

func handleResult(err error) (ctrl.Result, error) {
	if err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func getChart(r *slurmv1.SlurmDeployment) *chart.Chart {
	// 实现 Chart 下载逻辑（从仓库获取）
	// 可以使用 helm pull 命令或直接访问仓库
	chartName := r.Spec.Chart.Name
	repository := r.Spec.Chart.Repository
	version := r.Spec.Chart.Version

	// 构造 Chart 的下载 URL
	chartURL := fmt.Sprintf("%s/%s-%s.tgz", repository, chartName, version)
	log.Printf("Downloading chart from %s", chartURL)

	// 创建临时目录用于存储下载的 Chart 文件
	tempDir, err := os.MkdirTemp("", "helm-charts")
	if err != nil {
		log.Printf("Failed to create temporary directory: %v", err)
		return nil
	}
	defer os.RemoveAll(tempDir)

	// 下载 Chart 文件
	filePath := filepath.Join(tempDir, fmt.Sprintf("%s-%s.tgz", chartName, version))
	if err := downloadFile(chartURL, filePath); err != nil {
		log.Printf("Failed to download chart: %v", err)
		return nil
	} else {
		log.Printf("Downloaded chart to %s", filePath)
	}

	// 解压 Chart 文件
	if err := extractTarGz(filePath, tempDir); err != nil {
		log.Printf("Failed to extract chart: %v", err)
		return nil
	} else {
		log.Printf("Extracted chart to %s", tempDir)
	}

	// 加载 Chart
	chartPath := filepath.Join(tempDir, chartName)
	chrt, err := loader.Load(chartPath)
	if err != nil {
		log.Printf("Failed to load chart: %v", err)
		return nil
	} else {
		log.Printf("Loaded chart %s", chrt.Metadata.Name)
	}

	return chrt
}

// download chart File
func downloadFile(url, filePath string) error {
	log.Printf("Downloading file from %s to %s", url, filePath)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error downloading file: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to download file: %s", resp.Status)
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	out, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return err
	}
	defer out.Close()

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		log.Printf("Error writing file: %v", err)
		return err
	}
	log.Printf("Downloaded %d bytes to %s", written, filePath)

	return nil
}

// unzip helm chart.tgz file
func extractTarGz(filePath, destPath string) error {
	// 打开压缩文件
	gzFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer gzFile.Close()

	// 创建gzip reader
	gzReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return fmt.Errorf("创建gzip reader失败: %v", err)
	}
	defer gzReader.Close()

	// 创建tar reader
	tarReader := tar.NewReader(gzReader)

	// 获取目标路径绝对地址用于安全检查
	destAbs, err := filepath.Abs(destPath)
	if err != nil {
		return fmt.Errorf("获取绝对路径失败: %v", err)
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // 正常结束
		}
		if err != nil {
			return fmt.Errorf("读取tar条目失败: %v", err)
		}

		// 构建目标路径并进行安全检查
		targetPath := filepath.Join(destPath, header.Name)
		targetAbs, err := filepath.Abs(targetPath)
		if err != nil {
			return fmt.Errorf("路径解析失败: %v", err)
		}

		// 防止路径穿越攻击
		if !strings.HasPrefix(targetAbs, destAbs) {
			return fmt.Errorf("危险路径检测: %s 试图访问目标路径外", header.Name)
		}

		// 根据文件类型处理
		switch header.Typeflag {
		case tar.TypeDir:
			// 创建目录并设置权限
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("创建目录失败: %v", err)
			}
			if err := os.Chmod(targetPath, os.FileMode(header.Mode).Perm()); err != nil {
				return fmt.Errorf("设置目录权限失败: %v", err)
			}

		case tar.TypeReg:
			// 确保父目录存在
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("创建父目录失败: %v", err)
			}

			// 创建文件并设置内容
			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode).Perm())
			if err != nil {
				return fmt.Errorf("创建文件失败: %v", err)
			}

			// 复制文件内容
			if _, err := io.CopyN(outFile, tarReader, header.Size); err != nil {
				outFile.Close()
				return fmt.Errorf("写入文件内容失败: %v", err)
			}

			// 设置文件修改时间
			if err := os.Chtimes(targetPath, time.Time{}, header.ModTime); err != nil {
				outFile.Close()
				return fmt.Errorf("设置修改时间失败: %v", err)
			}

			outFile.Close()

		default:
			// 跳过非常规文件类型（如符号链接、设备文件等）
			continue
		}
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SlurmDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&slurmv1.SlurmDeployment{}).
		Complete(r)
}
