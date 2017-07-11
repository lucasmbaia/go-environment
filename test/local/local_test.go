package main

import (
  "testing"
  "log"
  "github.com/lucasmbaia/go-environment/local"
)

func TestLocalSetStringVariable(t *testing.T) {
  var (
    err	error
  )

  if err = local.Set("STRING", "teste", false); err != nil {
    log.Fatalf("Error to set local env: ", err)
  }
}

func TestLocalGetStringVariable(t *testing.T) {
  var (
    err	  error
    test  string
  )

  if err = local.Get("STRING", &test, false, false); err != nil {
    log.Fatalf("Erro to get local env: ", err)
  }

  log.Println(test)
}
