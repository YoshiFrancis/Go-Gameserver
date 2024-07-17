package messages

import (
	"bytes"
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	type args struct {
		req []byte
	}
	tests := []struct {
		name     string
		args     args
		wantFlag byte
		wantArgs []string
	}{
		{
			name: "Basic",
			args: args{
				req: []byte("-2\r\n8\r\nSHUTDOWN\r\n4\r\nCOOL\r\n\r\n"),
			},
			wantFlag: '-',
			wantArgs: []string{"SHUTDOWN", "COOL"},
		},
		{
			name: "Basic2",
			args: args{
				req: []byte("-4\r\n9\r\nBROADCAST\r\n2\r\n42\r\n12\r\nYoshiKing101\r\n18\r\nYoshi is the king!\r\n\r\n"),
			},
			wantFlag: '-',
			wantArgs: []string{"BROADCAST", "42", "YoshiKing101", "Yoshi is the king!"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlag, gotArgs := Decode(tt.args.req)
			if gotFlag != tt.wantFlag {
				t.Errorf("Decode() gotFlag = %v, want %v", gotFlag, tt.wantFlag)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("Decode() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func Test_readSize(t *testing.T) {
	type args struct {
		r *bytes.Reader
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 bool
	}{
		{
			name: "Valid",
			args: args{
				r: bytes.NewReader([]byte("12345\r\n")),
			},
			want:  12345,
			want1: true,
		},
		{
			name: "Invalid",
			args: args{
				r: bytes.NewReader([]byte("1s2345\r\n")),
			},
			want:  -1,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := readSize(tt.args.r)
			if got != tt.want {
				t.Errorf("readSize() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("readSize() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
