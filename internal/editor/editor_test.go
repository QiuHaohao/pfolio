package editor

import (
	"reflect"
	"testing"
)

func TestMustOpen(t *testing.T) {
	type args struct {
		editor         string
		initialContent []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "no change",
			args: struct {
				editor         string
				initialContent []byte
			}{editor: "vim -c ':wq'", initialContent: []byte("123\n")},
			want: []byte("123\n"),
		},
		{
			name: "replace text",
			args: struct {
				editor         string
				initialContent []byte
			}{editor: "vim -c ':%s/123/321/' -c ':wq'", initialContent: []byte("123\n")},
			want: []byte("321\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MustOpen(tt.args.editor, tt.args.initialContent); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustOpen() = %v, want %v", got, tt.want)
			}
		})
	}
}
