package main

import (
	"fmt"
	"math"
	"strconv"
	"flag"

	str "github.com/kenits/strength"

	excel "github.com/360EntSecGroup-Skylar/excelize/v2"
)

func main() {
	var (
		fileName string
	)
	flag.StringVar(&fileName, "file", "", "полное имя файла")
	flag.Parse()
	if fileName == "" {
		fmt.Println("Необходимо имя файла")
		return
	}
	err := calc(fileName)
	if err != nil {
		fmt.Println("Что-то пошло не так")
		fmt.Println(err)
		return
	}
	fmt.Println("Отработано")

}

func readBaseData(file *excel.File) (*str.BaseData, error) {
	nameSheet := "Исходные данные"

	addr := [8]string{
		"B1",
		"B2",
		"B4",
		"B5",
		"B6",
		"B7",
		"B8",
		"B9",
	}

	rez := [8]string{}

	for key, val := range addr {
		data, err := file.GetCellValue(nameSheet, val)
		if err != nil {
			return nil, err
		}
		rez[key] = data
	}

	project := rez[0]
	name := rez[1]
	momentFlag := false
	if rez[4] == "перегиб" {
		momentFlag = true
	}
	symmetry := false
	if rez[5] == "да" {
		symmetry = true
	}

	age, err := strconv.ParseFloat(rez[2], 64)
	if err != nil {
		return nil, err
	}
	elasticModul, err := strconv.ParseFloat(rez[7], 64)
	if err != nil {
		return nil, err
	}
	moment := 0.0
	if rez[3] != "" {
		moment, err = strconv.ParseFloat(rez[3], 64)
		if err != nil {
			return nil, err
		}
		moment = math.Abs(moment)
	}

	accuracy, err := strconv.ParseFloat(rez[6], 64)
	if err != nil {
		return nil, err
	}

	height, err := readVerticalArrayFloat(nameSheet, "A", 13, file)
	if err != nil {
		return nil, err
	}

	strain, err := readVerticalArrayFloat(nameSheet, "B", 13, file)
	if err != nil {
		return nil, err
	}

	if len(height) != len(strain) {
		return nil, fmt.Errorf("missing height/strain")
	}

	data := str.BaseData{
		Project:      project,
		Name:         name,
		Age:          age,
		Height:       height,
		Strain:       strain,
		ElasticModul: elasticModul,
		Symmetry:     symmetry,
		MomentFlag:   momentFlag,
		Moment:       moment,
		Accuracy:     accuracy,
	}
	return &data, nil

}

func readRigid(file *excel.File) (map[int]str.Rigid, error) {
	nameSheet := "Жёсткие связи"

	id, err := readVerticalArrayInt(nameSheet, "A", 2, file)
	if err != nil {
		return nil, err
	}
	name, err := readVerticalArray(nameSheet, "B", 2, file)
	if err != nil {
		return nil, err
	}
	area, err := readVerticalArrayFloat(nameSheet, "C", 2, file)
	if err != nil {
		return nil, err
	}
	corrosion, err := readVerticalArrayFloat(nameSheet, "D", 2, file)
	if err != nil {
		return nil, err
	}
	heigth, err := readVerticalArrayFloat(nameSheet, "E", 2, file)
	if err != nil {
		return nil, err
	}
	count, err := readVerticalArrayFloat(nameSheet, "F", 2, file)
	if err != nil {
		return nil, err
	}
	rigidMap := make(map[int]str.Rigid)
	for key := range id {
		rigid := str.Rigid{
			ID:        id[key],
			Name:      name[key],
			AreaStart: area[key],
			Corrosion: corrosion[key],
			Height:    heigth[key],
			Count:     count[key],
		}
		rigidMap[rigid.ID] = rigid

	}

	return rigidMap, nil

}
func readFlex(file *excel.File) (map[int]str.Flex, error) {
	nameSheet := "Гибкие связи"

	id, err := readVerticalArrayInt(nameSheet, "A", 2, file)
	if err != nil {
		return nil, err
	}
	name, err := readVerticalArray(nameSheet, "B", 2, file)
	if err != nil {
		return nil, err
	}
	length, err := readVerticalArrayFloat(nameSheet, "C", 2, file)
	if err != nil {
		return nil, err
	}
	width, err := readVerticalArrayFloat(nameSheet, "D", 2, file)
	if err != nil {
		return nil, err
	}
	thickness, err := readVerticalArrayFloat(nameSheet, "E", 2, file)
	if err != nil {
		return nil, err
	}
	corrosion, err := readVerticalArrayFloat(nameSheet, "F", 2, file)
	if err != nil {
		return nil, err
	}
	heigth, err := readVerticalArrayFloat(nameSheet, "G", 2, file)
	if err != nil {
		return nil, err
	}
	count, err := readVerticalArrayFloat(nameSheet, "H", 2, file)
	if err != nil {
		return nil, err
	}
	press, err := readVerticalArrayFloat(nameSheet, "I", 2, file)
	if err != nil {
		return nil, err
	}

	flexMap := make(map[int]str.Flex)
	for key := range id {
		flex := str.Flex{
			ID:             id[key],
			Name:           name[key],
			Length:         length[key],
			Width:          width[key],
			ThicknessStart: thickness[key],
			Corrosion:      corrosion[key],
			Height:         heigth[key],
			Count:          count[key],
			Pressure:       press[key],
		}
		flexMap[flex.ID] = flex

	}

	return flexMap, nil

}

