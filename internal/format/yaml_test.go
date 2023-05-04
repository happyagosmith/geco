package format_test

import (
	"testing"

	"github.com/happyagosmith/geco/internal"
	"github.com/happyagosmith/geco/internal/format"
)

func TestYaml(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		want := "key1: v1\n" +
			"key2:\n" +
			"  key2-1: v2-1\n" +
			"  key2-2: v2-2\n" +
			"key3:\n" +
			"  - a\n" +
			"  - b\n" +
			"key4:\n" +
			"  - t1: 1\n" +
			"    t2: 2"
		y, _ := format.NewYaml([]byte(want))
		got, _ := y.String()
		internal.AssertEqualString(t, "output", got, want)
	})
}

func TestMergeYaml(t *testing.T) {
	type testCase struct {
		name string
		y1   string
		y2   string
		want string
	}

	tests := []testCase{
		{
			name: "same key",
			y1:   "key1: v1",
			y2:   "key1: v1-overwritten",
			want: "key1: v1-overwritten",
		},
		{
			name: "same key in map",
			y1: "key4:\n" +
				"  key4-1: v4-1\n" +
				"  key4-2: v4-2",
			y2: "key4:\n" +
				"  key4-1: v4-overwritten\n" +
				"  key4-2: v4-2",
			want: "key4:\n" +
				"  key4-1: v4-overwritten\n" +
				"  key4-2: v4-2",
		},
		{
			name: "different keys in map",
			y1: "key4:\n" +
				"  key4-2: v4-1",
			y2: "key4:\n" +
				"  key4-1: v4-2",
			want: "key4:\n" +
				"  key4-2: v4-1\n" +
				"  key4-1: v4-2",
		},
		{
			name: "list",
			y1: "key3:\n" +
				"  - a\n" +
				"  - b",
			y2: "key3:\n" +
				"  - c",
			want: "key3:\n" +
				"  - c",
		},
		{
			name: "list of map",
			y1: "key4:\n" +
				"  - t1: 1\n" +
				"    t2: 2",
			y2: "key4:\n" +
				"  - t1: 2\n" +
				"    t3: 2",
			want: "key4:\n" +
				"  - t1: 2\n" +
				"    t3: 2",
		},
		{
			name: "generate sorted yaml",
			y1: "key2: v2\n" +
				"key3: v3\n" +
				"key4: v4",
			y2: "key1: v1\n" +
				"key2: v2-overwritten\n" +
				"key4: v4-overwritten",
			want: "key2: v2-overwritten\n" +
				"key3: v3\n" +
				"key4: v4-overwritten\n" +
				"key1: v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ny1, _ := format.NewYaml([]byte(tt.y1))
			ny2, _ := format.NewYaml([]byte(tt.y2))
			ny1.Merge(ny2)

			got, _ := ny1.String()
			internal.AssertEqualString(t, "merged yaml", got, tt.want)
		})
	}
}
