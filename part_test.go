package gopc

import (
	"reflect"
	"testing"
)

var fakeURL = "/doc/a.xml"

func Test_newPart(t *testing.T) {
	type args struct {
		uri               string
		contentType       string
		compressionOption CompressionOption
	}
	tests := []struct {
		name    string
		args    args
		want    *Part
		wantErr bool
	}{
		{"base", args{fakeURL, "application/HTML", CompressionNone}, &Part{fakeURL, "application/html", CompressionNone, nil}, false},
		{"baseWithParameters", args{fakeURL, "TEXT/html; charset=ISO-8859-4", CompressionNone}, &Part{fakeURL, "text/html; charset=ISO-8859-4", CompressionNone, nil}, false},
		{"baseWithTwoParameters", args{fakeURL, "TEXT/html; charset=ISO-8859-4;q=2", CompressionNone}, &Part{fakeURL, "text/html; charset=ISO-8859-4; q=2", CompressionNone, nil}, false},
		{"incorrectContentTypeInvalidMediaParameter", args{fakeURL, "TEXT/html; charset=ISO-8859-4 q=2", CompressionNone}, nil, true},
		{"incorrectContentTypeInvalidMediaParameterNoParamentreName", args{fakeURL, "TEXT/html; =ISO-8859-4", CompressionNone}, nil, true},
		{"incorrectContentTypeDuplicateParameterName", args{fakeURL, "TEXT/html; charset=ISO-8859-4; charset=ISO-8859-4", CompressionNone}, nil, true},
		{"incorrectContentTypeNoSlash", args{fakeURL, "application", CompressionNone}, nil, true},
		{"incorrectContentTypeUnexpectedContent", args{fakeURL, "application/html/html", CompressionNone}, nil, true},
		{"incorrectContentTypeNoMediaType", args{fakeURL, "/html", CompressionNone}, nil, true},
		{"incorrectContentTypeExpectedToken", args{fakeURL, "application/", CompressionNone}, nil, true},
		{"incorrectURI", args{"", "fakeContentType", CompressionNone}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newPart(tt.args.uri, tt.args.contentType, tt.args.compressionOption)
			if (err != nil) != tt.wantErr {
				t.Errorf("newPart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newPart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPart_AddRelationship(t *testing.T) {
	type args struct {
		id      string
		reltype string
		uri     string
	}
	tests := []struct {
		name    string
		p       *Part
		args    args
		want    *Part
		wantErr bool
	}{
		{"newRelationship", &Part{fakeURL, "fakeContentType", CompressionNone, nil}, args{"fakeId", "fakeType", "fakeTarget"}, &Part{fakeURL, "fakeContentType", CompressionNone, []*Relationship{&Relationship{"fakeId", "fakeType", "fakeTarget", ModeInternal}}}, false},
		{"existingID", &Part{fakeURL, "fakeContentType", CompressionNone, []*Relationship{&Relationship{"fakeId", "fakeType", "fakeTarget", ModeInternal}}}, args{"fakeId", "fakeType", "fakeTarget"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.AddRelationship(tt.args.id, tt.args.reltype, tt.args.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("Part.AddRelationship() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Part.AddRelationship() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPart_HasRelationship(t *testing.T) {
	tests := []struct {
		name string
		p    *Part
		want bool
	}{
		{"partRelationshipTrue", &Part{fakeURL, "fakeContentType", CompressionNone, []*Relationship{&Relationship{"fakeId", "fakeType", "fakeTarget", ModeInternal}}}, true},
		{"partRelationshipFalse", &Part{fakeURL, "fakeContentType", CompressionNone, nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.HasRelationship(); got != tt.want {
				t.Errorf("Part.HasRelationship() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPart_Relationships(t *testing.T) {
	tests := []struct {
		name string
		p    *Part
		want []*Relationship
	}{
		{"base", &Part{fakeURL, "fakeContentType", CompressionNone, make([]*Relationship, 0)}, make([]*Relationship, 0)},
		{"partRelationship", &Part{fakeURL, "fakeContentType", CompressionNone, []*Relationship{&Relationship{"fakeId", "fakeType", "fakeTarget", ModeInternal}}}, []*Relationship{&Relationship{"fakeId", "fakeType", "fakeTarget", ModeInternal}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Relationships(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Part.Relationships() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPart_URI(t *testing.T) {
	tests := []struct {
		name string
		p    *Part
		want string
	}{
		{"base", new(Part), ""},
		{"partURI", &Part{fakeURL, "fakeContentType", CompressionNone, nil}, fakeURL},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.URI(); got != tt.want {
				t.Errorf("Part.URI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPart_ContentType(t *testing.T) {
	tests := []struct {
		name string
		p    *Part
		want string
	}{
		{"base", new(Part), ""},
		{"partContentType", &Part{fakeURL, "fakeContentType", CompressionNone, nil}, "fakeContentType"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.ContentType(); got != tt.want {
				t.Errorf("Part.ContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPart_CompressionOption(t *testing.T) {
	tests := []struct {
		name string
		p    *Part
		want CompressionOption
	}{
		{"base", new(Part), CompressionNormal},
		{"partCompressionOption", &Part{fakeURL, "fakeContentType", CompressionNone, nil}, CompressionNone},
		{"partCompressionOption", &Part{fakeURL, "fakeContentType", CompressionMaximum, nil}, CompressionMaximum},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.CompressionOption(); got != tt.want {
				t.Errorf("Part.CompressionOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatePartName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{""}, true},
		{"invalidURL", args{"/docs%/a.xml"}, true},
		{"emptySegment", args{"/doc//a.xml"}, true},
		{"abs uri", args{"http://docs//a.xml"}, true},
		{"not rel uri", args{"docs/a.xml"}, true},
		{"endSlash", args{"/docs/a.xml/"}, true},
		{"endDot", args{"/docs/a.xml."}, true},
		{"dot", args{"/docs/./a.xml"}, true},
		{"twoDots", args{"/docs/../a.xml"}, true},
		{"reserved", args{"/docs/%7E/a.xml"}, true},
		{"withQuery", args{"/docs/a.xml?a=2"}, true},
		{"notencodechar", args{"/€/a.xml"}, true},
		{"encodedBSlash", args{"/%5C/a.xml"}, true},
		{"encodedBSlash", args{"/%2F/a.xml"}, true},
		{"encodechar", args{"/%E2%82%AC/a.xml"}, false},
		{"base", args{"/docs/a.xml"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidatePartName(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("ValidatePartName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
