package strength

import (
	"math"
)

// BaseData исходные данные по проекту.
// Количество напряжений равно количеству высот (рассматриваемые точки по высоте).
type BaseData struct {
	Project        string    // проект
	Name           string    // имя расчёта
	Age            float64   // срок службы судна лет
	Height, Strain []float64 // расчётные точки по высотам с допускаемыми напряжениями м и кН/см2
	ElasticModul   float64   // модуль упругости материала кПа
	Symmetry       bool      // признак симетрии
	MomentFlag     bool      // false прогиб, true перегиб
	Moment         float64   // расчётный момент всегда положительный (если не задан то считается предельный) кН*м
	Accuracy       float64   // точность расчёт в %
}

// calcAreaEnd считает площадь на срок службы.
func calcAreaEnd(area, corrosion, age float64) float64 {
	rez := area - age*corrosion
	return rez

}

// calcStaticMoment считает статический момент.
func calcStaticMoment(area, height float64) float64 {
	rez := height * area
	return rez
}

// calcMomentOfInertia считает момент инерции.
func calcMomentOfInertia(area, height float64) float64 {
	rez := math.Pow(height, 2) * area
	return rez
}

// accuracyCheck проверяет точности, если точность удовлетворительная -> true.
// moment как индикатор если 0 то проверка по предельному моменту иначе по напряжениям.
func accuracyCheck(old, new *Rezult, accuracy, moment float64) bool {
	var difference float64
	flag := true

	if moment != 0 {

		for i := range old.Strain {
			difference = math.Abs((new.Strain[i] - old.Strain[i]) / new.Strain[i] * 100)
			if difference > accuracy {
				flag = false
			}

		}
		return flag

	}

	difference = math.Abs((new.Moment - old.Moment) / new.Moment * 100)
	if difference > accuracy {
		flag = false
	}

	return flag
}

// Calculate считает всё и добавляет данные в Rigid и Flex и выдаёт карты результатов
func Calculate(baseData *BaseData, rigid map[int]Rigid, flex map[int]Flex) (map[int]map[int]Approx, map[int]Rezult) {

	// TODO: написать тесты

	approxData := make(map[int]map[int]Approx)
	rezultData := make(map[int]Rezult)

	area := calcSumRigidArea(rigid) + calcSumFlexArea(flex)
	staticMoment := calcSumRigidStaticMoment(rigid) + calcSumFlexStaticMoment(flex)
	momentOfInertia := calcSumRigidMomentOfInertia(rigid) + calcSumFlexMomentOfInertia(flex)

	// Считаем первое приближение
	rezultData[1] = createRezult(area, staticMoment, momentOfInertia, 0, 0, 0, baseData.Height, baseData.Strain, baseData.Symmetry, baseData.Moment)

	// Расчёт 2 и последующих приближений

	for id := 2; ; id++ {

		// С определением момента можно как-то лучше, но пока не пойму как
		var (
			moment float64
		)

		if baseData.Moment == 0 {
			moment = rezultData[id-1].Moment
		} else {
			moment = baseData.Moment
		}

		approxData[id] = createAllApprox(&flex, moment, rezultData[id-1].CenterOfMass, rezultData[id-1].MomentOfInertia, baseData.ElasticModul, baseData.MomentFlag)

		areaLoss := calcSumApproxArea(approxData[id])
		staticMomentLoss := calcSumApproxStaticMoment(approxData[id])
		momentOfInertiaLoss := calcSumApproxMomentOfInertia(approxData[id])

		rezultData[id] = createRezult(area, staticMoment, momentOfInertia,
			areaLoss, staticMomentLoss, momentOfInertiaLoss,
			baseData.Height, baseData.Strain, baseData.Symmetry, baseData.Moment)

		// Сравнение нового и старого результата для выхода цикла
		old := rezultData[id-1]
		new := rezultData[id]
		if accuracyCheck(&old, &new, baseData.Accuracy, baseData.Moment) {
			break
		}

	}

	return approxData, rezultData

}
