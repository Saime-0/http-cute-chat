package piper

/*
import (
	"fmt"
	"time"
)

type PointLine []*Point

type TimeLine struct {
	processName string
	start       time.Time
	points      *PointLine
}

func Start(pname ...interface{}) *TimeLine {
	return &TimeLine{
		processName: fmt.Sprint(pname),
		start:       time.Now(),
		points:      &PointLine{},
	}
}

func (t *TimeLine) Measure() {
	fmt.Println(t.processName, "execution time", time.Since(t.start))
}

func (n *Node) NewTimelinePoint(processName ...interface{}) (pointID int) {
	point := &Point{
		start:       time.Now(),
		processName: fmt.Sprint(processName),
	}
	pointID = len(*n.timeline.points)
	*n.timeline.points = append(*n.timeline.points, point)
	return
}

func (n *Node) CompletePoint(pointID int) {
	fmt.Println(n.timeline.points.String())
	(*n.timeline.points)[pointID] = nil
}

type Point struct {
	start       time.Time
	processName string
}

func (p PointLine) String() string {
	var str string
	for _, point := range p {
		if point != nil {
			str += point.processName + ":"
		}
	}
	return str
} */
