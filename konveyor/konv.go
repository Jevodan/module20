package konveyor

type Stadies func(channel chan int, done chan int) (filterChannel chan int) // Тип данных Стадия

type konveer struct {
	stage []Stadies
	done  chan int
}

func NewKonveer(done chan int, stadies ...Stadies) *konveer {
	k := konveer{stage: stadies, done: done}
	return &k
}

/*
Запуск конвеера. выходной  канал в цикле предыдущей стадии
является входным каналом слкдующей стадии. Возвращает выходной канал последней стадии
source - источник данных и входной канал первой стадии
с - обработанные данные после конвеера
*/
func (k *konveer) Start(source chan int) chan int {
	var c chan int = source
	for _, v := range k.stage {
		c = v(c, k.done) // Вызов внутренних функций filterOne filterSecond filterThree
	}
	return c
}
