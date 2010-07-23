package komoku

// ################ constants ####################
const (
    BoardSize = 19 // we support only quadratic boards at the moment
)

// ################ types ####################
type Point struct {
    x, y int
}

// ################ helper functions ####################

// TODO: use point!
func posToXY(pos int) (x, y int) {
    return pos%BoardSize, pos/BoardSize
}

// TODO: use point!
func xyToPos(x, y int) int {
    return 19*y + x
}

// Returns the neighbours of a field (x,y)
// TODO: implement this!
func neighbours(x, y int) []Point {
    return make([]Point, 1, 1)
}
