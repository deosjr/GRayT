package simd

func Add(a, b [4]float32) [4]float32

func Sub(a, b [4]float32) [4]float32

func Min(a, b [4]float32) [4]float32

func Max(a, b [4]float32) [4]float32

func Dot(a, b [4]float32) [4]float32

func Cross(a, b [4]float32) [4]float32

func BoxIntersect(origins, directions, mins, maxs [4]float32) ([4]float32, [4]float32)

func Box4Intersect(o4x, o4y, o4z, d4x, d4y, d4z, min4x, min4y, min4z, max4x, max4y, max4z [4]float32) [4]float32
