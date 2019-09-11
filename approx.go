package strength

import (
	"fmt"
	"math"
	"math/cmplx"
)

// Approx приближение одной пластины.
type Approx struct {
	ID                       int     // номер пластины
	Reducing                 float64 // редукционный коэффициент
	ReverseReducing          float64 // обратный редукционный коэффициент
	ReducingArea             float64 // площадь подлежащая редуцированию см2
	AreaLoss                 float64 // потеря площади см2
	Height                   float64 // высота м
	StaticMomentLoss         float64 // потеря статического момента см2*м
	MomentOfInertiaLoss      float64 // потеря момента инерции см2*м2
	length, width, thickness float64 // размеры связи для расчёта остальных параметров
	pressure                 float64 // расчтёное давление
	count                    float64 // количество связей
}

func createAllApprox(flex *map[int]Flex, moment, centerOfMass, momentOfInertia, elasticModul float64, momentFlag bool) map[int]Approx {
	rez := make(map[int]Approx)

	for key, val := range *flex {
		data := fillApprox(&val)
		actStrain := calcActualStrain(data.Height, moment, centerOfMass, momentOfInertia, momentFlag)
		startCurv := data.calcStartCurvature()
		reducing := data.calcReducing(actStrain, startCurv, elasticModul)
		data.Reducing = reducing
		data.calc(reducing)
		rez[key] = data
	}
	return rez

}

// fillApprox заполняет приближение.
func fillApprox(plate *Flex) Approx {
	var a Approx
	a.ID = plate.ID
	a.length = plate.Length
	a.width = plate.Width
	a.thickness = plate.ThicknessEnd
	a.Height = plate.Height
	a.pressure = plate.Pressure
	a.count = plate.Count
	return a
}

func (a *Approx) calc(reducing float64) {
	a.ReducingArea = (a.width - math.Min(a.length, a.width)/2) * (a.thickness / 10)
	a.ReverseReducing = 1 - reducing
	a.AreaLoss = a.ReducingArea * a.ReverseReducing * a.count
	a.StaticMomentLoss = calcStaticMoment(a.AreaLoss, a.Height)
	a.MomentOfInertiaLoss = calcMomentOfInertia(a.AreaLoss, a.Height)

}
func (a *Approx) calcReducing(actStrain, startCurv, elasticModul float64) float64 {
	var (
		pressCurv float64
	)

	// проверка коэффициента на физический смысл
	limitCheck := func(phi float64) float64 {
		if phi >= 1 {
			return 1
		}
		if phi <= 0 {
			return 0
		}
		return phi
	}

	// продольные связи
	if a.length > a.width {
		// растяжение
		if actStrain > 0 {
			return 1
		}
		// сжатие
		eulStrain := a.calcEulerianStrain()
		return limitCheck(-eulStrain / actStrain)
	}

	// поперечные связи
	eulStrain := a.calcEulerianStrain()
	rho := a.calcRho()
	if a.pressure != 0 {
		kappa, _ := calcKappa(a.width / a.length)
		pressCurv = a.calcPressCurvature(kappa, elasticModul)
	}
	x := a.calcX(rho, startCurv, pressCurv, eulStrain, actStrain)
	chnStrain := calcChainStrain(x, rho, eulStrain)

	// ратсяжение
	if actStrain > 0 {
		return limitCheck(chnStrain / actStrain)
	}
	// сжатие
	return math.Min(limitCheck(-eulStrain/actStrain), limitCheck(chnStrain/actStrain))

}

func (a *Approx) calcEulerianStrain() float64 {
	var rez float64
	if a.length >= a.width {
		rez = 7.6 * math.Pow(10*a.thickness/a.width, 2)
	} else {
		rez = 1.9 * math.Pow(10*a.thickness/a.length, 2) * math.Pow(1+math.Pow(a.length, 2)/math.Pow(a.width, 2), 2)
	}
	return rez
}

func (a *Approx) calcStartCurvature() float64 { // в см
	var rez float64
	rez = a.length / 60 * (1.5/a.thickness + 0.4)
	return rez
}

