package main

//Implement Pic. It should return a slice of length dy, each element of which is a slice of dx 8-bit unsigned integers. When you run the program, it will display your picture, interpreting the integers as grayscale (well, bluescale) values.
//https://tour.golang.org/moretypes/18

import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
	var row [][]uint8
	for i := 0; i < dx; i++ {
		var column []uint8
		for j := 0; j < dy; j++ {
			column = append(column, uint8(i*j))
		}
		row = append(row, column)
	}
	return row
}

func main() {
	pic.Show(Pic)
}
