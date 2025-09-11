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

package v1

import (
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type ChartSpec struct {
	Name       string `json:"name"`
	Repository string `json:"repository"`
	Version    string `json:"version"`
	Namespace  string `json:"namespace,omitempty"`
}

type MariaDBSpec struct {
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`
	// +kubebuilder:default=3306
	Port    int32              `json:"port"`
	Auth    *MariaDBAuthSpec   `json:"auth,omitempty"`
	Primary MariaDBPrimarySpec `json:"primary,omitempty"`
}

type MariaDBAuthSpec struct {
	// +kubebuilder:default="slurm"
	Username string `json:"username,omitempty"`
	// +kubebuilder:default="password-for-slurm"
	Password string `json:"password,omitempty"`
	// +kubebuilder:default="rootpassword-for-slurm"
	RootPassword string `json:"rootPassword,omitempty"`
	// +kubebuilder:default="slurm_acct_db"
	DatabaseName string `json:"database,omitempty"`
}

type MariaDBPrimarySpec struct {
	Persistence MariaDBPrimaryPersistenceSpec `json:"persistence"`
}

type MariaDBPrimaryPersistenceSpec struct {
	// +kubebuilder:default=false
	Enabled bool `json:"enabled"`
	// +kubebuilder:default=""
	StorageClass string `json:"storageClass"`
	// +kubebuilder:default="2Gi"
	Size string `json:"size"`
}

type AuthSpec struct {
	SSH AuthSSHSpec `json:"ssh,omitempty"`
}

type AuthSSHSpec struct {
	Secret    AuthSSHSecretSpec    `json:"secret,omitempty"`
	ConfigMap AuthSSHConfigmapSpec `json:"configmap"`
}

type AuthSSHSecretSpec struct {
	// +kubebuilder:default="slurm-ssh-keys"
	Name string                `json:"name"`
	Keys AuthSSHSecretKeysSpec `json:"keys"`
}

type AuthSSHSecretKeysSpec struct {
	// +kubebuilder:default="id_rsa.pub"
	Public string `json:"public"`
	// +kubebuilder:default="id_rsa"
	Private string `json:"private"`
	// +kubebuilder:default="authorized_keys"
	AuthorizedKeys string `json:"authorizedKeys"`
}

type AuthSSHConfigmapSpec struct {
	// +kubebuilder:default="slurm-ssh-auth-keys"
	Name          string   `json:"name"`
	PrefabPubKeys []string `json:"prefabPubKeys"`
}

type NodeAffinityPreset struct {
	Type   string   `json:"type,omitempty"`
	Key    string   `json:"key,omitempty"`
	Values []string `json:"values,omitempty"`
	Weight int32    `json:"weight,omitempty"`
}

type PersistenceSpec struct {
	Shared PersistenceSharedSpec `json:"shared"`
}

type PersistenceSharedSpec struct {
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`
	// +kubebuilder:default="slurm-shared-storage"
	Name string `json:"name"`
	// +kubebuilder:default=""
	ExistingClaim string   `json:"existingClaim"`
	AccessModes   []string `json:"accessModes"`
	// +kubebuilder:default=""
	StorageClass string `json:"storageClass"`
	// +kubebuilder:default="8Gi"
	Size string `json:"size"`
}

type ImageSpec struct {
	// +kubebuilder:default="localhost"
	Registry string `json:"registry"`
	// +kubebuilder:default="data-and-computing"
	Repository string `json:"repository"`
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format="string-or-int"
	// +kubebuilder:default="latest"
	Tag string `json:"tag"`
	// +kubebuilder:default="IfNotPresent"
	PullPolicy  string   `json:"pullPolicy,omitempty"`
	PullSecrets []string `json:"pullSecrets,omitempty"`
}

type DiagnosticModeSpec struct {
	// +kubebuilder:default=false
	Enabled bool     `json:"enabled"`
	Command []string `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
}

type ExtraVolumeMountsSpec struct {
	// +kubebuilder:default=""
	Name string `json:"name"`
	// +kubebuilder:default=""
	MountPath string `json:"mountPath"`
}

type MungedSpec struct {
	// +kubebuilder:default="munged"
	Name              string                  `json:"name"`
	CommonLabels      map[string]string       `json:"commonLabels,omitempty"`
	Image             ImageSpec               `json:"image"`
	DiagnosticMode    DiagnosticModeSpec      `json:"diagnosticMode,omitempty"`
	ExtraVolumes      []corev1.Volume         `json:"extraVolumes,omitempty"`
	ExtraVolumeMounts []ExtraVolumeMountsSpec `json:"extraVolumeMounts,omitempty"`
}

type SlurmctldSpec struct {
	// +kubebuilder:default="slurmctld"
	Name         string            `json:"name"`
	CommonLabels map[string]string `json:"commonLabels,omitempty"`
	Image        ImageSpec         `json:"image"`
	// +kubebuilder:default=1
	ReplicaCount       int32                   `json:"replicaCount"`
	Resources          *ResourceSpec           `json:"resources,omitempty"`
	NodeAffinityPreset NodeAffinityPreset      `json:"nodeAffinityPreset,omitempty"`
	NodeSelector       map[string]string       `json:"nodeSelector,omitempty"`
	DiagnosticMode     DiagnosticModeSpec      `json:"diagnosticMode,omitempty"`
	ExtraVolumes       []map[string]string     `json:"extraVolumes,omitempty"`
	ExtraVolumeMounts  []ExtraVolumeMountsSpec `json:"extraVolumeMounts,omitempty"`
}

type SlurmdCPUSpec struct {
	// +kubebuilder:default="slurmd"
	Name         string    `json:"name"`
	CommonLabels []string  `json:"commonLabels,omitempty"`
	Image        ImageSpec `json:"image"`
	// +kubebuilder:default=0
	ReplicaCount       int32                   `json:"replicaCount"`
	Resources          SlurmdResourceSpec      `json:"resources,omitempty"`
	NodeAffinityPreset NodeAffinityPreset      `json:"nodeAffinityPreset,omitempty"`
	NodeSelector       map[string]string       `json:"nodeSelector,omitempty"`
	DiagnosticMode     DiagnosticModeSpec      `json:"diagnosticMode,omitempty"`
	ExtraVolumes       []corev1.Volume         `json:"extraVolumes,omitempty"`
	ExtraVolumeMounts  []ExtraVolumeMountsSpec `json:"extraVolumeMounts,omitempty"`
}

type SlurmdGPUSpec struct {
	// +kubebuilder:default="slurmd"
	Name         string    `json:"name"`
	CommonLabels []string  `json:"commonLabels,omitempty"`
	Image        ImageSpec `json:"image"`
	// +kubebuilder:default=0
	ReplicaCount       int32                   `json:"replicaCount"`
	Resources          SlurmdResourceSpec      `json:"resources,omitempty"`
	NodeAffinityPreset NodeAffinityPreset      `json:"nodeAffinityPreset,omitempty"`
	NodeSelector       map[string]string       `json:"nodeSelector,omitempty"`
	DiagnosticMode     DiagnosticModeSpec      `json:"diagnosticMode,omitempty"`
	ExtraVolumes       []corev1.Volume         `json:"extraVolumes,omitempty"`
	ExtraVolumeMounts  []ExtraVolumeMountsSpec `json:"extraVolumeMounts,omitempty"`
}

type SlurmdResourceSpec struct {
	Requests *SlurmdResourceRequestSpec `json:"requests,omitempty"`
	Limits   *SlurmdResourceLimitSpec   `json:"limits,omitempty"`
}

type ResourceSpec struct {
	Requests *ResourceRequestSpec `json:"requests,omitempty"`
	Limits   *ResourceLimitSpec   `json:"limits,omitempty"`
}

type SlurmdResourceRequestSpec struct {
	// +kubebuilder:default=1
	Socket int32 `json:"socket,omitempty"`
	// +kubebuilder:default=1
	CorePerSocket int32 `json:"core-per-socket,omitempty"`
	// +kubebuilder:default=1
	ThreadPerCore int32 `json:"thread-per-core,omitempty"`
	// +kubebuilder:default="1Gi"
	Memory string `json:"memory"`
	// +kubebuilder:default="2Gi"
	EphemeralStorage string `json:"ephemeral-storage"`
	// +kubebuilder:default=0
	GPUShared int32 `json:"nvidia.com/gpu.shared,omitempty"`
}

type SlurmdResourceLimitSpec struct {
	// +kubebuilder:default=1
	Socket int32 `json:"socket,omitempty"`
	// +kubebuilder:default=1
	CorePerSocket int32 `json:"core-per-socket,omitempty"`
	// +kubebuilder:default=1
	ThreadPerCore int32 `json:"thread-per-core,omitempty"`
	// +kubebuilder:default="2Gi"
	Memory string `json:"memory"`
	// +kubebuilder:default="8Gi"
	EphemeralStorage string `json:"ephemeral-storage"`
	// +kubebuilder:default=0
	GPUShared int32 `json:"nvidia.com/gpu.shared,omitempty"`
}

type ResourceRequestSpec struct {
	// +kubebuilder:default="500m"
	CPU string `json:"cpu"`
	// +kubebuilder:default="1Gi"
	Memory string `json:"memory"`
	// +kubebuilder:default="2Gi"
	EphemeralStorage string `json:"ephemeral-storage"`
}

type ResourceLimitSpec struct {
	// +kubebuilder:default="3000m"
	CPU string `json:"cpu"`
	// +kubebuilder:default="2Gi"
	Memory string `json:"memory"`
	// +kubebuilder:default="8Gi"
	EphemeralStorage string `json:"ephemeral-storage"`
}

type SlurmdbdSpec struct {
	// +kubebuilder:default="slurmdbd"
	Name               string                  `json:"name"`
	CommonLabels       map[string]string       `json:"commonLabels,omitempty"`
	Image              ImageSpec               `json:"image"`
	DiagnosticMode     DiagnosticModeSpec      `json:"diagnosticMode,omitempty"`
	NodeAffinityPreset NodeAffinityPreset      `json:"nodeAffinityPreset,omitempty"`
	NodeSelector       map[string]string       `json:"nodeSelector,omitempty"`
	ExtraVolumes       []corev1.Volume         `json:"extraVolumes,omitempty"`
	ExtraVolumeMounts  []ExtraVolumeMountsSpec `json:"extraVolumeMounts,omitempty"`
}

type SlurmLogindSpec struct {
	// +kubebuilder:default="login"
	Name               string                  `json:"name"`
	CommonLabels       map[string]string       `json:"commonLabels,omitempty"`
	Image              ImageSpec               `json:"image"`
	Resources          ResourceSpec            `json:"resources,omitempty"`
	DiagnosticMode     DiagnosticModeSpec      `json:"diagnosticMode,omitempty"`
	NodeAffinityPreset NodeAffinityPreset      `json:"nodeAffinityPreset,omitempty"`
	NodeSelector       map[string]string       `json:"nodeSelector,omitempty"`
	ExtraVolumes       []corev1.Volume         `json:"extraVolumes,omitempty"`
	ExtraVolumeMounts  []ExtraVolumeMountsSpec `json:"extraVolumeMounts,omitempty"`
}

type ServiceAccountSpec struct {
	// +kubebuilder:default=true
	Automount   bool              `json:"automount"`
	Annotations map[string]string `json:"annotations,omitempty"`
	// +kubebuilder:default="slurm"
	Name        string                        `json:"name"`
	Role        ServiceAccountRoleSpec        `json:"role"`
	RoleBinding ServiceAccountRoleBindingSpec `json:"roleBinding"`
}

type ServiceAccountRoleSpec struct {
	// +kubebuilder:default="slurm"
	Name string `json:"name"`
}

type ServiceAccountRoleBindingSpec struct {
	// +kubebuilder:default="slurm"
	Name string `json:"name"`
}

type SlurmConfigSpec struct {
	Cgroup       CgroupSpec `json:"cgroup"`
	SlurmConf    string     `json:"slurmConf"`
	SlurmdbdConf string     `json:"slurmdbdConf"`
}

type CgroupSpec struct {
	// +kubebuilder:default="cgroup-conf"
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ImageMirrorSpec struct {
	Mirror MirrorSpec `json:"mirror"`
}

type MirrorSpec struct {
	Registry string `json:"registry"`
}

type ValuesSpec struct {
	Mariadb     MariaDBSpec     `json:"mariadb"`
	Auth        AuthSpec        `json:"auth,omitempty"`
	Persistence PersistenceSpec `json:"persistence,omitempty"`
	ImageMirror ImageMirrorSpec `json:"image,omitempty"`
	Munged      MungedSpec      `json:"munged"`
	Slurmctld   SlurmctldSpec   `json:"slurmctld"`
	SlurmdCPU   SlurmdCPUSpec   `json:"slurmdCPU,omitempty"`
	SlurmdGPU   SlurmdGPUSpec   `json:"slurmdGPU,omitempty"`
	Slurmdbd    SlurmdbdSpec    `json:"slurmdbd"`
	SlurmLogin  SlurmLogindSpec `json:"login"`
	// +kubebuilder:default="nano"
	ResourcesPreset string             `json:"resourcesPreset,omitempty"`
	ServiceAccount  ServiceAccountSpec `json:"serviceAccount,omitempty"`
	SlurmConfig     SlurmConfigSpec    `json:"configuration,omitempty"`
	// +kubebuilder:default=""
	NameOverride string `json:"nameOverride,omitempty"`
	// +kubebuilder:default=""
	FullnameOverride  string            `json:"fullnameOverride,omitempty"`
	CommonAnnotations map[string]string `json:"commonAnnotations,omitempty"`
	CommonLabels      map[string]string `json:"commonLabels,omitempty"`
}

// SlurmDeploymentSpec defines the desired state of SlurmDeployment.
type SlurmDeploymentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Chart  ChartSpec    `json:"chart"`
	Job    SlurmJobSpec `json:"job,omitempty"`
	Values ValuesSpec   `json:"values"`
}

type SlurmJobSpec struct {
	Command []string `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
}

// SlurmDeploymentStatus defines the observed state of SlurmDeployment.
type SlurmDeploymentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	CPUNodeCount        string `json:"cpuNodeCount,omitempty"`
	CPUNodeStsVersion   string `json:"cpuNodeStsVersion"`
	GPUNodeCount        string `json:"gpuNodeCount,omitempty"`
	GPUNodeStsVersion   string `json:"gpuNodeStsVersion"`
	LoginNodeCount      string `json:"loginNodeCount,omitempty"`
	ControldDeamonCount string `json:"ctldNodeCount,omitempty"`
	DatabaseDeamonCount string `json:"databaseDeamonCount,omitempty"`
	MariadbServiceCount string `json:"mariadbServiceCount,omitempty"`
	JobCommand          string `json:"jobCommand,omitempty"`
	ClusterStatus       string `json:"clusterStatus,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=sd;slurmdep
// +kubebuilder:printcolumn:name="CPU",type="string",JSONPath=".status.cpuNodeCount",description="Number of CPU nodes"
// +kubebuilder:printcolumn:name="GPU",type="string",JSONPath=".status.gpuNodeCount",description="Number of GPU nodes"
// +kubebuilder:printcolumn:name="Login",type="string",JSONPath=".status.loginNodeCount",description="Number of Login nodes"
// +kubebuilder:printcolumn:name="Ctld",type="string",JSONPath=".status.ctldNodeCount",description="Number of Ctld nodes"
// +kubebuilder:printcolumn:name="DBd",type="string",JSONPath=".status.databaseDeamonCount",description="Number of Db nodes"
// +kubebuilder:printcolumn:name="DBsvc",type="string",JSONPath=".status.mariadbServiceCount",description="Number of mariadb nodes"
// +kubebuilder:printcolumn:name="Job Command",type="string",JSONPath=".status.jobCommand",description="Current job command"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.clusterStatus",description="Cluster status"

// SlurmDeployment is the Schema for the slurmdeployments API.
type SlurmDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SlurmDeploymentSpec   `json:"spec,omitempty"`
	Status SlurmDeploymentStatus `json:"status,omitempty"`
}

func (v *ValuesSpec) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Mariadb           MariaDBSpec        `json:"mariadb"`
		Auth              AuthSpec           `json:"auth,omitempty"`
		Persistence       PersistenceSpec    `json:"persistence,omitempty"`
		ImageMirror       ImageMirrorSpec    `json:"image,omitempty"`
		Munged            MungedSpec         `json:"munged"`
		Slurmctld         SlurmctldSpec      `json:"slurmctld"`
		SlurmdCPU         SlurmdCPUSpec      `json:"slurmdCPU,omitempty"`
		SlurmdGPU         SlurmdGPUSpec      `json:"slurmdGPU,omitempty"`
		Slurmdbd          SlurmdbdSpec       `json:"slurmdbd"`
		SlurmLogin        SlurmLogindSpec    `json:"login"`
		ResourcesPreset   string             `json:"resourcesPreset,omitempty"`
		ServiceAccount    ServiceAccountSpec `json:"serviceAccount,omitempty"`
		SlurmConfig       SlurmConfigSpec    `json:"configuration,omitempty"`
		NameOverride      string             `json:"nameOverride,omitempty"`
		FullnameOverride  string             `json:"fullnameOverride,omitempty"`
		CommonAnnotations map[string]string  `json:"commonAnnotations,omitempty"`
		CommonLabels      map[string]string  `json:"commonLabels,omitempty"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	v.Mariadb = aux.Mariadb
	v.Auth = aux.Auth
	v.Persistence = aux.Persistence
	v.ImageMirror = aux.ImageMirror
	v.Munged = aux.Munged
	v.Slurmctld = aux.Slurmctld
	v.SlurmdCPU = aux.SlurmdCPU
	v.SlurmdGPU = aux.SlurmdGPU
	v.Slurmdbd = aux.Slurmdbd
	v.SlurmLogin = aux.SlurmLogin
	v.ResourcesPreset = aux.ResourcesPreset
	v.ServiceAccount = aux.ServiceAccount
	v.SlurmConfig = aux.SlurmConfig
	v.NameOverride = aux.NameOverride
	v.FullnameOverride = aux.FullnameOverride
	v.CommonAnnotations = aux.CommonAnnotations
	v.CommonLabels = aux.CommonLabels
	return nil
}

// +kubebuilder:object:root=true

// SlurmDeploymentList contains a list of SlurmDeployment.
type SlurmDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SlurmDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SlurmDeployment{}, &SlurmDeploymentList{})
}
