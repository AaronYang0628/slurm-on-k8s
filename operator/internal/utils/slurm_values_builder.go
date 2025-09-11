package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	slurmv1 "github.com/AaronYang0628/slurm-on-k8s/api/v1"
)

func GetLocalCPUInfo(prefix string) (int32, error) {
	// open /proc/cpuinfo
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return 0, fmt.Errorf("Cannot read /proc/cpuinfo: %v", err)
	}
	defer file.Close()
	physicalIDs := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, prefix) {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				physicalID := strings.TrimSpace(parts[1])
				physicalIDs[physicalID] = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("Failed to read /proc/cpuinfo: %v", err)
	}

	return int32(len(physicalIDs)), nil
}

func BuildSlurmValues(valuesSpec *slurmv1.ValuesSpec) map[string]interface{} {
	if valuesSpec.CommonAnnotations == nil {
		valuesSpec.CommonAnnotations = map[string]string{}
	}
	if valuesSpec.CommonLabels == nil {
		valuesSpec.CommonLabels = map[string]string{}
	}
	if valuesSpec.Mariadb.Auth == nil {
		valuesSpec.Mariadb.Auth = &slurmv1.MariaDBAuthSpec{
			RootPassword: "password-for-root",
			Username:     "slurm",
			Password:     "password-for-slurm",
			DatabaseName: "slurm_acct_db",
		}
	}

	if valuesSpec.Slurmctld.Resources == nil {
		valuesSpec.Slurmctld.Resources = &slurmv1.ResourceSpec{
			Requests: &slurmv1.ResourceRequestSpec{
				CPU:              "1000m",
				Memory:           "1Gi",
				EphemeralStorage: "10Gi",
			},
			Limits: &slurmv1.ResourceLimitSpec{
				CPU:              "2000m",
				Memory:           "2Gi",
				EphemeralStorage: "20Gi",
			},
		}
	}

	if valuesSpec.SlurmdCPU.Resources.Limits == nil {
		sockets, err := GetLocalCPUInfo("physical id")
		if err != nil {
			sockets = 1
		}

		cores, err2 := GetLocalCPUInfo("cpu cores")
		if err2 != nil {
			cores = 1
		}

		valuesSpec.SlurmdCPU.Resources.Limits = &slurmv1.SlurmdResourceLimitSpec{
			Socket:           sockets,
			CorePerSocket:    cores,
			ThreadPerCore:    1,
			Memory:           "8Gi",
			EphemeralStorage: "20Gi",
		}
	}

	if valuesSpec.SlurmdCPU.Resources.Requests == nil {
		sockets, err := GetLocalCPUInfo("physical id")
		if err != nil {
			sockets = 1
		}

		cores, err2 := GetLocalCPUInfo("cpu cores")
		if err2 != nil {
			cores = 1
		}

		valuesSpec.SlurmdCPU.Resources.Requests = &slurmv1.SlurmdResourceRequestSpec{
			Socket:           sockets,
			CorePerSocket:    cores,
			ThreadPerCore:    1,
			Memory:           "1Gi",
			EphemeralStorage: "2Gi",
		}
	}

	if valuesSpec.SlurmdGPU.Resources.Limits == nil {
		sockets, err := GetLocalCPUInfo("physical id")
		if err != nil {
			sockets = 1
		}

		cores, err2 := GetLocalCPUInfo("cpu cores")
		if err2 != nil {
			cores = 1
		}

		valuesSpec.SlurmdGPU.Resources.Limits = &slurmv1.SlurmdResourceLimitSpec{
			Socket:           sockets,
			CorePerSocket:    cores,
			ThreadPerCore:    1,
			Memory:           "8Gi",
			EphemeralStorage: "20Gi",
			GPUShared:        0,
		}
	}

	if valuesSpec.SlurmdGPU.Resources.Requests == nil {
		sockets, err := GetLocalCPUInfo("physical id")
		if err != nil {
			sockets = 1
		}

		cores, err2 := GetLocalCPUInfo("cpu cores")
		if err2 != nil {
			cores = 1
		}

		valuesSpec.SlurmdGPU.Resources.Requests = &slurmv1.SlurmdResourceRequestSpec{
			Socket:           sockets,
			CorePerSocket:    cores,
			ThreadPerCore:    1,
			Memory:           "1Gi",
			EphemeralStorage: "2Gi",
			GPUShared:        0,
		}
	}

	if valuesSpec.SlurmLogin.Resources.Limits == nil {
		valuesSpec.SlurmLogin.Resources.Limits = &slurmv1.ResourceLimitSpec{
			CPU:              "2000m",
			Memory:           "8Gi",
			EphemeralStorage: "20Gi",
		}
	}

	if valuesSpec.SlurmLogin.Resources.Requests == nil {
		valuesSpec.SlurmLogin.Resources.Requests = &slurmv1.ResourceRequestSpec{
			CPU:              "1000m",
			Memory:           "1Gi",
			EphemeralStorage: "2Gi",
		}
	}

	values := map[string]interface{}{
		"nameOverride":      valuesSpec.NameOverride,
		"fullnameOverride":  valuesSpec.FullnameOverride,
		"commonAnnotations": valuesSpec.CommonAnnotations,
		"commonLabels":      valuesSpec.CommonLabels,
		"image": map[string]interface{}{
			"mirror": map[string]string{
				"registry": valuesSpec.ImageMirror.Mirror.Registry,
			},
		},
		"mariadb": map[string]interface{}{
			"enabled": valuesSpec.Mariadb.Enabled,
			"port":    valuesSpec.Mariadb.Port,
			"auth": map[string]interface{}{
				"rootPassword": valuesSpec.Mariadb.Auth.RootPassword,
				"username":     valuesSpec.Mariadb.Auth.Username,
				"password":     valuesSpec.Mariadb.Auth.Password,
				"database":     valuesSpec.Mariadb.Auth.DatabaseName,
			},
			"primary": map[string]interface{}{
				"persistence": map[string]interface{}{
					"enabled":      valuesSpec.Mariadb.Primary.Persistence.Enabled,
					"storageClass": valuesSpec.Mariadb.Primary.Persistence.StorageClass,
					"size":         valuesSpec.Mariadb.Primary.Persistence.Size,
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
					"name":          valuesSpec.Auth.SSH.ConfigMap.Name,
					"prefabPubKeys": valuesSpec.Auth.SSH.ConfigMap.PrefabPubKeys,
				},
			},
		},
		"persistence": map[string]interface{}{
			"shared": map[string]interface{}{
				"enabled":       valuesSpec.Persistence.Shared.Enabled,
				"name":          valuesSpec.Persistence.Shared.Name,
				"existingClaim": valuesSpec.Persistence.Shared.ExistingClaim,
				"accessModes":   valuesSpec.Persistence.Shared.AccessModes,
				"storageClass":  valuesSpec.Persistence.Shared.StorageClass,
				"size":          valuesSpec.Persistence.Shared.Size,
			},
		},
		"resourcesPreset": valuesSpec.ResourcesPreset,
		"munged": map[string]interface{}{
			"name":         "munged",
			"commonLabels": map[string]string{},
			"image": map[string]interface{}{
				"registry":    valuesSpec.Munged.Image.Registry,
				"repository":  valuesSpec.Munged.Image.Repository,
				"tag":         valuesSpec.Munged.Image.Tag,
				"pullPolicy":  valuesSpec.Munged.Image.PullPolicy,
				"pullSecrets": valuesSpec.Munged.Image.PullSecrets,
			},
			"diagnosticMode": map[string]interface{}{
				"enabled": valuesSpec.Munged.DiagnosticMode.Enabled,
				"command": valuesSpec.Munged.DiagnosticMode.Command,
				"args":    valuesSpec.Munged.DiagnosticMode.Args,
			},
			"extraVolumes":      valuesSpec.Munged.ExtraVolumes,
			"extraVolumeMounts": valuesSpec.Munged.ExtraVolumeMounts,
		},
		"slurmctld": map[string]interface{}{
			"name":         "slurmctld",
			"commonLabels": map[string]string{},
			"replicaCount": valuesSpec.Slurmctld.ReplicaCount,
			"image": map[string]interface{}{
				"registry":    valuesSpec.Slurmctld.Image.Registry,
				"repository":  valuesSpec.Slurmctld.Image.Repository,
				"tag":         valuesSpec.Slurmctld.Image.Tag,
				"pullPolicy":  valuesSpec.Slurmctld.Image.PullPolicy,
				"pullSecrets": valuesSpec.Slurmctld.Image.PullSecrets,
			},
			"diagnosticMode": map[string]interface{}{
				"enabled": valuesSpec.Slurmctld.DiagnosticMode.Enabled,
				"command": valuesSpec.Slurmctld.DiagnosticMode.Command,
				"args":    valuesSpec.Slurmctld.DiagnosticMode.Args,
			},
			"automountServiceAccountToken": false,
			"podLabels":                    map[string]string{},
			"podAnnatations":               map[string]string{},
			"podAffinityPreset":            "",
			"podAntiAffinityPreset":        "soft",
			"nodeAffinityPreset": map[string]interface{}{
				"type":   valuesSpec.Slurmctld.NodeAffinityPreset.Type,
				"key":    valuesSpec.Slurmctld.NodeAffinityPreset.Key,
				"values": valuesSpec.Slurmctld.NodeAffinityPreset.Values,
				"weight": valuesSpec.Slurmctld.NodeAffinityPreset.Weight,
			},
			"hostNetwork":               false,
			"dnsPolicy":                 "",
			"dnsConfig":                 map[string]string{},
			"hostIPC":                   false,
			"priorityClassName":         "",
			"nodeSelector":              valuesSpec.Slurmctld.NodeSelector,
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
					"cpu":               valuesSpec.Slurmctld.Resources.Requests.CPU,
					"memory":            valuesSpec.Slurmctld.Resources.Requests.Memory,
					"ephemeral-storage": valuesSpec.Slurmctld.Resources.Requests.EphemeralStorage,
				},
				"limits": map[string]string{
					"cpu":               valuesSpec.Slurmctld.Resources.Limits.CPU,
					"memory":            valuesSpec.Slurmctld.Resources.Limits.Memory,
					"ephemeral-storage": valuesSpec.Slurmctld.Resources.Limits.EphemeralStorage,
				},
			},
			"extraVolumes":      valuesSpec.Slurmctld.ExtraVolumes,
			"extraVolumeMounts": valuesSpec.Slurmctld.ExtraVolumeMounts,
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
		"slurmdCPU": map[string]interface{}{
			"name":         "slurmd",
			"commonLabels": map[string]string{},
			"replicaCount": valuesSpec.SlurmdCPU.ReplicaCount,
			"image": map[string]interface{}{
				"registry":    valuesSpec.SlurmdCPU.Image.Registry,
				"repository":  valuesSpec.SlurmdCPU.Image.Repository,
				"tag":         valuesSpec.SlurmdCPU.Image.Tag,
				"pullPolicy":  valuesSpec.SlurmdCPU.Image.PullPolicy,
				"pullSecrets": valuesSpec.SlurmdCPU.Image.PullSecrets,
			},
			"diagnosticMode": map[string]interface{}{
				"enabled": valuesSpec.SlurmdCPU.DiagnosticMode.Enabled,
				"command": valuesSpec.SlurmdCPU.DiagnosticMode.Command,
				"args":    valuesSpec.SlurmdCPU.DiagnosticMode.Args,
			},
			"automountServiceAccountToken": false,
			"podLabels":                    map[string]string{},
			"podAnnatations":               map[string]string{},
			"podAffinityPreset":            "",
			"podAntiAffinityPreset":        "soft",
			"nodeAffinityPreset": map[string]interface{}{
				"type":   valuesSpec.SlurmdCPU.NodeAffinityPreset.Type,
				"key":    valuesSpec.SlurmdCPU.NodeAffinityPreset.Key,
				"values": valuesSpec.SlurmdCPU.NodeAffinityPreset.Values,
				"weight": valuesSpec.SlurmdCPU.NodeAffinityPreset.Weight,
			},
			"hostNetwork":               false,
			"dnsPolicy":                 "",
			"dnsConfig":                 map[string]string{},
			"hostIPC":                   false,
			"priorityClassName":         "",
			"nodeSelector":              valuesSpec.SlurmdCPU.NodeSelector,
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
					"cpu":               fmt.Sprintf("%dm", valuesSpec.SlurmdCPU.Resources.Requests.Socket*valuesSpec.SlurmdCPU.Resources.Requests.CorePerSocket*valuesSpec.SlurmdCPU.Resources.Requests.ThreadPerCore*1000),
					"memory":            valuesSpec.SlurmdCPU.Resources.Requests.Memory,
					"ephemeral-storage": valuesSpec.SlurmdCPU.Resources.Requests.EphemeralStorage,
				},
				"limits": map[string]string{
					"cpu":               fmt.Sprintf("%dm", valuesSpec.SlurmdCPU.Resources.Limits.Socket*valuesSpec.SlurmdCPU.Resources.Limits.CorePerSocket*valuesSpec.SlurmdCPU.Resources.Limits.ThreadPerCore*1000),
					"memory":            valuesSpec.SlurmdCPU.Resources.Limits.Memory,
					"ephemeral-storage": valuesSpec.SlurmdCPU.Resources.Limits.EphemeralStorage,
				},
			},
			"extraVolumes":      valuesSpec.SlurmdCPU.ExtraVolumes,
			"extraVolumeMounts": valuesSpec.SlurmdCPU.ExtraVolumeMounts,
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
				"name": "slurmd-cpu-headless",
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
		"slurmdGPU": map[string]interface{}{
			"name":         "slurmd",
			"commonLabels": map[string]string{},
			"replicaCount": valuesSpec.SlurmdGPU.ReplicaCount,
			"image": map[string]interface{}{
				"registry":    valuesSpec.SlurmdGPU.Image.Registry,
				"repository":  valuesSpec.SlurmdGPU.Image.Repository,
				"tag":         valuesSpec.SlurmdGPU.Image.Tag,
				"pullPolicy":  valuesSpec.SlurmdGPU.Image.PullPolicy,
				"pullSecrets": valuesSpec.SlurmdGPU.Image.PullSecrets,
			},
			"diagnosticMode": map[string]interface{}{
				"enabled": valuesSpec.SlurmdGPU.DiagnosticMode.Enabled,
				"command": valuesSpec.SlurmdGPU.DiagnosticMode.Command,
				"args":    valuesSpec.SlurmdGPU.DiagnosticMode.Args,
			},
			"automountServiceAccountToken": false,
			"podLabels":                    map[string]string{},
			"podAnnatations":               map[string]string{},
			"podAffinityPreset":            "",
			"podAntiAffinityPreset":        "soft",
			"nodeAffinityPreset": map[string]interface{}{
				"type":   valuesSpec.SlurmdGPU.NodeAffinityPreset.Type,
				"key":    valuesSpec.SlurmdGPU.NodeAffinityPreset.Key,
				"values": valuesSpec.SlurmdGPU.NodeAffinityPreset.Values,
				"weight": valuesSpec.SlurmdGPU.NodeAffinityPreset.Weight,
			},
			"hostNetwork":               false,
			"dnsPolicy":                 "",
			"dnsConfig":                 map[string]string{},
			"hostIPC":                   false,
			"priorityClassName":         "",
			"nodeSelector":              valuesSpec.SlurmdGPU.NodeSelector,
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
					"cpu":                   fmt.Sprintf("%dm", valuesSpec.SlurmdGPU.Resources.Requests.Socket*valuesSpec.SlurmdGPU.Resources.Requests.CorePerSocket*valuesSpec.SlurmdGPU.Resources.Requests.ThreadPerCore*1000),
					"memory":                valuesSpec.SlurmdGPU.Resources.Requests.Memory,
					"ephemeral-storage":     valuesSpec.SlurmdGPU.Resources.Requests.EphemeralStorage,
					"nvidia.com/gpu.shared": fmt.Sprintf("%d", valuesSpec.SlurmdGPU.Resources.Requests.GPUShared),
				},
				"limits": map[string]string{
					"cpu":                   fmt.Sprintf("%dm", valuesSpec.SlurmdGPU.Resources.Limits.Socket*valuesSpec.SlurmdGPU.Resources.Limits.CorePerSocket*valuesSpec.SlurmdGPU.Resources.Limits.ThreadPerCore*1000),
					"memory":                valuesSpec.SlurmdGPU.Resources.Limits.Memory,
					"ephemeral-storage":     valuesSpec.SlurmdGPU.Resources.Limits.EphemeralStorage,
					"nvidia.com/gpu.shared": fmt.Sprintf("%d", valuesSpec.SlurmdGPU.Resources.Limits.GPUShared),
				},
			},
			"extraVolumes":      valuesSpec.SlurmdGPU.ExtraVolumes,
			"extraVolumeMounts": valuesSpec.SlurmdGPU.ExtraVolumeMounts,
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
				"name": "slurmd-gpu-headless",
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
				"registry":    valuesSpec.Slurmdbd.Image.Registry,
				"repository":  valuesSpec.Slurmdbd.Image.Repository,
				"tag":         valuesSpec.Slurmdbd.Image.Tag,
				"pullPolicy":  valuesSpec.Slurmdbd.Image.PullPolicy,
				"pullSecrets": valuesSpec.Slurmdbd.Image.PullSecrets,
			},
			"diagnosticMode": map[string]interface{}{
				"enabled": valuesSpec.Slurmdbd.DiagnosticMode.Enabled,
				"command": valuesSpec.Slurmdbd.DiagnosticMode.Command,
				"args":    valuesSpec.Slurmdbd.DiagnosticMode.Args,
			},
			"automountServiceAccountToken": false,
			"podLabels":                    map[string]string{},
			"podAnnatations":               map[string]string{},
			"podAffinityPreset":            "",
			"podAntiAffinityPreset":        "soft",
			"nodeAffinityPreset": map[string]interface{}{
				"type":   valuesSpec.Slurmdbd.NodeAffinityPreset.Type,
				"key":    valuesSpec.Slurmdbd.NodeAffinityPreset.Key,
				"values": valuesSpec.Slurmdbd.NodeAffinityPreset.Values,
				"weight": valuesSpec.Slurmdbd.NodeAffinityPreset.Weight,
			},
			"hostNetwork":               false,
			"dnsPolicy":                 "",
			"dnsConfig":                 map[string]string{},
			"hostIPC":                   false,
			"priorityClassName":         "",
			"nodeSelector":              valuesSpec.Slurmdbd.NodeSelector,
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
			"extraVolumes":      valuesSpec.Slurmdbd.ExtraVolumes,
			"extraVolumeMounts": valuesSpec.Slurmdbd.ExtraVolumeMounts,
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
				"registry":    valuesSpec.SlurmLogin.Image.Registry,
				"repository":  valuesSpec.SlurmLogin.Image.Repository,
				"tag":         valuesSpec.SlurmLogin.Image.Tag,
				"pullPolicy":  valuesSpec.SlurmLogin.Image.PullPolicy,
				"pullSecrets": valuesSpec.SlurmLogin.Image.PullSecrets,
			},
			"diagnosticMode": map[string]interface{}{
				"enabled": valuesSpec.SlurmLogin.DiagnosticMode.Enabled,
				"command": valuesSpec.SlurmLogin.DiagnosticMode.Command,
				"args":    valuesSpec.SlurmLogin.DiagnosticMode.Args,
			},
			"automountServiceAccountToken": false,
			"podLabels":                    map[string]string{},
			"podAnnatations":               map[string]string{},
			"podAffinityPreset":            "",
			"podAntiAffinityPreset":        "soft",
			"nodeAffinityPreset": map[string]interface{}{
				"type":   valuesSpec.SlurmLogin.NodeAffinityPreset.Type,
				"key":    valuesSpec.SlurmLogin.NodeAffinityPreset.Key,
				"values": valuesSpec.SlurmLogin.NodeAffinityPreset.Values,
				"weight": valuesSpec.SlurmLogin.NodeAffinityPreset.Weight,
			},
			"hostNetwork":               false,
			"dnsPolicy":                 "",
			"dnsConfig":                 map[string]string{},
			"hostIPC":                   false,
			"priorityClassName":         "",
			"nodeSelector":              valuesSpec.SlurmLogin.NodeSelector,
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
					"cpu":               valuesSpec.SlurmLogin.Resources.Requests.CPU,
					"memory":            valuesSpec.SlurmLogin.Resources.Requests.Memory,
					"ephemeral-storage": valuesSpec.SlurmLogin.Resources.Requests.EphemeralStorage,
				},
				"limits": map[string]string{
					"cpu":               valuesSpec.SlurmLogin.Resources.Limits.CPU,
					"memory":            valuesSpec.SlurmLogin.Resources.Limits.Memory,
					"ephemeral-storage": valuesSpec.SlurmLogin.Resources.Limits.EphemeralStorage,
				},
			},
			"extraVolumes":      valuesSpec.SlurmLogin.ExtraVolumes,
			"extraVolumeMounts": valuesSpec.SlurmLogin.ExtraVolumeMounts,
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
CgroupPlugin=autodetect
IgnoreSystemd=yes
IgnoreSystemdOnFailure=yes
ConstrainDevices=yes
ConstrainRAMSpace=yes
ConstrainSwapSpace=no`,
			},
			"slurmConf": `ClusterName=slurm-cluster