func (a *Approx) calcPressCurvature(k, e float64) float64 { // в см
	var rez float64
	rez = k * a.pressure * math.Pow(a.length, 4) / (e * math.Pow(a.thickness/10, 3))
	return rez
}

func (a *Approx) calcRho() float64 {
	var rez float64
	if a.pressure == 0 {
		rez = 1
	} else {
		rez = 4 - 2.81*a.length/a.width + 1.34*math.Pow(a.length, 2)/math.Pow(a.width, 2)
	}
	return rez
}

func calcChainStrain(x, rho, eulerianStrain float64) float64 {
	rez := rho * (x - 1) * eulerianStrain
	return rez

}

func calcKappa(relations float64) (float64, error) {
	// TODO: Придумать что-то получше под эту таблицу
	var (
		tableOfK, tableOfRelations [14]float64
		startKey, endKey           int
	)
	interpolation := func(key, firstKey, firstVal, secondKey, secondVal float64) float64 {
		val := firstVal + (secondVal-firstVal)/(secondKey-firstKey)*(key-firstKey)
		return val
	}
	tableOfK[0] = 0.0138
	tableOfK[1] = 0.0165
	tableOfK[2] = 0.0191
	tableOfK[3] = 0.0210
	tableOfK[4] = 0.0227
	tableOfK[5] = 0.0241
	tableOfK[6] = 0.0251
	tableOfK[7] = 0.0260
	tableOfK[8] = 0.0267
	tableOfK[9] = 0.0272
	tableOfK[10] = 0.0276
	tableOfK[11] = 0.0279
	tableOfK[12] = 0.0282
	tableOfK[13] = 0.0284
	tableOfRelations[0] = 1
	tableOfRelations[1] = 1.1
	tableOfRelations[2] = 1.2
	tableOfRelations[3] = 1.3
	tableOfRelations[4] = 1.4
	tableOfRelations[5] = 1.5
	tableOfRelations[6] = 1.6
	tableOfRelations[7] = 1.7
	tableOfRelations[8] = 1.8
	tableOfRelations[9] = 1.9
	tableOfRelations[10] = 2
	tableOfRelations[11] = 3
	tableOfRelations[12] = 4
	tableOfRelations[13] = 5
	if relations < 1 {
		return 0, fmt.Errorf("bad ratio")
	}
	if relations >= 5 {
		return tableOfK[13], nil
	}
	// Поиск между какими значениями попадает relations
	for k, v := range tableOfRelations {
		if v > relations {
			endKey = k
			startKey = endKey - 1
			break
		}
	}

	rez := interpolation(relations,
		tableOfRelations[startKey],
		tableOfK[startKey],
		tableOfRelations[endKey],
		tableOfK[endKey],
	)
	return rez, nil

}

func calcCubicEquation(a, b, c, d float64) (float64, complex128, complex128, error) {
	var (
		y1                 float64
		y2, y3             complex128
		x1                 float64
		x2, x3             complex128
		realPart, imagPart float64
		err                error
	)
	if a == 0 {
		err = fmt.Errorf("Not cubic equation")
		return x1, x2, x3, err
	}

	p := (3*a*c - math.Pow(b, 2)) / (3 * math.Pow(a, 2))
	q := (2*math.Pow(b, 3) - 9*a*b*c + 27*math.Pow(a, 2)*d) / (27 * math.Pow(a, 3))
	Q := math.Pow(p/3, 3) + math.Pow(q/2, 2)
	if Q >= 0 {
		// Если корень не целый и под корнем отрицательное число то выдаёт NaN
		// Эта фигня чтоб это обойти, степень задавать с точками иначе жрёт как int (?) подумать на досуге
		var (
			alpha, beta float64
		)
		alphaCubicSqrt := -q/2 + math.Sqrt(Q)
		if alphaCubicSqrt < 0 {
			alpha = -math.Pow(math.Abs(alphaCubicSqrt), 1.0/3.0)
		} else {
			alpha = math.Pow(alphaCubicSqrt, 1.0/3.0)
		}
		betaCubicSqrt := -q/2 - math.Sqrt(Q)
		if betaCubicSqrt < 0 {
			beta = -math.Pow(math.Abs(betaCubicSqrt), 1.0/3.0)
		} else {
			beta = math.Pow(betaCubicSqrt, 1.0/3.0)
		}
		y1 = alpha + beta
		realPart = -(alpha + beta) / 2
		imagPart = math.Sqrt(3) * (alpha - beta) / 2
		y2 = complex(realPart, imagPart)
		y3 = complex(realPart, -imagPart)

	} else {
		alpha := cmplx.Pow(complex(-q/2, 0)+cmplx.Sqrt(complex(Q, 0)), 1.0/3.0)
		beta := cmplx.Pow(complex(-q/2, 0)-cmplx.Sqrt(complex(Q, 0)), 1.0/3.0)
		y1 = real(alpha + beta)
		realPart = real(-(alpha + beta) / 2)
		imagPart = imag(cmplx.Sqrt(3) * (alpha - beta) / 2)
		y2 = complex(realPart+imagPart, 0)
		y3 = complex(realPart-imagPart, 0)

	}
	x1 = y1 - b/(3*a)
	x2 = y2 - complex(b/(3*a), 0)
	x3 = y3 - complex(b/(3*a), 0)
	return x1, x2, x3, nil

}

