package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	ErrPathNotFound = errors.New("bfs error: no paths found")
	ErrBufioScan = errors.New("scan error: attemp to scan data failed")
)

type Point struct {
	X, Y int
}

// directions - напрваления возможного движения
var directions = [4]Point{
	{0, 1}, // Вправо
	{1, 0}, // Вниз
	{0, -1}, // Влево
	{-1, 0}, // Вверх
}

// isValid проверяет, является ли клетка с координатами (x, y) допустимой для перемещения.
func isValid(x, y, rows, cols int) bool {
	return x >= 0 && x < rows && y >= 0 && y < cols
}

// BFSCheck - поиск в ширину; находит кратчайший путь от стартовой точки до финишной в лабиринте,
// двигаясь по клеткам, используя только горизонтальные и вертикальные перемещения.
func BFSCheck(maze [][]int, start, end Point) ([]Point, error) {
	rows := len(maze)
	cols := len(maze[0])
	
	visited := make([][]bool, rows) // Матрица посещенных точек
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	queue := []Point{start}

	visited[start.X][start.Y] = true // Считаем начальную точку посещенной

	parent := make(map[Point]Point)
	found := false
	for len(queue) > 0 {
		currentPoint := queue[0]
		queue = queue[1:]

		if currentPoint == end {
			found = true
			break
		}

		for _, direction := range directions {
			nextX, nextY := currentPoint.X + direction.X, currentPoint.Y + direction.Y

			if isValid(nextX, nextY, rows, cols) && maze[nextX][nextY] != 0 && !visited[nextX][nextY] {
				visited[nextX][nextY] = true
				parent[Point{nextX, nextY}] = currentPoint

				queue = append(queue, Point{nextX, nextY})
			}
		}
	}

	if !found {
		return nil, ErrPathNotFound
	}

	// Воссоздаем путь
	path := []Point{}
	
	var ok bool
	currentPoint := end
	for currentPoint != start {
		path = append([]Point{currentPoint}, path...)

		if currentPoint, ok = parent[currentPoint]; !ok {
			return nil, ErrPathNotFound
		}
	}

	path = append([]Point{start}, path...)

	return path, nil
}


func startScan(r io.Reader) error {
	scanner := bufio.NewScanner(r)

	// Считываем размер лабиринта
	if ok := scanner.Scan(); !ok {
		return ErrBufioScan
	} 
	dimensions := strings.Fields(scanner.Text())
	rows, _ := strconv.Atoi(dimensions[0]) // checkerr: пропускаем проверку ошибки - по условию
	cols, _ := strconv.Atoi(dimensions[1]) // checkerr: пропускаем проверку ошибки - по условию

	// Считываем сам лабиринт
	maze := make([][]int, rows)
	for i := 0; i < rows; i++ {
		if ok := scanner.Scan(); !ok {
			return ErrBufioScan
		} 
		
		line := strings.Fields(scanner.Text())
		maze[i] = make([]int, cols)
		for j := 0; j < cols; j++ {
			maze[i][j], _ = strconv.Atoi(line[j])
		}
	}

	// Считываем координаты старта и финиша
	if ok := scanner.Scan(); !ok {
		return ErrBufioScan
	} 
	startCoords := strings.Fields(scanner.Text())
	startX, _ := strconv.Atoi(startCoords[0]) // checkerr: пропускаем проверку ошибки - по условию
	startY, _ := strconv.Atoi(startCoords[1]) // checkerr: пропускаем проверку ошибки - по условию
	endX, _ := strconv.Atoi(startCoords[2]) // checkerr: пропускаем проверку ошибки - по условию
	endY, _ := strconv.Atoi(startCoords[3]) // checkerr: пропускаем проверку ошибки - по условию

	start, end := Point{startX, startY}, Point{endX, endY}

	path, err := BFSCheck(maze, start, end) // Находим путь с помощью BFSCheck
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	// Выводим путь
	for _, point := range path {
		fmt.Printf("%d %d\n", point.X, point.Y)
	}
	fmt.Println(".") // Точка в конце нужна, чтобы отделить результат от побочных данных (логи и т.п.)

	return nil
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case <-stop:
			fmt.Println("Программа завершена")
			os.Exit(0)
		case <-time.After(2 * time.Second): 
			fmt.Println("_______________________________") // Отделение текущего блока 
			fmt.Println("Ввод разрешен:") // Выводим служебное сообщение, чтобы сообщить пользователю - ввод разрешен

			if err := startScan(os.Stdin); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			// Можем вывести служебные логи, т.к. основная информация отделена точкой
		}
	}
}