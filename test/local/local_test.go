package main

import (
  "testing"
  "log"
  "github.com/lucasmbaia/go-environment/local"
)

type Service struct {
  ServiceName         string  `env:"SERVICE_NAME" envDefault:"grpc-orchestration"`
  EtcdURL             string  `env:"ETCD_URL" envDefault:"http://127.0.0.1:2379"`
  LinkerdURL          string  `env:"LINKERD_URL" envDefault:"127.0.0.1:4140"`
  CAFile              string  `env:"CA_FILE" envDefault:""`
  ServerNameAuthority string  `env:"SERVER_NAME_AUTHORITY" envDefault:""`
}

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

func TestLocalGetStruct(t *testing.T) {
  var (
    err	    error
    service Service
  )

  if err = local.Get("", &service, true, false); err != nil {
    log.Fatalf("Erro to get local env: ", err)
  }

  log.Println(service)
}
