package controllers

import(
	"testing"
)

func TestGetPriorityI18n(t *testing.T) {
  cases := []struct{
    in int
    expected string
  }{
    {1, "Highest"},
    {2, "High"},
    {3, "Medium"},
    {4, "Low"},
    {5, "Lowest"},
    {10, "Unknown Status"},
    {22, "Unknown Status"},
  }

  for _, c := range cases {
    got := getPriorityI18n(c.in)
    if got != c.expected {
      t.Errorf("getPriorityI18n(%q) == %q, expected %q", c.in, got, c.expected)
    }
  }
}
