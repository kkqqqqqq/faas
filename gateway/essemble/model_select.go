package essemble

// '''
// Terms:

//  latency:	max latency amont all instances,
//  accuracy:	aggreagted aacuracy of all instances,
// 	total cost of all instances

// slo_latency, slo_accuracy, slo_cost : constraints provided by user. (FIXME: Need some defaut values for each)

// latency_margin, accuracy_margin, cost_margin: margins provided by user. Defaut value = 0 if not provided

// Scaling Priority: User specified. Cost > Accuracy > latency (Default)

// README
// 	Function Lists:
// 		instace.init					:	initializes base object of class instance
// 		instance.add_model_to_instance	:	keeps track of all models in an instace and adds models as scale-in is called #can be called instance_model_manager
// 											returns 1 in case of successful scale in, returns 0 if fail
// 		scale_up						: 	adds one more instace to the instance list

// 	TODO:
// 		scale_down_policy				: incase the user changes the SLO in the runtime
// '''

// input: latency,acc, of user

var model_list = []string{
	"nasnetlarge",
	"inceptionresnetv2",
	"xception",
	"inceptionv3",
	"densenet201",
	"resnet50v2",
	"densenet121",
	"resnet50",
	"nasnetmobile",
	"mobilenetv2",
	"vgg16",
	"mobilenet"}
var models map[string]model

type model struct {
	name      string
	latency   float32
	accuracy  float32
	coldstart float32
	memory    int
	cpu       float32
}

//var model_latency = [315,151.96, 119.2, 74, 152.21, 89.5, 102.35, 98.22, 78.18, 41.5, 259, 43.45]

//var model_accuracy =[74.6,73, 69.75, 67.9, 72.83, 66, 70, 65,71.1, 68.05, 71.30, 68.36 ]

//var model_coldstart =[]

func ModelSelection(slo_latency float32, slo_accuracy float32, mode string) (model_selected []string) {

	if mode == "cocktail" {
		// cocketail
	}

	if mode == "inference" {
		//使用满足acc的单个模型
	}

	if mode == "efaas" {
		//efaas
		//input: latency ,acc
		//factors: latency,acc,cpu,memory,cold/warm

		modelInit()

		essemble_latency := float32(0)
		essemble_accuracy := float32(0)

		for essemble_accuracy < slo_accuracy && essemble_latency > slo_latency {

		}

		model_selected = []string{"mobilenet", "mobilenetv2", "resnet50v2", "nasnetmobile"}

	}
	return

}

func getEssembleLatency(model_selected []string) (essembleLatency float32) {
	essembleLatency = 0
	for _, model := range model_selected {
		essembleLatency = max(essembleLatency, getSingleLatency(model))
	}
	return essembleLatency
}

func getSingleLatency(model string) (latency float32) {
	latency = getModelLatency(model)
	if coldStart(model) {
		latency = getColdStartLatency(model) + latency
	}
	return latency
}

func coldStart(model string) bool {
	return false
}

func getColdStartLatency(model string) (latency float32) {
	return models[model].coldstart
}
func getModelLatency(model string) (latency float32) {
	return models[model].latency
}

func max(a float32, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
func modelInit() {
	models["nast"] = model{"nasnet", 1.0, 1.0, 22.0, 100, 0.2}
}
