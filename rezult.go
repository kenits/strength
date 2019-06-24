package strength

import (
	"fmt"
	"math"
	"strings"

	"github.com/recoilme/pudge"
)

// Rezult результаты приближения.
type Rezult struct {
	Area                float64   // площадь сечения см2
	StaticMoment        float64   // статический момент сечения см2*м
	CenterOfMass        float64   // высота центра масс относительно ОП м
	MomentOfInertia     float64   // момент инерции сечения см2*м2
	MomentsOfResistance []float64 // моменты сопротивления в контрольных точеках см2*м
	Moments             []float64 // предельные моменты в контрольных точеках кН*м
	Moment              float64   // предельный момент кН*м
	Strain              []float64 // действующие напряжения в контрольных точеках кН/см2
	AreaLoss            float64   // потеря площади сечения корпуса см2
	StaticMomentLoss    float64   // потеря статического момента сечения корпуса см2*м
	MomentOfInertiaLoss float64   // потеря момента инерции сечения корпуса см2*м2
	Heigth              []float64 // высоты контрольных точек сечения м
}

// calcLimitMoments расчитывает предельные моменты.
func calcLimitMoments(strains, momentsOfResistance []float64) []float64 {
	rez := make([]float64, len(strains))
	for key := range strains {
		rez[key] = strains[key] * momentsOfResistance[key]
	}
	return rez

}

// calcStrain расчитывает напряжения.
func calcStrain(momentOfResistance []float64, moment float64) []float64 {
	data := make([]float64, len(momentOfResistance))
	for key, val := range momentOfResistance {
		data[key] = moment / val
	}
	return data
}

// calcMomentOfResistance инициализирует моменты сопротивления и рассчитывает их в на заданных высотах.
func calcMomentOfResistance(momentOfInertia, centerOfMass float64, height []float64) []float64 {

	rez := make([]float64, len(height))

	for key, val := range height {
		rez[key] = momentOfResistance(momentOfInertia, centerOfMass, val)
	}
	return rez
}

// createRezult создаёт результат по входным данным.
func createRezult(area, staticMoment, momentOfInertia float64,
	areaLoss, staticMomentLoss, momentOfInertiaLoss float64,
	height, strain []float64,
	simmetry bool, moment float64) Rezult {
	var (
		rez Rezult
	)
	rez.Heigth = height
	rez.Area = area - areaLoss
	rez.StaticMoment = staticMoment - staticMomentLoss
	deltaMoment := momentOfInertia - momentOfInertiaLoss
	rez.MomentOfInertia = deltaMoment - math.Pow(rez.StaticMoment, 2)/rez.Area

	rez.CenterOfMass = rez.StaticMoment / rez.Area

	rez.AreaLoss = areaLoss
	rez.StaticMomentLoss = staticMomentLoss
	rez.MomentOfInertiaLoss = momentOfInertiaLoss

	if simmetry {
		rez.Area *= 2
		rez.StaticMoment *= 2
		rez.MomentOfInertia *= 2
	}

	rez.MomentsOfResistance = calcMomentOfResistance(rez.MomentOfInertia, rez.CenterOfMass, height)

	if moment != 0 {
		rez.Strain = calcStrain(rez.MomentsOfResistance, moment)

	} else {
		rez.Moments = calcLimitMoments(strain, rez.MomentsOfResistance)
		rez.Moment = rez.Moments[0]
		for _, val := range rez.Moments {
			rez.Moment = math.Min(rez.Moment, val)
		}
	}

	return rez

}

// write пишет результат в базу.
func (r *Rezult) write(addrDB string, id int) error {
	err := pudge.Set(strings.Join([]string{addrDB, "rezult"}, sep), id, r)
	return err
}

// ReadAllRezult читает все результирующие таблицы в виде [step]Rezult.
func ReadAllRezult(addrDB string) (map[int]Rezult, error) {
	db, err := pudge.Open(strings.Join([]string{addrDB, "rezult"}, sep), nil)
	defer db.Close()

	if err != nil {
		return nil, err
	}

	data := make(map[int]Rezult)

	count, err := db.Count()
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, fmt.Errorf("no items in rezult")
	}

	for i := 1; i <= count; i++ {
		val := data[i]
		err = db.Get(i, &val)
		if err != nil {
			return nil, err
		}
		data[i] = val
	}

	return data, nil

}

// ReadRezultByStep читает из базы результат конкретного шага.
func ReadRezultByStep(step int, addrDB string) (Rezult, error) {
	var rez Rezult
	err := pudge.Get(strings.Join([]string{addrDB, "rezult"}, sep), step, &rez)
	return rez, err

}
