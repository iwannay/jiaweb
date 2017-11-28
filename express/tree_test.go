package express

import (
	"testing"
)

var tree = NewTree()

func TestInsertChild(t *testing.T) {
	tree.insertChild("/hello/go/toscrr", func(t *HttpCtx) {})
	tree.insertChild("/a/b/c/d", func(t *HttpCtx) {})
	tree.insertChild("/", func(t *HttpCtx) {})
	tree.insertChild("/hello/go/haah", func(t *HttpCtx) {})
	tree.insertChild("/hello/test", func(t *HttpCtx) {})
	tree.insertChild("/home/like-a/", func(t *HttpCtx) {})
	tree.insertChild("/home/like-b/", func(t *HttpCtx) {})
	tree.insertChild("/home/like-a/haha", func(t *HttpCtx) {})
	tree.insertChild("/home/like-a/good", func(t *HttpCtx) {})
	tree.insertChild(`/home/<a1:\d>/good`, func(t *HttpCtx) {})
	tree.insertChild("/home/<a:.{3}>/good", func(t *HttpCtx) {})
	tree.insertChild("/hohah/go", func(t *HttpCtx) {})
	tree.insertChild("/hohah/goto", func(t *HttpCtx) {})
	tree.insertChild(`/home/<a:\d>/lala`, func(t *HttpCtx) {})
	tree.insertChild(`/home/<age:\d>/<name:\d>/<like:.*>`, func(t *HttpCtx) {})
	tree.insertChild(`/student-<name:.*>-<age:\d*>`, func(t *HttpCtx) {})
	tree.insertChild(`/student/你好`, func(t *HttpCtx) {})
	// str, err := json.MarshalIndent(tree, "", "    ")
	// if err != nil {
	// 	t.Fatalf("%v", err)
	// }
	// f, err := os.OpenFile("./test.json", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	// if err != nil {
	// 	t.Fatalf("%v", err)
	// }
	// if _, err := f.WriteString(string(str)); err != nil {
	// 	t.Fatalf("%v", err)
	// }
}

func TestGetNode(t *testing.T) {
	if n, _, _ := tree.getNode("/home/like-a/"); n == nil {
		t.Fatal("cannot found /home/like-a/")
	}
	if n, _, fpath := tree.getNode("/haha"); n != nil {
		t.Fatal("fullpath /haha find " + fpath)
	}

	if n, p, _ := tree.getNode("/student-john-16"); n != nil {
		if v, ok := p["name"]; !ok || v != "john" {
			t.Fatal("fullpath /student-john-16 name err!")
		}
		if v, ok := p["age"]; !ok || v != "16" {
			t.Fatal("fullpath /student-john-16 age err!")
		}
	}

	if n, p, _ := tree.getNode("/home/2/3/apple"); n != nil {
		if v, ok := p["name"]; !ok || v != "3" {
			t.Fatal("fullpath /home/2/3/apple name err!")
		}
		if v, ok := p["age"]; !ok || v != "2" {
			t.Fatal("fullpath /home/2/3/apple age err!")
		}
		if v, ok := p["like"]; !ok || v != "apple" {
			t.Fatal("fullpath /home/2/3/apple like err!")
		}
	} else {
		t.Fatal("fullpath /home/2/3/apple cannot found!")
	}

	t.Log(tree.getNode("/home/232/good"))

}
