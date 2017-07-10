package main

import (
  "log"
  "testing"
  "github.com/lucasmbaia/go-environment/etcd"
)

type Test struct {
  ValueA  string  `env:"VALUEA"`
  ValueB  int	  `env:"VALUEB"`
  ValueC  int32	  `env:"VALUEC"`
  ValueD  int64	  `env:"VALUED"`
  ValueE  float64 `env:"VALUEE"`
  ValueF  bool	  `env:"VALUEF"`
  ValueG  uint	  `env:"VALUEG"`
  ValueH  uint32  `env:"VALUEH"`
  ValueI  uint64  `env:"VALUEI"`
  TS	  TestStruct  `env:"TS"`
  SS	  []Test2Struct	`env:"SS"`
}

type TestStruct struct {
  ValueJ  string  `env:"VALUEJ"`
  ValueK  int	  `env:"VALUEK"`
}

type Test2Struct struct {
  ValueL  string      `env:"VALUEL"`
  ValueM  TestStruct  `env:"VALUEM"`
}

func TestEtcdConnect(t *testing.T) {
  var (
    config  etcd.Config
    err	    error
  )

  config = etcd.Config {
    Endpoints:	[]string{"127.0.0.1:2379"},
    TimeOut:	5,
  }

  if _, err = config.NewClient(); err != nil {
    log.Fatalf("Error to init new client: ", err)
  }
}

func TestEtcdSetStringVariable(t *testing.T) {
  var (
    client  etcd.Client
    err	    error
  )

  if client, err = connectEtcd(); err != nil {
    log.Fatalf("Error to init new client: ", err)
  }

  if err = client.Set("string", "teste", false); err != nil {
    log.Fatalf("Error to set etcd: ", err)
  }
}

func TestEtcdSetIntVariable(t *testing.T) {
  var (
    client  etcd.Client
    err	    error
    teste   = 10
  )

  if client, err = connectEtcd(); err != nil {
    log.Fatalf("Error to init new client: ", err)
  }

  if err = client.Set("int", teste, false); err != nil {
    log.Fatalf("Error to set etcd: ", err)
  }
}

func TestEtcdSetMapVariable(t *testing.T) {
  var (
    client  etcd.Client
    err	    error
    teste   = map[string]string{"teste": "teste"}
  )

  if client, err = connectEtcd(); err != nil {
    log.Fatalf("Error to init new client: ", err)
  }

  if err = client.Set("map", teste, false); err != nil {
    log.Fatalf("Error to set etcd: ", err)
  }
}

func TestEtcdSetSliceVariable(t *testing.T) {
  var (
    client  etcd.Client
    err	    error
    teste   = []string{"number1", "number2"}
  )

  if client, err = connectEtcd(); err != nil {
    log.Fatalf("Error to init new client: ", err)
  }

  if err = client.Set("slice_variable", teste, false); err != nil {
    log.Fatalf("Error to set etcd: ", err)
  }
}

func TestEtcdSetStruct(t *testing.T) {
  var (
    client  etcd.Client
    err	    error
    test    Test
    test2   []Test2Struct
  )

  test2 = append(test2, Test2Struct{ValueL: "valueL", ValueM: TestStruct{ValueJ: "valueJ", ValueK: 80}})
  test2 = append(test2, Test2Struct{ValueL: "valueL", ValueM: TestStruct{ValueJ: "valueJ", ValueK: 80}})

  test = Test{
    ValueA: "valueA",
    ValueB: 10,
    ValueC: 20,
    ValueD: 30,
    ValueE: 3.14,
    ValueF: true,
    ValueG: 40,
    ValueH: 50,
    ValueI: 60,
    TS: TestStruct{ValueJ: "valueJ", ValueK: 70},
    SS: test2,
  }

  if client, err = connectEtcd(); err != nil {
    log.Fatalf("Error to init new client: ", err)
  }

  if err = client.Set("struct", test, true); err != nil {
    log.Fatalf("Error to set etcd: ", err)
  }
}

func TestEtcdSetSliceStruct(t *testing.T) {
  var (
    client	etcd.Client
    err		error
    testSlice   []Test
    test	Test
    test2	[]Test2Struct
  )

  test2 = append(test2, Test2Struct{ValueL: "valueL", ValueM: TestStruct{ValueJ: "valueJ", ValueK: 80}})
  test2 = append(test2, Test2Struct{ValueL: "valueL", ValueM: TestStruct{ValueJ: "valueJ", ValueK: 80}})

  test = Test{
    ValueA: "valueA",
    ValueB: 10,
    ValueC: 20,
    ValueD: 30,
    ValueE: 3.14,
    ValueF: true,
    ValueG: 40,
    ValueH: 50,
    ValueI: 60,
    TS: TestStruct{ValueJ: "valueJ", ValueK: 70},
    SS: test2,
  }

  testSlice = append(testSlice, test)
  testSlice = append(testSlice, test)

  if client, err = connectEtcd(); err != nil {
    log.Fatalf("Error to init new client: ", err)
  }

  if err = client.Set("slice", testSlice, false); err != nil {
    log.Fatalf("Error to set etcd: ", err)
  }
}

func TestEtcdGetStringVariable(t *testing.T) {
  var (
    client  etcd.Client
    err	    error
    teste   string
  )

  if client, err = connectEtcd(); err != nil {
    log.Fatalf("Error to init new client: ", err)
  }

  if err = client.Get("string", &teste, false, false); err != nil {
    log.Fatalf("Error to set etcd: ", err)
  }

  log.Println(teste)
}

func TestEtcdGetIntVariable(t *testing.T) {
  var (
    client  etcd.Client
    err	    error
    teste   int
  )

  if client, err = connectEtcd(); err != nil {
    log.Fatalf("Error to init new client: ", err)
  }

  if err = client.Get("int", &teste, false, false); err != nil {
    log.Fatalf("Error to set etcd: ", err)
  }

  log.Println(teste)
}

func TestEtcdGetStruct(t *testing.T) {
  var (
    client  etcd.Client
    err	    error
    teste   Test
  )

  if client, err = connectEtcd(); err != nil {
    log.Fatalf("Error to init new client: ", err)
  }

  if err = client.Get("struct", &teste, true, false); err != nil {
    log.Fatalf("Error to set etcd: ", err)
  }

  log.Println(teste)
}

func TestEtcdGetSliceStruct(t *testing.T) {
  var (
    client  etcd.Client
    err	    error
    teste   []Test
  )

  if client, err = connectEtcd(); err != nil {
    log.Fatalf("Error to init new client: ", err)
  }

  if err = client.Get("slice", &teste, false, false); err != nil {
    log.Fatalf("Error to set etcd: ", err)
  }

  log.Println(teste)
}

func connectEtcd() (etcd.Client, error) {
  var (
    config  etcd.Config
  )

  config = etcd.Config {
    Endpoints:	[]string{"http://127.0.0.1:2379"},
    TimeOut:	5,
  }

  return config.NewClient()
}
