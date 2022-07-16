package object

type Object interface {
	GetID() (ID int)
	SetID(ID int)
}

// Books
type Book struct {
	ID     int
	Author string
	Title  string
}

func (book *Book) GetID() (ID int) {
	return book.ID
}

func (book *Book) SetID(ID int) {
	book.ID = ID
}

// Cars
type Car struct {
	ID    int
	Brand string
	Model string
}

func (car *Car) GetID() (ID int) {
	return car.ID
}

func (car *Car) SetID(ID int) {
	car.ID = ID
}
