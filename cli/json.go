package main

import (
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

// ParseBaseData разбираем базовые данные
func ParseBaseData(data *JSONData) str.BaseData {
	var (
		baseData str.BaseData
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
	Height = make([]float64,1)
	Strain = make([]float64,1)
	for _,v := range data.BaseData.ControlPoints {
		Height = append(Height, v.Height)
		Strain = append(Strain, v.Stress)
	}

	baseData.Height = Height
	baseData.Strain = Strain

	return baseData
}
