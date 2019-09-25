package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	str "github.com/kenits/strength"
)

// JSONData структура получения данных из GUI
type JSONData struct {
	BaseData struct {
		ProjectName    string  `json:"projectName"`
		ReportName     string  `json:"reportName"`
		MomentCase     bool    `json:"momentCase"`
		CalculateCase  bool    `json:"calculateCase"`
		MomentValue    float64 `json:"momentValue"`
		Symmetry       bool    `json:"symmetry"`
		LifeTime       float64 `json:"lifeTime"`
		Accurasity     float64 `json:"accurasity"`
		ElasticModules float64 `json:"elasticModules"`
		ControlPoints  []struct {
			Height float64 `json:"height"`
			Stress float64 `json:"stress"`
		} `json:"controlPoints"`
	} `json:"baseData"`
	BeamData []struct {
		ID        int     `json:"id"`
		Name      string  `json:"name"`
		Area      float64 `json:"area"`
		Corrosion float64 `json:"corrosion"`
		Height    float64 `json:"height"`
		Count     float64 `json:"count"`
	} `json:"beamData"`
	PlateData []struct {
		ID        int     `json:"id"`
		Name      string  `json:"name"`
		Lenght    float64 `json:"lenght"`
		Width     float64 `json:"width"`
		Tickness  float64 `json:"tickness"`
		Corrosion float64 `json:"corrosion"`
		Height    float64 `json:"height"`
		Count     float64 `json:"count"`
		Press     float64 `json:"press"`
	} `json:"plateData"`
}

// Parse разбираем структуру
func Parse(dataStruct *JSONData) (str.BaseData, map[int]str.Rigid, map[int]str.Flex) {
	base := parseBaseData(dataStruct)
	rigidMap := parseBeamData(dataStruct)
	flexMap := parsePlateData(dataStruct)
	return base, rigidMap, flexMap
}

// ParseFile распарсить файл по имени
func ParseFile(fileName string) (str.BaseData, map[int]str.Rigid, map[int]str.Flex, error) {
	var (
		base     str.BaseData
		rigidMap map[int]str.Rigid
		flexMap  map[int]str.Flex
		data     JSONData
	)
	file, err := os.Open(fileName)
	if err != nil {
		return base, rigidMap, flexMap, fmt.Errorf("Не смогли открыть файл данных")
	}
	defer file.Close()
	jsonByteData, err := ioutil.ReadAll(file)
	if err != nil {
		return base, rigidMap, flexMap, fmt.Errorf("Не смогли прочиать файл")
	}
	err = json.Unmarshal(jsonByteData, &data)
	if err != nil {
		return base, rigidMap, flexMap, err // fmt.Errorf("Не смогли разобрать файл")
	}
	base, rigidMap, flexMap = Parse(&data)
	return base, rigidMap, flexMap, nil
}

// ParseBaseData разбираем базовые данные
func parseBaseData(data *JSONData) str.BaseData {
	var (
		baseData       str.BaseData
		Height, Strain []float64
	)
	baseData.Project = data.BaseData.ProjectName
	baseData.Name = data.BaseData.ReportName
	baseData.Age = data.BaseData.LifeTime
	baseData.Symmetry = data.BaseData.Symmetry
	baseData.MomentFlag = data.BaseData.MomentCase
	baseData.Accuracy = data.BaseData.Accurasity
	baseData.ElasticModul = data.BaseData.ElasticModules
	// CalculateCase если true расчёт на напряжений, если false предельного момента
	if data.BaseData.CalculateCase {
		baseData.Moment = data.BaseData.MomentValue
	}
	Height = make([]float64, 0)
	Strain = make([]float64, 0)
	for _, v := range data.BaseData.ControlPoints {
		Height = append(Height, v.Height)
		Strain = append(Strain, v.Stress)
	}

	baseData.Height = Height
	baseData.Strain = Strain

	return baseData
}

// ParseBeamData разбираем входные данные жёстких связей
// TODO: надо потестить, а то есть ощущение что все будут ссылатся на один обьекит
// то есть все будут одинаковы и равны последнему записанному
func parseBeamData(data *JSONData) map[int]str.Rigid {
	var (
		rigids map[int]str.Rigid
		rigid  str.Rigid
	)
	rigids = make(map[int]str.Rigid)

	for _, v := range data.BeamData {
		rigid.ID = v.ID
		rigid.Name = v.Name
		rigid.AreaStart = v.Area
		rigid.Corrosion = v.Corrosion
		rigid.Height = v.Height
		rigid.Count = v.Count
		rigids[rigid.ID] = rigid
	}

	return rigids
}

// ParsePlateData аналог ParseBeamData толко для гибких связей
func parsePlateData(data *JSONData) map[int]str.Flex {
	var (
		flexs map[int]str.Flex
		flex  str.Flex
	)
	flexs = make(map[int]str.Flex)

	for _, v := range data.PlateData {
		flex.ID = v.ID
		flex.Name = v.Name
		flex.Length = v.Lenght
		flex.Width = v.Width
		flex.ThicknessStart = v.Tickness
		flex.Corrosion = v.Corrosion
		flex.Height = v.Height
		flex.Count = v.Count
		flex.Pressure = v.Press
		flexs[flex.ID] = flex
	}

	return flexs
}

func writeJSON(aprox map[int]map[int]str.Approx, rezult map[int]str.Rezult) {
	// TODO: сделать это дерьмо

}
