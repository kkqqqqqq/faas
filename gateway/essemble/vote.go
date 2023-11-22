package essemble

import (
	"strconv"
	"strings"
)

func Vote(results [][]byte) (ans string) {

	//加权多数投票策略：为了提高多数投票的准确度，作者设计了一种加权多数投票策略。
	//在这种策略中，权重矩阵由考虑每个模型在每个类别上的准确度而生成。
	//这个权重矩阵的维度是L×N，其中L是唯一标签的数量，N是集成中使用的模型数量。
	//多数投票计算：多数投票的计算是基于每个模型的权重来实现的。对于每个唯一类别的预测，它将该类别的所有模型的权重相加，然后选择具有最大权重的类别作为多数投票的输出。
	//这意味着即使某些类别没有获得最高票数，如果与这些类别相关的模型具有更高的权重，那么这些类别仍然有可能成为最终输出。
	//考虑类别信息：这一策略与通常基于整体正确预测分配权重的投票策略不同。它考虑了每个类别的信息，以便更好地适应不同类别的图像。
	//预处理
	var ansList [][]float64
	//column, err = tf.NewTensor([1][100]int64{ {0} })
	for _, j := range results {
		ansList = append(ansList, stringToFloat64List(string(j)))
	}

	//// 创建一个存储结果的 tensor
	//resultTensor := tensorList[0] // 初始化为第一个 tensor

	//for _, tensor := range tensorList[1:] {
	//	resultTensor, _ = resultTensor + tensor
	//}
	//
	ans = avarage(ansList)

	return
}

func avarage(ansList [][]float64) (ans string) {
	m := len(ansList[0])
	ansslice := make([]float64, m)
	for _, j := range ansList {
		for pos, v := range j {
			ansslice[pos] += v
		}
	}
	return float64SliceToString(ansslice)
}

func weight(ansList [][]float64) (ans string) {
	ans = "weight"
	return ans

}

func stringToFloat64List(input string) []float64 {
	// 使用 strings.Split 将字符串拆分为切片
	strList := strings.Split(input, ",")

	// 创建 float64 切片
	floatList := make([]float64, len(strList))

	// 将字符串转换为 float64
	for i, str := range strList {
		floatVal, _ := strconv.ParseFloat(strings.TrimSpace(str), 64)

		floatList[i] = floatVal
	}

	return floatList
}

func float64SliceToString(floatSlice []float64) string {
	// Convert each float64 to a string
	stringSlice := make([]string, len(floatSlice))
	for i, v := range floatSlice {
		stringSlice[i] = strconv.FormatFloat(v, 'f', -1, 64)
	}

	// Join the string representations with a separator (e.g., comma)
	resultString := strings.Join(stringSlice, ",")

	return resultString
}
