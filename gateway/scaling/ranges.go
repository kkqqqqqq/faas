package scaling

const (
	// DefaultMinReplicas is the minimal amount of replicas for a service.
	DefaultMinReplicas = 0

	// DefaultMaxReplicas is the amount of replicas a service will auto-scale up to.
	DefaultMaxReplicas = 20

	// DefaultScalingFactor is the defining proportion for the scaling increments.
	DefaultScalingFactor = 20

	// DefaultTargetLoad
	DefaultTargetLoad = 10

	DefaultTypeScale = "rps"

	// MinScaleLabel label indicating min scale for a function
	MinScaleLabel = "com.openfaas.scale.min"

	// MaxScaleLabel label indicating max scale for a function
	MaxScaleLabel = "com.openfaas.scale.max"

	// ScalingFactorLabel label indicates the scaling factor for a function
	ScalingFactorLabel = "com.openfaas.scale.factor"

	// TargetLoadLabel see also DefaultTargetScale
	TargetLoadLabel = "com.openfaas.scale.target"

	// ScaleTypeLabel see also DefaultScaleType
	ScaleTypeLabel = "com.openfaas.scale.type"
)
