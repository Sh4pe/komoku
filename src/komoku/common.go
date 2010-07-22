package komoku

const (
    BoardSize = 19 // we support only quadratic boards at the moment
)

func posToXY(pos int) (x, y int) {
    return pos%BoardSize, pos/BoardSize
}

func xyToPos(x, y int) int {
    return 19*y + x
}
