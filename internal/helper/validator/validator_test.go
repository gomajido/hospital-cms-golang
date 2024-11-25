package validator

import (
	"reflect"
	"testing"
)

func TestParseValidationError(t *testing.T) {
	type args struct {
		errorString string
	}
	tests := []struct {
		name    string
		args    args
		want    *ValidationError
		wantErr bool
	}{
		{
			name: "Valid error string",
			args: args{
				errorString: "Key: 'InputLogin.Password' Error:Field validation for 'Password' failed on the 'required' tag",
			},
			want: &ValidationError{
				Key:          "InputLogin.Password",
				ErrorMessage: "Field validation for 'Password' failed on the 'required' tag",
			},
			wantErr: false,
		},
		{
			name: "Invalid error string - missing key",
			args: args{
				errorString: "Error:Field validation for 'Password' failed on the 'required' tag",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid error string - missing error message",
			args: args{
				errorString: "Key: 'InputLogin.Password'",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid error string - completely malformed",
			args: args{
				errorString: "This is not a valid error string",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseValidationError(tt.args.errorString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseValidationError() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseValidationError() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNPWPValid(t *testing.T) {
	type args struct {
		npwp string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid NPWP with dots and dashes",
			args: args{
				npwp: "12.345.678.9-012.345",
			},
			want: true,
		},
		{
			name: "Valid NPWP without dots and dashes",
			args: args{
				npwp: "123456789012345",
			},
			want: true,
		},
		{
			name: "Invalid NPWP with letters",
			args: args{
				npwp: "12.345.678.9-012.34A",
			},
			want: false,
		},
		{
			name: "Invalid NPWP with special characters",
			args: args{
				npwp: "12.345.678.9-012.34@",
			},
			want: false,
		},
		{
			name: "Empty NPWP",
			args: args{
				npwp: "",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNPWPValid(tt.args.npwp); got != tt.want {
				t.Errorf("IsNPWPValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMimeTypeValid(t *testing.T) {
	type args struct {
		mimeType     string
		allowedMimes []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid MIME type",
			args: args{
				mimeType:     "image/jpeg",
				allowedMimes: []string{"image/jpeg", "image/png"},
			},
			want: true,
		},
		{
			name: "Invalid MIME type",
			args: args{
				mimeType:     "application/json",
				allowedMimes: []string{"image/jpeg", "image/png"},
			},
			want: false,
		},
		{
			name: "Empty MIME type",
			args: args{
				mimeType:     "",
				allowedMimes: []string{"image/jpeg", "image/png"},
			},
			want: false,
		},
		{
			name: "Empty allowed MIME types",
			args: args{
				mimeType:     "image/jpeg",
				allowedMimes: []string{},
			},
			want: false,
		},
		{
			name: "MIME type with wildcard",
			args: args{
				mimeType:     "image/jpeg",
				allowedMimes: []string{"image/*"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMimeTypeValid(tt.args.mimeType, tt.args.allowedMimes); got != tt.want {
				t.Errorf("IsMimeTypeValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