func writeRigid(rigid map[int]str.Rigid, file *excel.File) error {
	nameSheet := "Жёсткие связи"
	for _, val := range rigid {
		err := writeRigidRow(nameSheet, &val, file)
		if err != nil {
			return err
		}
	}

	return nil

}
func writeFlex(flex map[int]str.Flex, file *excel.File) error {
	nameSheet := "Гибкие связи"

	for _, val := range flex {
		err := writeFlexRow(nameSheet, &val, file)
		if err != nil {
			return err
		}
	}

	return nil

}

func writeAllRezult(rezult map[int]str.Rezult, file *excel.File) error {
	for key, val := range rezult {
		err := writeRezult(key, &val, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeAllApprox(approx map[int]map[int]str.Approx, file *excel.File) error {
	for key, val := range approx {
		err := writeApprox(key, val, file)
		if err != nil {
			return err
		}
	}
	return nil

}

// calc сосчитать файл
func calc(fileName string) error {

	file, err := excel.OpenFile(fileName)
	if err != nil {
		return err
	}

	basedata, err := readBaseData(file)
	if err != nil {
		return err
	}

	rigid, err := readRigid(file)
	if err != nil {
		return err
	}

	str.CalcAllRigid(rigid, basedata.Age)
	err = writeRigid(rigid, file)
	if err != nil {
		return err
	}

	flex, err := readFlex(file)
	if err != nil {
		return err
	}

	str.CalcAllFlex(flex, basedata.Age)
	err = writeFlex(flex, file)
	if err != nil {
		return err
	}

	approx, rezult := str.Calculate(basedata, rigid, flex)
	if err != nil {
		return err
	}

	err = writeAllRezult(rezult, file)
	if err != nil {
		return err
	}

	err = writeAllApprox(approx, file)
	if err != nil {
		return err
	}
	file.SaveAs("rezult.xlsx")

	return nil
}
func readVerticalArrayFloat(sheetName, column string, row int, file *excel.File) ([]float64, error) {

	arr := make([]float64, 0)
	data, err := readVerticalArray(sheetName, column, row, file)
	if err != nil {
		return nil, err
	}
	for _, val := range data {
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return arr, err
		}
		arr = append(arr, floatVal)
	}
	return arr, nil

}

func readVerticalArrayInt(sheetName, column string, row int, file *excel.File) ([]int, error) {

	arr := make([]int, 0)
	data, err := readVerticalArray(sheetName, column, row, file)
	if err != nil {
		return nil, err
	}
	for _, val := range data {
		intVal, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return arr, err
		}
		arr = append(arr, int(intVal))
	}
	return arr, nil
}

func writeRigidRow(sheetName string, rigid *str.Rigid, file *excel.File) error {
	var (
		err error
	)
	row := rigid.ID + 1

	addrArea := fmt.Sprintf("G%d", row)
	addrStaticMoment := fmt.Sprintf("H%d", row)
	addrMomentOfInertia := fmt.Sprintf("I%d", row)
	err = file.SetCellValue(sheetName, addrArea, rigid.AreaEnd)
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, addrStaticMoment, rigid.StaticMoment)
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, addrMomentOfInertia, rigid.MomentOfInertia)
	if err != nil {
		return err
	}
	return nil

}

func writeFlexRow(sheetName string, flex *str.Flex, file *excel.File) error {
	var (
		err error
	)
	row := flex.ID + 1
	addrArea := fmt.Sprintf("J%d", row)
	addrStaticMoment := fmt.Sprintf("K%d", row)
	addrMomentOfInertia := fmt.Sprintf("L%d", row)
	err = file.SetCellValue(sheetName, addrArea, flex.AreaEnd)
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, addrStaticMoment, flex.StaticMoment)
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, addrMomentOfInertia, flex.MomentOfInertia)
	if err != nil {
		return err
	}
	return nil

}

