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
	before  *Node
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
		route = invertRoute(route)
	} else {
		node := new(Node)
		node = &Node{&hotels[firstID], nil, nil}
		route = Route{node}
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
	var currentElement *Node
	for i := startID; i <= endID; i++ {
		if i == startID {
			genesisNode := new(Node)
			genesisNode = &Node{&hotels[startID], nil, nil}
			route = Route{genesisNode}
			currentElement = route.First
		}
		node := new(Node)
		node = &Node{&hotels[i], nil, currentElement}
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
	for i := startX; i < endX; i++ {
		for n := 0; n < 2; n++ {
			img.SetColorIndex(i, startY+n, 1)
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
	palette := []color.Color{color.White, color.Black}
	rect := image.Rect(0, 0, width, 40)
	img := image.NewPaletted(rect, palette)
	anim := gif.GIF{Delay: []int{0}, Image: []*image.Paletted{img}}
	currentElement := route.First
	startY := 10
	startX := 10
	for {
		if currentElement.next == nil {
			drawHotel(img, startX, startY)
			break
		} else {
			drawHotel(img, startX, startY)
			startX = startX + 19
			drawConnection(img, startX, startX+currentElement.element.distance, startY+10)
			startX = startX + currentElement.element.distance
			currentElement = currentElement.next
		}
	}
	gif.EncodeAll(f, &anim)
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

func invertRoute(route Route) Route {
	currentElement := route.First
	for {
		if currentElement.next == nil {
			currentElement.next = currentElement.before
			currentElement.before = nil
			route.First = currentElement
			break
		}
		nextTemp := currentElement.next
		currentElement.next = currentElement.before
		currentElement.before = nextTemp
		currentElement = nextTemp
	}
	return route
}
