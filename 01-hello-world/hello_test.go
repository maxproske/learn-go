package main

import "testing"

func TestHello(t *testing.T) {
	// Subtests for different scenarios
	t.Run("saying hello to people", func(t *testing.T) {
		got := Hello("Max", "English")
		want := "Hello, Max"

		assertCorrectMessage(t, got, want)
	})
	t.Run("say hello to empty string", func(t *testing.T) {
		got := Hello("", "English")
		want := "Hello, World"

		assertCorrectMessage(t, got, want)
	})
	t.Run("in Spanish", func(t *testing.T) {
		got := Hello("Elodie", "Spanish")
		want := "Hola, Elodie"
		assertCorrectMessage(t, got, want)
	})
	t.Run("in French", func(t *testing.T) {
		got := Hello("James", "French")
		want := "Bonjour, James"
		assertCorrectMessage(t, got, want)
	})
}

func assertCorrectMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
