package strength

import (
	"math"
	"strings"

	"github.com/recoilme/pudge"
)

const (
	sep string = "/"
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

// write пишет основные данные по расчёту в базу.
func (b *BaseData) write(addrDB string) error {
	err := pudge.Set(strings.Join([]string{addrDB, "baseData"}, sep), 1, b)
	return err
}

// WriteBaseData записывает в базу основные данные по расчёту.
func WriteBaseData(base *BaseData, addrDB string) error {
	err := base.write(addrDB)
	return err
}

// ReadBaseData читает из базы основные данные по расчёту.
func ReadBaseData(addrDB string) (BaseData, error) {
	var rez BaseData
	err := pudge.Get(strings.Join([]string{addrDB, "baseData"}, sep), 1, &rez)
	return rez, err
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

// CalculateToDB считает обую прочность записывая все данные в базу.
func CalculateToDB(baseData *BaseData, rigid map[int]Rigid, flex map[int]Flex, addressDB string) error {

	err := writeImputData(baseData, rigid, flex, addressDB)
	if err != nil {
		return err
	}

	area := calcSumRigidArea(rigid) + calcSumFlexArea(flex)
	staticMoment := calcSumRigidStaticMoment(rigid) + calcSumFlexStaticMoment(flex)
	momentOfInertia := calcSumRigidMomentOfInertia(rigid) + calcSumFlexMomentOfInertia(flex)

	rezult := createRezult(area, staticMoment, momentOfInertia, 0, 0, 0, baseData.Height, baseData.Strain, baseData.Symmetry, baseData.Moment)
	err = rezult.write(addressDB, 1)
	if err != nil {
		return err
	}

	for id := 2; ; id++ {

		moment := 0.0
		if baseData.Moment == 0 {
			moment = rezult.Moment
		} else {
			moment = baseData.Moment
		}

		approx := createAllApprox(&flex, moment, rezult.CenterOfMass, rezult.MomentOfInertia, baseData.ElasticModul, baseData.MomentFlag)

		err = writeAllApprox(&approx, addressDB, id)
		if err != nil {
			return err
		}

		areaLoss := calcSumApproxArea(approx)
		staticMomentLoss := calcSumApproxStaticMoment(approx)
		momentOfInertiaLoss := calcSumApproxMomentOfInertia(approx)

		newRezult := createRezult(area, staticMoment, momentOfInertia,
			areaLoss, staticMomentLoss, momentOfInertiaLoss,
			baseData.Height, baseData.Strain, baseData.Symmetry, baseData.Moment)

		err = newRezult.write(addressDB, id)
		if err != nil {
			return err
		}

		// Сравнение нового и старого результата для выхода цикла
		if accuracyCheck(&rezult, &newRezult, baseData.Accuracy, baseData.Moment) {
			break
		} else {
			rezult = newRezult
		}
	}

	return nil

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

// writeImputData записывает все исходные данные в базу.
func writeImputData(baseData *BaseData, rigid map[int]Rigid, flex map[int]Flex, addressDB string) error {

	err := WriteBaseData(baseData, addressDB)
	if err != nil {
		return err
	}
	err = WriteAllRigid(&rigid, addressDB)
	if err != nil {
		return err
	}
	err = WriteAllFlex(&flex, addressDB)
	if err != nil {
		return err
	}
	return nil
}

// Calculate считает всё и добавляет данные в Rigid и Flex и выдаёт карты результатов
func Calculate(baseData *BaseData, rigid map[int]Rigid, flex map[int]Flex) (map[int]map[int]Approx, map[int]Rezult) {

	// TODO: проверить это всё и написать тесты

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
