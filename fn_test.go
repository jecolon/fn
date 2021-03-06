package fn

import (
	"testing"
	"unicode/utf8"
)

func TestStripControl(t *testing.T) {
	s := "\a\b\f\n\r\v\t test 123.45 _ () {} :; | *? <> \" '"
	o := stripControl(s)
	w := " test 123.45 _ () {} :; | *? <> \" '"
	if o != w {
		t.Fatalf("wanted %s but got %s", w, o)
	}
}

func BenchmarkStripControl(b *testing.B) {
	s := "\a\b\f\n\r\v\t test 123.45 _ () {} :; | *? <> \" '"
	for i := 0; i < b.N; i++ {
		stripControl(s)
	}
}

func TestStripSpecial(t *testing.T) {
	s := "\a\b\f\n\r\v\t test 123.45 _ () {} :; | *? <> \" '"
	o := stripSpecial(s)
	w := "\a\b\f\n\r\v\t test 123.45 _        "
	if o != w {
		t.Fatalf("wanted %s but got %s", w, o)
	}
}

func BenchmarkStripSpecial(b *testing.B) {
	s := "\a\b\f\n\r\v\t test 123.45 _ () {} :; | *? <> \" '"
	for i := 0; i < b.N; i++ {
		stripSpecial(s)
	}
}

func TestReplaceSpaces(t *testing.T) {
	s := "\a\b\f\n\r\v\t test 123.45 _ () {} :; | *? <> \" '"
	r := '_'
	o := replaceSpaces(s, r)
	w := "\a\b\f\n\r\v\t_test_123.45___()_{}_:;_|_*?_<>_\"_'"
	if o != w {
		t.Fatalf("wanted %s but got %s", w, o)
	}
}

func BenchmarkReplaceSpaces(b *testing.B) {
	s := "\a\b\f\n\r\v\t test 123.45 _ () {} :; | *? <> \" '"
	r := '_'
	for i := 0; i < b.N; i++ {
		replaceSpaces(s, r)
	}
}

func TestTrim(t *testing.T) {
	s := "-this-is--a----test-.png-"
	c := "-"
	o := trim(s, c)
	w := "this-is-a-test.png"
	if o != w {
		t.Fatalf("wanted %s but got %s", w, o)
	}

	s = "_this_is__a___test_.png_"
	c = "_"
	o = trim(s, c)
	w = "this_is_a_test.png"
	if o != w {
		t.Fatalf("wanted %s but got %s", w, o)
	}
}

func BenchmarkTrim(b *testing.B) {
	s := "-this-is--a----test-.png-"
	c := "-"
	for i := 0; i < b.N; i++ {
		trim(s, c)
	}
}

func TestFixForShell(t *testing.T) {
	s := "--  __\a\b\f\n\r\v\t  test 123.45 _ () {} :; | *? <> \" '--.png-"
	o := FixForShell(s)
	w := "test_123.45.png"
	if o != w {
		t.Fatalf("wanted %s but got %s", w, o)
	}

	s = "--thí$ is (a really<bad>){file*/name}- .png "
	o = FixForShell(s)
	w = "thí_is_a_reallybadfilename.png"
	if o != w {
		t.Fatalf("wanted %s but got %s", w, o)
	}

	s = "    "
	o = FixForShell(s)
	w = "FN_NO_NAME"
	if o != w {
		t.Fatalf("wanted %q but got %q", w, o)
	}
}

func BenchmarkFixForShell(b *testing.B) {
	s := "--  __\a\b\f\n\r\v\t  test 123.45 _ () {} :; | *? <> \" '--.png-"
	for i := 0; i < b.N; i++ {
		FixForShell(s)
	}
}

func TestFixForURL(t *testing.T) {
	s := "--  __\a\b\f\n\r\v\t  áéíóúñ test 123.45 _ () {} :; | *? <> \" '--.png-"
	o := FixForURL(s)
	w := "aeioun-test-123.45.png"
	if o != w {
		t.Fatalf("wanted %s but got %s", w, o)
	}

	s = "--thí$ is (á réallý<bad>){file*/ñame}- .png "
	o = FixForURL(s)
	w = "thi-is-a-reallybadfilename.png"
	if o != w {
		t.Fatalf("wanted %s but got %s", w, o)
	}

	s = "    "
	o = FixForURL(s)
	w = "FN-NO-NAME"
	if o != w {
		t.Fatalf("wanted %q but got %q", w, o)
	}
}

func BenchmarkFixForURL(b *testing.B) {
	s := "--  __\a\b\f\n\r\v\t  áéíóúñ test 123.45 _ () {} :; | *? <> \" '--.png-"
	for i := 0; i < b.N; i++ {
		FixForURL(s)
	}
}

func TestTruncate(t *testing.T) {
	s := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer nec odio. Praesent libero. Sed cursus ante dapibus diam. Sed nisi. Nulla quis sem at nibh elementum imperdiet. Duis sagittis ipsum. Praesent mauris. Fusce nec tellus sed augue semper portab."
	o := truncate(s, MaxLenURL)
	ol := utf8.RuneCountInString(o)
	if !(ol >= MaxLenURL) {
		t.Fatalf("wanted %d but got %d", MaxLenURL, ol)
	}
}

func BenchmarkTruncate(b *testing.B) {
	s := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer nec odio. Praesent libero. Sed cursus ante dapibus diam. Sed nisi. Nulla quis sem at nibh elementum imperdiet. Duis sagittis ipsum. Praesent mauris. Fusce nec tellus sed augue semper portab."
	for i := 0; i < b.N; i++ {
		truncate(s, MaxLenURL)
	}
}
