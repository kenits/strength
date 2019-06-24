package strength

import (
	"fmt"
	"strings"

	"github.com/recoilme/pudge"
)

// Flex гибкая связь.
type Flex struct {
	ID              int     // номер
	Name            string  // имя
	Length          float64 // длинна пластины вдоль судна см
	Width           float64 // ширина пластины поперёк судна см
	ThicknessStart  float64 // толщина пластины мм
	Corrosion       float64 // годовая коррозия мм/год
	ThicknessEnd    float64 // толщина с учётом коррозии мм
	Height          float64 // положение центра тяжести относительно ОП м
	Count           float64 // колличество связей
	AreaStart       float64 // площадь в начале срока службы см2
	AreaEnd         float64 // площадь в конце срока службы с учётом колличества связей см2
	Pressure        float64 // поперечная нагрузка на пластину кПа
	StaticMoment    float64 // статический момент см2*м
	MomentOfInertia float64 // момент инерции см2*м2
}

// calc считает площадь, статический момент и момент инерции с учётом коррозии на срок службы.
func (f *Flex) calc(age float64) {
	f.AreaStart = (f.ThicknessStart / 10) * f.Width
	f.ThicknessEnd = f.ThicknessStart - f.Corrosion*age
	f.AreaEnd = calcAreaEnd(f.AreaStart, f.Width*f.Corrosion/10, age) * f.Count
	f.StaticMoment = calcStaticMoment(f.AreaEnd, f.Height)
	f.MomentOfInertia = calcMomentOfInertia(f.AreaEnd, f.Height)
}

// CalcAllFlex просчитать все гибкие связи.
func CalcAllFlex(data map[int]Flex, age float64) {
	for key, flex := range data {
		flex.calc(age)
		data[key] = flex
	}
}

// write пишет связь в базу.
func (f *Flex) write(addrDB string) error {
	err := pudge.Set(strings.Join([]string{addrDB, "flex"}, sep), f.ID, f)
	return err
}

// WriteFlex записать связь в базу данных.
func WriteFlex(flex *Flex, addrDB string) error {
	err := flex.write(addrDB)
	return err
}

// WriteAllFlex записать все связи в базу данных.
func WriteAllFlex(data *map[int]Flex, addrDB string) error {
	db, err := pudge.Open(strings.Join([]string{addrDB, "flex"}, sep), nil)
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

// ReadFlex прочитать связь из базы с заданным id.
func ReadFlex(id int, addrDB string) (Flex, error) {
	var rez Flex
	err := pudge.Get(strings.Join([]string{addrDB, "flex"}, sep), id, &rez)
	return rez, err
}

// ReadAllFlex прочитать данные всех связей из базы в виде карты [id]Rigid.
func ReadAllFlex(addrDB string) (map[int]Flex, error) {
	var (
		val Flex
	)
	db, err := pudge.Open(strings.Join([]string{addrDB, "flex"}, sep), nil)
	defer db.Close()

	if err != nil {
		return nil, err
	}

	data := make(map[int]Flex)

	count, err := db.Count()
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, fmt.Errorf("no items in flex")
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

// calcSumFlexArea рассчитывает суммарную площадь всех переданных связей.
func calcSumFlexArea(data map[int]Flex) float64 {
	var sum float64
	for _, val := range data {
		sum += val.AreaEnd
	}
	return sum
}

// calcSumFlexStaticMoment рассчитывает суммарный статический момент всех переданных связей.
func calcSumFlexStaticMoment(data map[int]Flex) float64 {
	var sum float64
	for _, val := range data {
		sum += val.StaticMoment
	}
	return sum
}

// calcSumFlexMomentOfInertia рассчитывает суммарный момент инерции всех переданных связей.
func calcSumFlexMomentOfInertia(data map[int]Flex) float64 {
	var sum float64
	for _, val := range data {
		sum += val.MomentOfInertia
	}
	return sum
}
