package internal

import (
	"bytes"
	"encoding/xml"
	"reflect"
	"strings"
	"testing"
	"time"
)

// https://tools.ietf.org/html/rfc4918#section-9.6.2
const exampleDeleteMultistatusStr = `<?xml version="1.0" encoding="utf-8" ?>
<d:multistatus xmlns:d="DAV:">
  <d:response>
    <d:href>http://www.example.com/container/resource3</d:href>
    <d:status>HTTP/1.1 423 Locked</d:status>
    <d:error><d:lock-token-submitted/></d:error>
  </d:response>
</d:multistatus>`

func TestResponse_Err_error(t *testing.T) {
	r := strings.NewReader(exampleDeleteMultistatusStr)
	var ms MultiStatus
	if err := xml.NewDecoder(r).Decode(&ms); err != nil {
		t.Fatalf("Decode() = %v", err)
	}

	if len(ms.Responses) != 1 {
		t.Fatalf("expected 1 <response>, got %v", len(ms.Responses))
	}

	resp := ms.Responses[0]

	err := resp.Err()
	if err == nil {
		t.Errorf("Multistatus.Get() returned a nil error, expected non-nil")
	} else if httpErr, ok := err.(*HTTPError); !ok {
		t.Errorf("Multistatus.Get() = %T, expected an *HTTPError", err)
	} else if httpErr.Code != 423 {
		t.Errorf("HTTPError.Code = %v, expected 423", httpErr.Code)
	}
}

func TestTimeRoundTrip(t *testing.T) {
	now := Time(time.Now().UTC())
	want, err := now.MarshalText()
	if err != nil {
		t.Fatalf("could not marshal time: %+v", err)
	}

	var got Time
	err = got.UnmarshalText(want)
	if err != nil {
		t.Fatalf("could not unmarshal time: %+v", err)
	}

	raw, err := got.MarshalText()
	if err != nil {
		t.Fatalf("could not marshal back: %+v", err)
	}

	if got, want := raw, want; !bytes.Equal(got, want) {
		t.Fatalf("invalid round-trip:\ngot= %s\nwant=%s", got, want)
	}
}

type Property struct {
	XMLName xml.Name
	// SkipOnAllprop skips this property on PROPFIND DAV:allprop
	SkipOnAllprop bool `xml:"-"`
}

func (p Property) WithChildren(children ...interface{}) *RawXMLValue {
	raws := make([]RawXMLValue, 0, len(children))
	for _, c := range children {
		switch r := c.(type) {
		case RawXMLValue:
			raws = append(raws, r)
		case *RawXMLValue:
			raws = append(raws, *r)
		case xml.Name:
			raws = append(raws, *NewRawXMLElement(r, nil, nil))
		default:
			raws = append(raws, *EncodeRawXMLElementMust(c))
		}
	}
	return NewRawXMLElement(p.XMLName, nil, raws)

}

func xmlDAV(name string) xml.Name {
	return xml.Name{Space: "DAV:", Local: name}
}

var CurrentUserPrincipal2 = Property{
	XMLName:       xmlDAV("2urrent-user-principal"),
	SkipOnAllprop: true,
}
var ResourceType2 = Property{
	XMLName:       xmlDAV("2esourcetype"),
	SkipOnAllprop: false,
}

func TestExtraction(t *testing.T) {
	marshal(t, CurrentUserPrincipal2)
	marshal(t, CurrentUserPrincipal2.WithChildren(HrefWrapperstruct{Href: &Href{Path: "principalPath"}}))

	marshal(t, NewResourceType(CollectionName))
	marshal(t, ResourceType2)
	marshal(t, ResourceType2.WithChildren(CollectionName))
	marshal(t, ResourceType2.WithChildren())

	// cup := Property{
	// 	XMLName: xml.Name{"DAV:", "lol"},
	// }
	// raw, _ := EncodeRawXMLElement(cup)
	// buf, err := xml.Marshal(raw)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(string(buf))

	// h, _ := EncodeRawXMLElement((HrefWrapperstruct{Href: &Href{Path: "/bim"}}))
	// // child := NewRawXMLElement(xml.Name{"DAV:", "href"}, nil, []RawXMLValue{*h})
	// raw = NewRawXMLElement(cup.XMLName, nil, []RawXMLValue{*h})
	// buf, err = xml.Marshal(raw)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(string(buf))
	// buf, err := xml.Marshal(NewOKResponse("/wtf"))
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(string(buf))
	// buf, err = xml.Marshal(Location{Href: Href{Path: "/wtf"}})
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(string(buf))

	// // got := parseXMLName(CurrentUserPrincipal{})
	// // want := xml.Name{"lo", "lo"}
	// // if want != got {

	// // 	t.Fatalf("invalid round-trip:\ngot= %s\nwant=%s", got, want)
	// }
	// t.Fail()

}

func marshal(t *testing.T, v interface{}) {
	t.Helper()
	buf, err := xml.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(buf))
}

func parseXMLName(v interface{}) xml.Name {
	var name xml.Name
	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return name
	}
	field, ok := t.FieldByName("XMLName")
	if !ok {
		return name
	}
	tag := field.Tag.Get("xml")

	if ns, t, ok := strings.Cut(tag, " "); ok {
		name.Space, tag = ns, t
	}
	name.Local = tag
	return name

}
