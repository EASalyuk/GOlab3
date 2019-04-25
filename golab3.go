package main

import "fmt"                                    //библиотека консольного вывода
import "os"                                     //библиотека чтения файлов
import "io"                                     //библиотека записи файлов
import geojson "github.com/paulmach/go.geojson" //библиотека парсинга geojson-данных
import "github.com/fogleman/gg"                 //библиотека рендера изображения

//Структура данных для хранения параметров изображения
type settings struct {
	stroke         string  //цвет линии
	stroke_width   float64 //толщина линии
	stroke_opacity float64 //прозрачность линии
	fill           string  //цвет заполнения
	fill_opacity   float64 //прозрачность заполнения
}

//Структура данных для хранения координат
type position struct {
	x float64 //хранение X-координаты
	y float64 //хранение Y-координаты
}

//Переменная для хранения данных из файла geojson в байтовом представлении
var data64 = []byte(ReadFile("map.geojson"))

//Переменная для хранения параметров изображения
var Properties = settings{}

//Массив для хранения координат изображения
var Coordinates [255]position

//Переменная для хранения количества координат изображения
var Num = 0

//Функция чтения данных из файла в строчном представлении
func ReadFile(geofile string) string {
	file, err := os.Open(geofile) //открытие файла geojson
	if err != nil {               //если файл не существует
		fmt.Println(err) //вывод в консоль сообщения об ошибке
		os.Exit(1)       //завершение действия с кодом ошибки 1
	}
	defer file.Close() //закрытие файла

	data64 := make([]byte, 64) //присваивание переменной для хранения данных
	var data string            //переменная для хранения данных из файла geojson в строковом представлении

	for {
		fc, err := file.Read(data64) //считывание
		if err == io.EOF {
			break
		}
		data += string(data64[:fc]) //запись строки в переменную
	}

	return data //возвращение значения переменной хранения данных из файла geojson в строковом представлении
}

//Функция структурирования данных из файла geojson в байтовом представлении для последующего парсинга
func Collection() *geojson.FeatureCollection {
	fc, _ := geojson.UnmarshalFeatureCollection(data64)
	return fc
}

//Функция рендера изображения
func Drowing() {
	const width = 1366                //значение ширины холста
	const height = 1024               //значение высоты холста
	var coefW float64 = 0.007 * width //коэффицент преобразования координат по ширине
	var coefH float64 = 0.01 * height //коэффицент преобразования координат по высоте
	//Подготовка пространства вывода
	dc := gg.NewContext(width, height) //создание новоего холста заданного размера
	dc.SetRGB(1, 1, 1)                 //выбор белого цвета
	dc.Clear()                         //завершение создания холста заданным цветом
	//Отрисовка фигуры по полученным координатам
	for coord := 0; coord < Num-1; coord++ {
		x := (Coordinates[coord].x + 55) * coefW  //применение к считанной координате X коэффициента изменения значения по ширине
		y := (-Coordinates[coord].y + 50) * coefH //применение к считанной координате Y коэффициента изменения значения по высоте
		if coord == 0 {                           //если первая координата
			dc.MoveTo(x, y) //установить ее как начальную координату фигуры
		} else { //иначе
			dc.LineTo(x, y) //установить ее как координату конца текущего отрезка
		}
	}
	dc.ClosePath()                  //завершить обозначение границ фигуры путем соединения координат конца последнего обозначенного отрезка с координатой начала фигуры
	dc.SetHexColor(Properties.fill) //установка цвета заливки фигуры
	dc.Fill()                       //заливка заданном цветом внутри обозначенной границы фигуры

	//Отрисовка контура фигуры по полученным координатам
	dc.SetHexColor(Properties.stroke)        //установка цвета контура фигуры
	dc.SetLineWidth(Properties.stroke_width) //установка толщины линий
	for coord := 0; coord < Num-1; coord++ {
		x := (Coordinates[coord].x + 55) * coefW  //применение к считанной координате X коэффициента изменения значения по ширине
		y := (-Coordinates[coord].y + 50) * coefH //применение к считанной координате Y коэффициента изменения значения по высоте
		if coord == 0 {                           //если первая координата
			dc.MoveTo(x, y) //установить ее как начальную координату фигуры
		} else { //иначе
			dc.LineTo(x, y) //установить ее как координату конца текущего отрезка
		}
	}
	dc.ClosePath() //завершить обозначение границ фигуры путем соединения координат конца последнего обозначенного отрезка с координатой начала фигуры
	dc.Stroke()    //установка контура вдоль обозначенной границы фигуры заданным цветом и толщиной

	//Сохранение в файл
	dc.SavePNG("out.png")
}

//Функция парсинга данных из файла geojson
func Parsing() {
	geojson_properties := Collection()

	//Парсинг параметров
	Properties.stroke = geojson_properties.Features[0].Properties["stroke"].(string)                  //парсинг параметра цвета линий
	Properties.stroke_width = geojson_properties.Features[0].Properties["stroke-width"].(float64)     //парсинг параметра ширины линий
	Properties.stroke_opacity = geojson_properties.Features[0].Properties["stroke-opacity"].(float64) //парсинг параметра празрачности линий
	Properties.fill = geojson_properties.Features[0].Properties["fill"].(string)                      //парсинг параметра цвета фигуры
	Properties.fill_opacity = geojson_properties.Features[0].Properties["fill-opacity"].(float64)     //парсинг параметра празрачности фигуры

	//Парсинг координат
	CoordinatesJSON := geojson_properties.Features[0].Geometry.Polygon //парсинг всех координат фигуры
	Num = len(CoordinatesJSON[0])                                      //установка количества координат в массиве
	for coord := 0; coord < Num; coord++ {
		Coordinates[coord].x = CoordinatesJSON[0][coord][0] //запись координаты X
		Coordinates[coord].y = CoordinatesJSON[0][coord][1] //запись координаты Y
	}
}

//Главная функция программы, реализующая функцию парсинга geojson-файла и функцию рендера
func main() {
	Parsing() //вызов функции парсинга файла geojson
	Drowing() //вызов функции рендера изображения
}
