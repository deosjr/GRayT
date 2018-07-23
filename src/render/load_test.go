package render

import (
	"bufio"
	"reflect"
	"strings"
	"testing"

	"model"
)

func TestLoadObj(t *testing.T) {
	for i, tt := range []struct {
		obj  string
		want model.Object
	}{
		{
			obj:  `# empty file`,
			want: nil,
		},
		{
			obj: `# note reversed vertex order
			v 1.0 -0.02 2.1754370e-002
			v 2 3 4
			v 4 5 6.0
			f 1 2 3`,
			want: model.NewComplexObject([]model.Object{
				model.NewTriangle(
					model.Vector{4, 5, 6},
					model.Vector{2, 3, 4},
					model.Vector{1.0, -0.02, 2.1754370e-002},
					model.NewColor(255, 0, 0)),
			}),
		},
	} {
		reader := strings.NewReader(tt.obj)
		scanner := bufio.NewScanner(reader)
		got, err := loadObj(scanner, model.NewColor(255, 0, 0))
		if err != nil {
			if tt.want == nil {
				continue
			}
			t.Errorf("%d): error in load: %s", i, err.Error())
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d): got %v want %v", i, got, tt.want)
		}
	}
}