func calcActualStrain(height, moment, centerOfMass, momentOfInertia float64, momentFlag bool) float64 {

	rez := moment / momentOfResistance(
		momentOfInertia,
		centerOfMass,
		height,
	)

	// расстановка знаков  действующих напряжений в соответствии с растяжением(+) сжатием (-)
	if (height < centerOfMass) && momentFlag { //перегиб и ниже цт
		rez = -rez

	}
	if (height > centerOfMass) && !momentFlag { //прогиб и выше цт
		rez = -rez
	}

	return rez

}

func momentOfResistance(momentOfInertia, centerOfMass, height float64) float64 {
	rez := momentOfInertia / math.Abs(height-centerOfMass)
	return rez
}

// calcX решает кубическое уравнение для нахождения цепных напряжений
func (a *Approx) calcX(rho, startCurv, pressCurv, eulStrain, actStrain float64) float64 {
	var (
		roots [3]float64
		root  float64
	)
	devisor := rho * math.Pow(1+math.Pow(a.length, 2)/math.Pow(a.width, 2), 2)
	squareFactor := 2.73/devisor*math.Pow(startCurv*10/a.thickness, 2) - actStrain/(rho*eulStrain) - 1
	freeFactor := 2.73 * math.Pow(pressCurv+startCurv, 2) / (devisor * math.Pow(a.thickness/10, 2))
	x1, x2, x3, _ := calcCubicEquation(1, squareFactor, 0, -freeFactor)

	// если корень меньше 0 или больше 1 то цепные напряжения становятся больше эйлеровых что не имеет смысла
	// так как они не достигаются и пластина теряет устойчивость по эйлеровым напряжениям

	// FIXME: зедесь точно намудил, надо поправить, нужен первый положительный корень, возможно вынести в отдельную функцию

	// отсечение мнимых и отрицательных корней
	if x1 > 0 {
		roots[0] = x1
	}
	if imag(x2) == 0 && real(x2) > 0 {
		roots[1] = real(x2)
	}
	if imag(x3) == 0 && real(x2) > 0 {
		roots[2] = real(x3)
	}

	// выбор наименьшего положительного корня
	for _, val := range roots {
		if val == 0 {
			continue
		}
		if root == 0 {
			root = val
			continue
		}
		root = math.Min(root, val)
	}

	return root
}

func calcSumApproxArea(data map[int]Approx) float64 {
	var sum float64
	for _, val := range data {
		sum += val.AreaLoss
	}
	return sum
}

func calcSumApproxStaticMoment(data map[int]Approx) float64 {
	var sum float64
	for _, val := range data {
		sum += val.StaticMomentLoss
	}
	return sum
}
func calcSumApproxMomentOfInertia(data map[int]Approx) float64 {
	var sum float64
	for _, val := range data {
		sum += val.MomentOfInertiaLoss
	}
	return sum
}
