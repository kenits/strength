package main

// APIData структура получения данных из GUI
type APIData struct {
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
