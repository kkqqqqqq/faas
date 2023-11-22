package essemble

import (
	"fmt"
	"log"
	"sort"
)

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
var Models = make(map[string]model)

func modelInit() {

	Models["nasetlarge"] = model{"nasnetlarge", 1.0, 1.0, 22.0, 100, 0.2, 11}
	Models["inceptionresnetv2"] = model{"inceptionresnetv2", 1.0, 1.0, 22.0, 100, 0.2, 10}
	Models["xception"] = model{"xception", 1.0, 1.0, 22.0, 100, 0.2, 8}
	Models["inceptionv3"] = model{"inceptionv3", 1.0, 1.0, 22.0, 100, 0.2, 9}
	Models["densenet201"] = model{"densenet201", 1.0, 1.0, 22.0, 100, 0.2, 5}
	Models["resnet50v2"] = model{"resnet50v2", 1.0, 1.0, 22.0, 100, 0.2, 7}
	Models["densenet121"] = model{"densenet121", 1.0, 1.0, 22.0, 100, 0.2, 4}
	Models["nasnetmobile"] = model{"nasnetmobile", 1.0, 1.0, 22.0, 100, 0.2, 3}
	Models["mobilenetv2"] = model{"mobilenetv2", 1.0, 1.0, 22.0, 100, 0.2, 2}
	//Models["vgg16"] = model{"vgg16", 1.0, 1.0, 22.0, 100, 0.2,}
	Models["mobilenet"] = model{"mobilenet", 1.0, 1.0, 22.0, 100, 0.2, 1}
}

type model struct {
	name      string
	latency   float32
	accuracy  float32
	coldstart float32
	memory    int
	cpu       float32
	inputtype int
}
type ModelSelectedInfo struct {
	Name      string
	Inputtype int
}

//	type modelMetrics struct {
//		name string
//		metricValue float32
//	}
func ModelSelection(slo_latency float32, slo_accuracy float32, mode string) (modelSelected []ModelSelectedInfo) {
	if mode == "cocktail" {
		log.Println("mode:cocktail..")
		// cocketail
	}

	if mode == "inference" {
		//使用满足acc的单个模型
	}

	if mode == "efaas" {
		log.Println("mode:effs..")
		//efaas
		//input: latency ,acc
		//factors: latency,acc,cpu,memory,cold/warm
		// 贪心策略，cocktail用的是滑动窗口
		// 优先选择什么样的模型：
		// acc 高，latency低，cost低
		modelInit()
		//model_selected = []string{"mobilenet", "mobilenetv2", "resnet50v2", "nasnetmobile"}
		//var modelSelected []string
		// 初始化当前状态下每个模型的μal，并排序，返回按照μAL排序后的modelname
		// sortedmodels := getSortedModels(models)
		// for _, j := range sortedmodels {
		// 	if essemble_accuracy < slo_accuracy || essemble_latency > slo_latency {
		// 		modelSelected = append(modelSelected, j)
		// 		essemble_accuracy = getEssembleAccuracy(modelSelected)
		// 		essemble_latency = getEssembleLatency(modelSelected)
		// 	} else {
		// 		break
		// 	}
		// }
		//modelSelected = []string{"mobilenet", "mobilenetv2", "resnet50v2"}
		modelSelected = []ModelSelectedInfo{
			{"mobilenet", 1},
			{"mobilenetv2", 2},
			{"resnet50v2", 7},
			{"nasnetlarge", 11}}
		//modelSelected = []ModelSelectedInfo{
		//	{"nasnetlarge", 11}}
		log.Println(modelSelected)

	}
	return

}

func getSortedModels(models map[string]model) []string {
	// 把model 按照metirc排序，返model回name的数组
	//求model的metics
	var metrics map[string]float32
	for key, value := range models {
		metricValue := value.accuracy / getSingleLatency(key)
		metrics[key] = metricValue
	}
	//把models按照meitrc进行排序

	// 提取 map 的键值对到切片
	var sortMetrics []struct {
		Key   string
		Value float32
	}
	for key, value := range metrics {
		sortMetrics = append(sortMetrics, struct {
			Key   string
			Value float32
		}{key, value})
	}
	// 自定义排序函数，按值排序
	sort.Slice(sortMetrics, func(i, j int) bool {
		return sortMetrics[i].Value < sortMetrics[j].Value
	})

	var sortedModels []string
	for _, modelMetric := range sortMetrics {
		sortedModels = append(sortedModels, modelMetric.Key)
		// 打印排序后的键值对
		fmt.Printf("%s: %d\n", modelMetric.Key, modelMetric.Value)
	}
	return sortedModels
}

func getEssembleAccuracy(model_selected []string) (essembleLatency float32) {
	essembleLatency = 0.0
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
	return Models[model].coldstart
}
func getModelLatency(model string) (latency float32) {
	return Models[model].latency
}

func max(a float32, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
