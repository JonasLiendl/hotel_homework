package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"os"
	"strconv"
)

type Hotel struct {
	name     string
	distance int
}

type Node struct {
	element *Hotel
	next    *Node
}

type Route struct {
	First *Node
}

func main() {
	var firstHotelName string
	var lastHotelName string
	hotels := readDataFromCSV("Hotels.csv")
	fmt.Println("Name des Starts:")
	fmt.Scanln(&firstHotelName)
	fmt.Println("Name des Ziels:")
	fmt.Scanln(&lastHotelName)
	firstID, lastID, err := getIndexOfHotels(firstHotelName, lastHotelName, hotels)
	if err != nil {
		log.Fatalln(err)
	}
	var route Route
	var distance int
	if firstID < lastID {
		route, distance = calculateRoute(firstID, lastID, hotels)
	} else if firstID > lastID {
		route, distance = calculateRoute(lastID, firstID, hotels)
	} else {
		route = Route{nil}
		distance = 0
	}
	fmt.Printf("Der Abstand zwischen %s und %s beträgt %dkm\n", firstHotelName, lastHotelName, distance)
	drawGIF(route)
}

func readDataFromCSV(name string) []Hotel {
	file, err := os.Open(name)
	if err != nil {
		log.Fatalln(err)
	}
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}
	var hotels []Hotel
	var hotel Hotel
	for _, data := range records {
		hotel.name = data[0]
		hotel.distance, err = strconv.Atoi(data[1])
		hotels = append(hotels, hotel)
	}
	return hotels
}

func getIndexOfHotels(firstHotelName string, lastHotelName string, hotels []Hotel) (firstHotelID int, lastHotelID int, err error) {
	firstHotelID = -1
	lastHotelID = -1
	for index := range hotels {
		if firstHotelName == hotels[index].name {
			firstHotelID = index
		} else if lastHotelName == hotels[index].name {
			lastHotelID = index
		}
	}
	if (firstHotelID == -1) || (lastHotelID == -1) {
		return firstHotelID, lastHotelID, errors.New("Not all hotels found")
	}
	return firstHotelID, lastHotelID, nil
}

func calculateRoute(startID int, endID int, hotels []Hotel) (route Route, distance int) {
	distance = hotels[startID].distance
	route = Route{&Node{&hotels[startID], nil}}
	currentElement := route.First
	for i := startID + 1; i < endID; i++ {
		node := new(Node)
		node = &Node{&hotels[i], nil}
		currentElement.next = node
		currentElement = currentElement.next
		if i != endID {
			distance = distance + currentElement.element.distance
		}
	}
	return route, distance
}

func initGIF(out io.Writer, width int, height int) *image.Paletted {
	palette := []color.Color{color.White, color.Black}
	rect := image.Rect(0, 0, width, height)
	img := image.NewPaletted(rect, palette)
	anim := gif.GIF{Delay: []int{0}, Image: []*image.Paletted{img}}
	gif.EncodeAll(out, &anim)
	return img
}

func drawHotel(img *image.Paletted, startX int, startY int) {
	for i := 0; i < 20; i++ {
		for n := 0; n < 20; n++ {
			img.SetColorIndex(startX+i, startY+n, 1)
		}
	}
}

func drawConnection(img *image.Paletted, startX int, endX int, startY int) {
	for i := 0; i < endX; i++ {
		for n := 0; n < 2; n++ {
			img.SetColorIndex(startX+i, startY+n, 1)
		}
	}
}

func drawGIF(route Route) {
	f, err := os.Create("route.gif")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	width := calculateWidth(route)
	img := initGIF(f, width, 100)
	currentElement := route.First
	startY := 10
	startX := 10
	for {
		drawHotel(img, startX, startY)
		drawConnection(img, startX, startX+currentElement.element.distance, startY+10)
		startX = startX + 20 + currentElement.element.distance
		if currentElement.next == nil {
			drawHotel(img, startX, startY)
			drawConnection(img, startX, startX+currentElement.element.distance, startY+10)
			break
		}
	}
}

func calculateWidth(route Route) int {
	width := 20
	currentElement := route.First
	for {
		if currentElement.next == nil {
			width = width + 20
			break
		}
		width = width + 20 + currentElement.element.distance
		currentElement = currentElement.next
	}
	return width
}