func writeRezult(id int, rezult *str.Rezult, file *excel.File) error {
	sheetName := fmt.Sprintf("Результаты %d приближения", id)
	file.NewSheet(sheetName)
	names := make([]string, 4)
	names[0] = "Площадь"
	names[1] = "Статический момент"
	names[2] = "Центр масс"
	names[3] = "Момент инерции"
	vals := make([]float64, 4)
	vals[0] = rezult.Area
	vals[1] = rezult.StaticMoment
	vals[2] = rezult.CenterOfMass
	vals[3] = rezult.MomentOfInertia
	if rezult.Moment != 0 {
		names = append(names, "Предельный момент")
		vals = append(vals, rezult.Moment)
	}
	rowNumNames, err := writeVerticalArrayStrings(sheetName, "A", 1, names, file)
	if err != nil {
		return err
	}
	rowNumVals, err := writeVerticalArrayFloat(sheetName, "B", 1, vals, file)
	if err != nil {
		return err
	}

	if rowNumNames != rowNumVals {
		return fmt.Errorf("not simetry write")
	}

	tableStartRow := 7
	err = file.SetCellValue(sheetName, fmt.Sprintf("A%d", tableStartRow), "Высота")
	if err != nil {
		return err
	}
	_, err = writeVerticalArrayFloat(sheetName, "A", tableStartRow+1, rezult.Heigth, file)
	if err != nil {
		return err
	}

	err = file.SetCellValue(sheetName, fmt.Sprintf("B%d", tableStartRow), "Момент сопротивления")
	if err != nil {
		return err
	}
	_, err = writeVerticalArrayFloat(sheetName, "B", tableStartRow+1, rezult.MomentsOfResistance, file)
	if err != nil {
		return err
	}

	if rezult.Moment != 0 {
		err = file.SetCellValue(sheetName, fmt.Sprintf("C%d", tableStartRow), "Максимальные моменты")
		if err != nil {
			return err
		}
		_, err = writeVerticalArrayFloat(sheetName, "C", tableStartRow+1, rezult.Moments, file)
		if err != nil {
			return err
		}
	} else {
		err = file.SetCellValue(sheetName, fmt.Sprintf("C%d", tableStartRow), "Напряжения")
		if err != nil {
			return err
		}
		_, err = writeVerticalArrayFloat(sheetName, "C", tableStartRow+1, rezult.Strain, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeApproxRow(sheetName string, approx *str.Approx, file *excel.File) error {
	var (
		err error
	)
	row := approx.ID + 1
	addrID := fmt.Sprintf("A%d", row)
	addrReducing := fmt.Sprintf("B%d", row)
	addrReverseReducing := fmt.Sprintf("C%d", row)
	addrReducingArea := fmt.Sprintf("D%d", row)
	addrHeight := fmt.Sprintf("E%d", row)
	addrAreaLoss := fmt.Sprintf("F%d", row)
	addrStaticMomentLoss := fmt.Sprintf("G%d", row)
	addrMomentOfInertiaLoss := fmt.Sprintf("H%d", row)
	err = file.SetCellValue(sheetName, addrID, approx.ID)
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, addrReducing, approx.Reducing)
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, addrReverseReducing, approx.ReverseReducing)
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, addrReducingArea, approx.ReducingArea)
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, addrAreaLoss, approx.AreaLoss)
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, addrHeight, approx.Height)
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, addrStaticMomentLoss, approx.StaticMomentLoss)
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, addrMomentOfInertiaLoss, approx.MomentOfInertiaLoss)
	if err != nil {
		return err
	}
	return nil

}

func writeApprox(id int, approx map[int]str.Approx, file *excel.File) error {
	sheetName := fmt.Sprintf("Приближение %d", id)
	file.NewSheet(sheetName)
	err := writeApproxHead(sheetName, file)
	if err != nil {
		return err
	}
	for _, val := range approx {
		err = writeApproxRow(sheetName, &val, file)
		if err != nil {
			return err
		}
	}
	return nil

}

func writeApproxHead(sheetName string, file *excel.File) error {

	err := file.SetCellValue(sheetName, "A1", "№")
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, "B1", "Редукционый коэффициент")
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, "C1", "обратный редукционный коэффициент")
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, "D1", "Редуцируемая площадь")
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, "F1", "Потеря площади")
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, "E1", "Высота")
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, "G1", "Потеря статического момента")
	if err != nil {
		return err
	}
	err = file.SetCellValue(sheetName, "H1", "Потеря момента инерции")
	if err != nil {
		return err
	}
	return nil

}

func readVerticalArray(sheetName, column string, row int, file *excel.File) ([]string, error) {

	arr := make([]string, 0)

	for {
		addr := fmt.Sprintf("%s%d", column, row)

		val, err := file.GetCellValue(sheetName, addr)
		if err != nil {
			return nil, err
		}
		if val == "" {
			return arr, nil
		}
		arr = append(arr, val)
		row++
	}
}

func writeVerticalArrayStrings(sheetName, column string, row int, data []string, file *excel.File) (int, error) {

	for key, val := range data {
		addr := fmt.Sprintf("%s%d", column, row)

		err := file.SetCellValue(sheetName, addr, val)
		if err != nil {
			return key + row, err
		}
		row++
	}
	return len(data) + row, nil
}

func writeVerticalArrayFloat(sheetName, column string, row int, data []float64, file *excel.File) (int, error) {

	for key, val := range data {
		addr := fmt.Sprintf("%s%d", column, row)

		err := file.SetCellValue(sheetName, addr, val)
		if err != nil {
			return key + row, err
		}
		row++
	}
	return len(data) + row, nil
}
