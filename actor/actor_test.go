package actor

import (
	"testing"
)

func TestActor(t *testing.T) {
	a1 := NewActor(10)
	a2 := NewActor(20)
	a3 := NewActor(30)
	a4 := NewActor(40)
	a5 := NewActor(50)

	a1.Next()
	if len(a1.next) != 0 {
		t.Fatal("Calling Actor.Next() without arguments shouldn't do anything")
	}
	a1.Next(nil)
	if len(a1.next) != 0 {
		t.Fatal("Calling Actor.Next(nil) shouldn't do anything")
	}

	a1.Next(a2, a4)
	a2.Next(a3)
	a4.Next(a5)

	if len(a1.next) != 2 {
		t.Fatal("Adding 2 actors via Actor.Next(...) should leave 3 actors in Actor.next")
	}
	if len(a2.next) != 1 {
		t.Fatal("Adding 1 actor via Actor.Next(...) should leave 1 actor in Actor.next")
	}
	if len(a4.next) != 1 {
		t.Fatal("Adding 1 actor via Actor.Next(...) should leave 1 actor in Actor.next")
	}
}
