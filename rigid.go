package strength

// Rigid жёсткая связь.
type Rigid struct {
	ID              int     // номер
	Name            string  // имя
	AreaStart       float64 // площадь в начале срока службы см2
	Corrosion       float64 // годовая коррозия см2/год
	Height          float64 // положение центра тяжести относительно ОП м
	Count           float64 // колличество связей
	AreaEnd         float64 // площадь в конце срока службы с учётом колличества связей см2
	StaticMoment    float64 // статический момент см2*м
	MomentOfInertia float64 // момент инерции см2*м2
}

// calc считает площадь, статический момент и момент инерции с учётом коррозии на срок службы.
func (r *Rigid) calc(age float64) {
	r.AreaEnd = calcAreaEnd(r.AreaStart, r.Corrosion, age) * r.Count
	r.StaticMoment = calcStaticMoment(r.AreaEnd, r.Height)
	r.MomentOfInertia = calcMomentOfInertia(r.AreaEnd, r.Height)
}

// CalcAllRigid просчитать все жёсткие связи.
func CalcAllRigid(data map[int]Rigid, age float64) {
	for key, rigid := range data {
		rigid.calc(age)
		data[key] = rigid
	}
}

// calcSumRigidArea рассчитывает суммарную площадь всех переданных связей.
func calcSumRigidArea(data map[int]Rigid) float64 {
	var sum float64
	for _, val := range data {
		sum += val.AreaEnd
	}
	return sum
}

// calcSumRigidStaticMoment рассчитывает суммарный статический момент всех переданных связей.
func calcSumRigidStaticMoment(data map[int]Rigid) float64 {
	var sum float64
	for _, val := range data {
		sum += val.StaticMoment
	}
	return sum
}

// calcSumRigidMomentOfInertia рассчитывает суммарный момент инерции всех переданных связей.
func calcSumRigidMomentOfInertia(data map[int]Rigid) float64 {
	var sum float64
	for _, val := range data {
		sum += val.MomentOfInertia
	}
	return sum
}
