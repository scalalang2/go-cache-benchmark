package main

type Generator interface {
	Name() string
	Next() string
}
