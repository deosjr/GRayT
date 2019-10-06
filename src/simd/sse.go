package simd

func Add(a, b [4]float32) [4]float32

func Sub(a, b [4]float32) [4]float32

func Min(a, b [4]float32) [4]float32

func Max(a, b [4]float32) [4]float32

func Dot(a, b [4]float32) float32

func Cross(a, b [4]float32) [4]float32

func Normalize(a [4]float32) [4]float32

func Normalize4(x4, y4, z4 [4]float32) ([4]float32, [4]float32, [4]float32)

func BoxIntersect(origins, directions, mins, maxs [4]float32) ([4]float32, [4]float32)

func Box4Intersect(o4x, o4y, o4z, d4x, d4y, d4z, min4x, min4y, min4z, max4x, max4y, max4z [4]float32) [4]float32

func TriangleIntersect(p0, p1, p2, ro, rd [4]float32) float32

func Triangle4Intersect(p0x, p0y, p0z, p1x, p1y, p1z, p2x, p2y, p2z, rox, roy, roz, rdx, rdy, rdz [4]float32) [4]float32
