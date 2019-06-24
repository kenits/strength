package strength

import (
	"fmt"
	"strings"

	"github.com/recoilme/pudge"
)

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

// write пишет связь в базу.
func (r *Rigid) write(addrDB string) error {
	err := pudge.Set(strings.Join([]string{addrDB, "rigid"}, sep), r.ID, r)
	return err
}

// WriteRigid записать связь в базу данных.
func WriteRigid(rigid *Rigid, addrDB string) error {
	err := rigid.write(addrDB)
	return err
}

// WriteAllRigid записать все связи в базу данных.
func WriteAllRigid(data *map[int]Rigid, addrDB string) error {
	db, err := pudge.Open(strings.Join([]string{addrDB, "rigid"}, sep), nil)
	defer db.Close()
	if err != nil {
		return err
	}
	for key, val := range *data {
		err = db.Set(key, val)
	}
	if err != nil {
		return err
	}
	return nil
}

// ReadRigid прочитать связь из базы с заданным id.
func ReadRigid(id int, addrDB string) (Rigid, error) {
	var rez Rigid
	err := pudge.Get(strings.Join([]string{addrDB, "rigid"}, sep), id, &rez)
	return rez, err
}

// ReadAllRigid прочитать данные всех связей из базы в виде карты [id]Rigid.
func ReadAllRigid(addrDB string) (map[int]Rigid, error) {
	var (
		val Rigid
	)
	db, err := pudge.Open(strings.Join([]string{addrDB, "rigid"}, sep), nil)
	defer db.Close()

	if err != nil {
		return nil, err
	}

	data := make(map[int]Rigid)

	count, err := db.Count()
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, fmt.Errorf("no items in rigid")
	}

	for i := 1; i <= count; i++ {
		err = db.Get(i, &val)
		if err != nil {
			return nil, err
		}
		data[i] = val
	}

	return data, nil
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