SlurmctldHost={{ include "slurm.fullname" . }}-{{ .Values.slurmctld.name }}-0
MpiDefault=pmi2
DebugFlags=cgroup
SlurmdDebug=debug
ProctrackType=proctrack/cgroup
ReturnToService=1
SlurmctldPidFile=/var/run/slurmctld.pid
SlurmctldPort={{ .Values.slurmctld.service.slurmctld.port }}
SlurmdPidFile=/var/run/slurmd.pid
SlurmdPort=6818
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
AccountingStorageHost={{ include "slurm.fullname" . }}-{{ .Values.slurmdbd.name }}-0
AccountingStoragePort={{ .Values.slurmdbd.service.slurmdbd.port }}
AccountingStorageType=accounting_storage/slurmdbd
AccountingStoreFlags=job_comment
JobAcctGatherType=jobacct_gather/linux
JobAcctGatherFrequency=30
SlurmctldDebug=info
SlurmctldLogFile=/var/log/slurm/slurmctld.log
SlurmdLogFile=/var/log/slurm/slurmd.log
NodeName={{ include "slurm.fullname" . }}-slurmd-cpu-[0-` + fmt.Sprintf("%d", valuesSpec.SlurmdCPU.ReplicaCount+10) + `] CPUs=` + fmt.Sprintf("%d", valuesSpec.SlurmdCPU.Resources.Requests.Socket*valuesSpec.SlurmdCPU.Resources.Requests.CorePerSocket*valuesSpec.SlurmdCPU.Resources.Requests.ThreadPerCore) + ` Sockets=` + fmt.Sprintf("%d", valuesSpec.SlurmdCPU.Resources.Requests.Socket) + ` CoresPerSocket=` + fmt.Sprintf("%d", valuesSpec.SlurmdCPU.Resources.Requests.CorePerSocket) + ` ThreadsPerCore=` + fmt.Sprintf("%d", valuesSpec.SlurmdCPU.Resources.Requests.ThreadPerCore) + ` RealMemory=` + fmt.Sprintf("%d", ParseRAMstr(valuesSpec.SlurmdCPU.Resources.Requests.Memory)) + ` State=UNKNOWN
NodeName={{ include "slurm.fullname" . }}-slurmd-gpu-[0-` + fmt.Sprintf("%d", valuesSpec.SlurmdGPU.ReplicaCount+10) + `] CPUs=` + fmt.Sprintf("%d", valuesSpec.SlurmdGPU.Resources.Requests.Socket*valuesSpec.SlurmdGPU.Resources.Requests.CorePerSocket*valuesSpec.SlurmdGPU.Resources.Requests.ThreadPerCore) + ` Sockets=` + fmt.Sprintf("%d", valuesSpec.SlurmdGPU.Resources.Requests.Socket) + ` CoresPerSocket=` + fmt.Sprintf("%d", valuesSpec.SlurmdGPU.Resources.Requests.CorePerSocket) + ` ThreadsPerCore=` + fmt.Sprintf("%d", valuesSpec.SlurmdGPU.Resources.Requests.ThreadPerCore) + ` RealMemory=` + fmt.Sprintf("%d", ParseRAMstr(valuesSpec.SlurmdGPU.Resources.Requests.Memory)) + ` State=UNKNOWN
PartitionName=compute Nodes=ALL Default=YES MaxTime=INFINITE State=UP`,
			"slurmdbdConf": `AuthType=auth/munge
AuthInfo=/var/run/munge/munge.socket.2
SlurmUser=slurm
DebugLevel=verbose
LogFile=/var/log/slurm/slurmdbd.log
PidFile=/var/run/slurmdbd.pid
DbdHost={{ include "slurm.fullname" . }}-{{ .Values.slurmdbd.name }}-0
DbdPort={{ .Values.slurmdbd.service.slurmdbd.port }}
StorageType=accounting_storage/mysql
StorageHost={{ .Release.Name }}-mariadb
StoragePort={{ .Values.mariadb.port }}
StoragePass={{ .Values.mariadb.auth.password }}
StorageUser={{ .Values.mariadb.auth.username }}
StorageLoc={{ .Values.mariadb.auth.database }}`,
		},
	}
	return values
}
