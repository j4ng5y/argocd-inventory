package kube

type (
	DeprecationsRemovalsAndAdditions struct {
		IsDeprecation bool
		IsRemoval     bool
		IsAddition    bool
		Resource      *GVK
		ReplacedBy    *GVK
		NotableChanges []string
	}
	GVK struct {
		Group   string
		Version string
		Kind    string
	}
)

var (
	// Groups
	extensions string = "extensions"
	networkingK8sIo string = "networking.k8s.io"
	apps string = "apps"
	policy string = "policy"

	// Versions
	v1 string = "v1"
	v1beta1 string = "v1beta1"
	v1beta2 string = "v1beta2"

	// Kinds
	Ingress string = "Ingress"
	IngressClass string = "IngressClass"
	NetworkPolicy string = "NetworkPolicy"
	DaemonSet string = "DaemonSet"
	Deployment string = "Deployment"
	StatefulSet string = "StatefulSet"
	ReplicaSet string = "ReplicaSet"
	PodSecurityPolicy string = "PodSecurityPolicy"


	deprecationsAndAdditions = map[string][]DeprecationsRemovalsAndAdditions{
		"1.16": {
			{
				IsRemoval: true,
				Resource: &GVK{
					Group: extensions,
					Version: v1beta1,
					Kind: NetworkPolicy,
				},
				ReplacedBy: &GVK{
					Group: networkingK8sIo,
					Version: v1,
					Kind: NetworkPolicy,
				},
			},
			{
				IsRemoval: true,
				Resource: &GVK{
					Group: extensions,
					Version: v1beta1,
					Kind: DaemonSet,
				},
				ReplacedBy: &GVK{
					Group: apps,
					Version: v1,
					Kind: DaemonSet,
				},
				NotableChanges: []string{
					"spec.templateGeneration is removed",
					"spec.selector is now required and immutable after creation; use the existing template labels as the selector for seamless upgrades",
					"spec.updateStrategy.type now defaults to RollingUpdate (the default in extensions/v1beta1 was OnDelete)",
				},
			},
			{
				IsRemoval: true,
				Resource: &GVK{
					Group: extensions,
					Version: v1beta2,
					Kind: DaemonSet,
				},
				ReplacedBy: &GVK{
					Group: apps,
					Version: v1,
					Kind: DaemonSet,
				},
				NotableChanges: []string{
					"spec.templateGeneration is removed",
					"spec.selector is now required and immutable after creation; use the existing template labels as the selector for seamless upgrades",
					"spec.updateStrategy.type now defaults to RollingUpdate (the default in extensions/v1beta1 was OnDelete)",
				},
			},
			{
				IsRemoval: true,
				Resource: &GVK{
					Group: extensions,
					Version: v1beta1,
					Kind: Deployment,
				},
				ReplacedBy: &GVK{
					Group: apps,
					Version: v1,
					Kind: Deployment,
				},
				NotableChanges: []string{
					"spec.rollbackTo is removed",
					"spec.selector is now required and immutable after creation; use the existing template labels as the selector for seamless upgrades",
					"spec.progressDeadlineSeconds now defaults to 600 seconds (the default in extensions/v1beta1 was no deadline)",
					"spec.revisionHistoryLimit now defaults to 10 (the default in apps/v1beta1 was 2, the default in extensions/v1beta1 was to retain all)",
					"maxSurge and maxUnavailable now default to 25% (the default in extensions/v1beta1 was 1)",
				},
			},
			{
				IsRemoval: true,
				Resource: &GVK{
					Group: apps,
					Version: v1beta1,
					Kind: Deployment,
				},
				ReplacedBy: &GVK{
					Group: apps,
					Version: v1,
					Kind: Deployment,
				},
				NotableChanges: []string{
					"spec.rollbackTo is removed",
					"spec.selector is now required and immutable after creation; use the existing template labels as the selector for seamless upgrades",
					"spec.progressDeadlineSeconds now defaults to 600 seconds (the default in extensions/v1beta1 was no deadline)",
					"spec.revisionHistoryLimit now defaults to 10 (the default in apps/v1beta1 was 2, the default in extensions/v1beta1 was to retain all)",
					"maxSurge and maxUnavailable now default to 25% (the default in extensions/v1beta1 was 1)",
				},
			},
			{
				IsRemoval: true,
				Resource: &GVK{
					Group: apps,
					Version: v1beta2,
					Kind: Deployment,
				},
				ReplacedBy: &GVK{
					Group: apps,
					Version: v1,
					Kind: Deployment,
				},
				NotableChanges: []string{
					"spec.rollbackTo is removed",
					"spec.selector is now required and immutable after creation; use the existing template labels as the selector for seamless upgrades",
					"spec.progressDeadlineSeconds now defaults to 600 seconds (the default in extensions/v1beta1 was no deadline)",
					"spec.revisionHistoryLimit now defaults to 10 (the default in apps/v1beta1 was 2, the default in extensions/v1beta1 was to retain all)",
					"maxSurge and maxUnavailable now default to 25% (the default in extensions/v1beta1 was 1)",
				},
			},
			{
				IsRemoval: true,
				Resource: &GVK{
					Group: apps,
					Version: v1beta1,
					Kind: StatefulSet,
				},
				ReplacedBy: &GVK{
					Group: apps,
					Version: v1,
					Kind: StatefulSet,
				},
				NotableChanges: []string{
					"spec.selector is now required and immutable after creation; use the existing template labels as the selector for seamless upgrades",
					"spec.updateStrategy.type now defaults to RollingUpdate (the default in apps/v1beta1 was OnDelete)",
				},
			},
			{
				IsRemoval: true,
				Resource: &GVK{
					Group: apps,
					Version: v1beta2,
					Kind: StatefulSet,
				},
				ReplacedBy: &GVK{
					Group: apps,
					Version: v1,
					Kind: StatefulSet,
				},
				NotableChanges: []string{
					"spec.selector is now required and immutable after creation; use the existing template labels as the selector for seamless upgrades",
					"spec.updateStrategy.type now defaults to RollingUpdate (the default in apps/v1beta1 was OnDelete)",
				},
			},
			{
				IsRemoval: true,
				Resource: &GVK{
					Group: extensions,
					Version: v1beta1,
					Kind: ReplicaSet,
				},
				ReplacedBy: &GVK{
					Group: apps,
					Version: v1,
					Kind: ReplicaSet,
				},
				NotableChanges: []string{
					"spec.selector is now required and immutable after creation; use the existing template labels as the selector for seamless upgrades",
				},
			},
			{
				IsRemoval: true,
				Resource: &GVK{
					Group: apps,
					Version: v1beta1,
					Kind: ReplicaSet,
				},
				ReplacedBy: &GVK{
					Group: apps,
					Version: v1,
					Kind: ReplicaSet,
				},
				NotableChanges: []string{
					"spec.selector is now required and immutable after creation; use the existing template labels as the selector for seamless upgrades",
				},
			},
			{
				IsRemoval: true,
				Resource: &GVK{
					Group: apps,
					Version: v1beta2,
					Kind: ReplicaSet,
				},
				ReplacedBy: &GVK{
					Group: apps,
					Version: v1,
					Kind: ReplicaSet,
				},
				NotableChanges: []string{
					"spec.selector is now required and immutable after creation; use the existing template labels as the selector for seamless upgrades",
				},
			},
			{
				IsRemoval: true,
				Resource: &GVK{
					Group: extensions,
					Version: v1beta1,
					Kind: PodSecurityPolicy,
				},
				ReplacedBy: &GVK{
					Group: policy,
					Version: v1beta1,
					Kind: PodSecurityPolicy,
				},
			},
		"1.17": {},
		"1.18": {},
		"1.19": {
			{
				IsDeprecation: true,
				Resource: &GVK{
					Group: extensions,
					Version: v1beta1,
					Kind: Ingress,
				},
				ReplacedBy: &GVK{
					Group: networkingK8sIo,
					Version: v1,
					Kind: Ingress,
				},
			},
			{
				IsAddition: true,
				Resource: &GVK{
					Group: networkingK8sIo,
					Version: v1,
					Kind: IngressClass,
				},
			},
			{
				IsDeprecation: true,
				Resource: &GVK{
					Group: "extensions",
					Version: "v1beta1",
					Kind: "Ingress",
				},
				ReplacedBy: &GVK{
					Group: "networking.k8s.io",
					Version: "v1",
					Kind: "Ingress",
				},
			},
			{
				IsDeprecation: true,
				Resource: &GVK{
					Group: "extensions",
					Version: "v1beta1",
					Kind: "Ingress",
				},
				ReplacedBy: &GVK{
					Group: "networking.k8s.io",
					Version: "v1",
					Kind: "Ingress",
				},
			},
		},

		"1.20": {},
		"1.21": {},
		"1.22": {},
		"1.23": {},
		"1.24": {},
		"1.25": {},
		"1.26": {},
		"1.27": {},
		"1.28": {},
		"1.29": {},
		"1.30": {},
	},
)
