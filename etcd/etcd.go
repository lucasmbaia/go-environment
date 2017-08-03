package etcd

import (
  "time"
  "reflect"
  "errors"
  "strconv"
  "encoding/json"

  client "github.com/coreos/etcd/client"
  "golang.org/x/net/context"
)

type Config struct {
  Endpoints []string
  TimeOut   int32
}

type Client struct {
  Client	  client.KeysAPI
  ContextTimeout  int
  Endpoints	  []string
  TimeOut	  int32
}

func (c Config) NewClient() (Client, error) {
  var (
    cliETCD client.Client
    err	    error
    cli	    Client
    keysAPI client.KeysAPI
  )

  if cliETCD, err = client.New(client.Config{
    Endpoints:	c.Endpoints,
    Transport:	client.DefaultTransport,
    HeaderTimeoutPerRequest: time.Duration(c.TimeOut) * time.Second,
  }); err != nil {
    return cli, err
  }

  keysAPI = client.NewKeysAPI(cliETCD)

  cli = Client{
    Client:	keysAPI,
    Endpoints:	c.Endpoints,
    TimeOut:	c.TimeOut,
  }

  return cli, err
}

func (c Client) Set(key string, value interface{}, tag bool) error {
  var (
    err	    error
    cancel  context.CancelFunc
    ctx     = context.Background()
  )

  ctx, cancel = context.WithTimeout(context.Background(), time.Duration(c.TimeOut) * time.Second)
  defer cancel()

  switch reflect.ValueOf(value).Kind() {
  case reflect.String:
    _, err = c.Client.Set(ctx, key, value.(string), nil)
  case reflect.Int:
    _, err = c.Client.Set(ctx, key, strconv.Itoa(value.(int)), nil)
  case reflect.Int8:
    _, err = c.Client.Set(ctx, key, strconv.Itoa(int(value.(int8))), nil)
  case reflect.Int16:
    _, err = c.Client.Set(ctx, key, strconv.Itoa(int(value.(int16))), nil)
  case reflect.Int32:
    _, err = c.Client.Set(ctx, key, strconv.Itoa(int(value.(int32))), nil)
  case reflect.Int64:
    _, err = c.Client.Set(ctx, key, strconv.Itoa(int(value.(int64))), nil)
  case reflect.Bool:
    _, err = c.Client.Set(ctx, key, strconv.FormatBool(value.(bool)), nil)
  case reflect.Float32, reflect.Float64:
    _, err = c.Client.Set(ctx, key, strconv.FormatFloat(value.(float64), 'f', -1, 64), nil)
  case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
    _, err = c.Client.Set(ctx, key, strconv.FormatUint(value.(uint64), 64), nil)
  case reflect.Uintptr:
  case reflect.Complex64:
  case reflect.Complex128:
  case reflect.Map:
    err = c.mapToString(ctx, key, value, tag)
  case reflect.Struct:
    err = c.mapToString(ctx, key, value, tag)
  case reflect.Ptr:
    err = c.mapToString(ctx, key, reflect.ValueOf(value).Elem().Interface(), tag)
  case reflect.Slice:
    err = c.mapToString(ctx, key, value, tag)
  default:
    return errors.New("Type is not supported")
  }

  return err
}

func (c Client) Get(key string, value interface{}, tag, defaultTag bool) error {
  var (
    err	      error
    cancel    context.CancelFunc
    ctx	      = context.Background()
    response  *client.Response
    v	      = reflect.ValueOf(value)
  )

  if v.Kind() != reflect.Ptr {
    return errors.New("Expected a pointer to a variable")
  }

  ctx, cancel = context.WithTimeout(context.Background(), time.Duration(c.TimeOut) * time.Second)
  defer cancel()

  if response, err = c.Client.Get(ctx, key, nil); err != nil {
    return err
  }

  switch v.Elem().Kind() {
  case reflect.String, reflect.Int, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Float64, reflect.Bool:
    if err = setField(v.Elem(), response.Node.Value); err != nil {
      return err
    }
  case reflect.Struct, reflect.Slice:
    if err = c.stringToMap(ctx, value, response.Node.Value, tag, defaultTag); err != nil {
      return err
    }
  default:
    return errors.New("Type is not supported")
  }

  return nil
}

func (c Client) stringToMap(ctx context.Context, value interface{}, response string, tag, defaultTag bool) error {
  var (
    err	    error
    objMap  map[string]interface{}
    v	    = reflect.ValueOf(value)
  )

  if tag {
    if err = json.Unmarshal([]byte(response), &objMap); err != nil {
      return err
    }

    switch v.Elem().Kind() {
    case reflect.Struct, reflect.Slice:
      if err = setField(v.Elem(), objMap); err != nil {
	return err
      }
    }
  } else {
    if err = json.Unmarshal([]byte(response), &value); err != nil {
      return err
    }
  }

  return nil
}

func setField(f reflect.Value, value interface{}) error {
  var (
    err		error
    i		interface{}
    t		reflect.Type
    r		reflect.Value
    objMap	map[string]interface{}
    objSlice	[]interface{}
    str		string
  )

  switch f.Type().Kind() {
  case reflect.String:
    f.Set(reflect.ValueOf(value.(string)))
  case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
    if i, err = strconv.ParseInt(strconv.FormatFloat(value.(float64), 'f', -1, 64), 10, 64); err != nil {
      return err
    }

    f.SetInt(i.(int64))
  case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
    if i, err = strconv.ParseUint(strconv.FormatFloat(value.(float64), 'f', -1, 64), 10, 64); err != nil {
      return err
    }

    f.SetUint(i.(uint64))
  case reflect.Float32, reflect.Float64:
    f.SetFloat(value.(float64))
  case reflect.Bool:
    f.SetBool(value.(bool))
  case reflect.Struct:
    objMap = value.(map[string]interface{})

    for i := 0; i < f.NumField(); i++ {
      t = f.Type()

      if tag, ok := t.FieldByName(t.Field(i).Name); ok {
	str = tag.Tag.Get("env")

	if _, ok := objMap[str]; ok {
	  if err = setField(f.Field(i), objMap[str]); err != nil {
	    return err
	  }
	}
      } else {
	return errors.New("Set tag env in variable of struct")
      }
    }
  case reflect.Slice:
    objSlice = value.([]interface{})

    for _, obj := range objSlice {
      r = reflect.New(f.Type().Elem())
      t = r.Type()

      switch reflect.ValueOf(obj).Kind() {
      case reflect.String, reflect.Int, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Float64, reflect.Bool:
	if err = setField(r.Elem(), obj); err != nil {
	  return err
	}
      default:
	objMap = obj.(map[string]interface{})

	for i := 0; i < r.Elem().NumField(); i++ {
	  if tag, ok := t.Elem().FieldByName(t.Elem().Field(i).Name); ok {
	    str = tag.Tag.Get("env")

	    if _, ok := objMap[str]; ok {
	      if err = setField(r.Elem().Field(i), objMap[str]); err != nil {
		return err
	      }
	    }
	  }
	}
      }

      f.Set(reflect.Append(f, r.Elem()))
    }
  }

  return nil
}

func (c Client) mapToString(ctx context.Context, key string, value interface{}, tag bool) error {
  var (
    body	  []byte
    err		  error
  )

  if tag {
    switch reflect.ValueOf(value).Kind() {
    case reflect.Struct, reflect.Slice:
      if body, err = json.Marshal(fieldToMap(reflect.ValueOf(value))); err != nil {
	return err
      }

      _, err = c.Client.Set(ctx, key, string(body), nil)
    }
  } else {
    if body, err = json.Marshal(value); err != nil {
      return err
    }

    _, err = c.Client.Set(ctx, key, string(body), nil)
  }

  return err
}

func fieldToMap(f reflect.Value) interface{} {
  var (
    mapKeysSlice  = make([]map[string]interface{}, 0)
    t		  reflect.Type
    r		  reflect.Value
    str		  string
  )

  switch f.Type().Kind() {
  case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Bool:
    return f.Interface()
  case reflect.Struct:
    var (
      mapKeys = make(map[string]interface{})
    )

    t = f.Type()

    for i := 0; i < f.NumField(); i++ {
      if tag, ok := t.FieldByName(t.Field(i).Name); ok {
	str = tag.Tag.Get("env")

	if str != "" {
	  mapKeys[str] = fieldToMap(f.Field(i))
	}
      }
    }

    return mapKeys
  case reflect.Slice:
    for i := 0; i < f.Len(); i++ {
      var (
	mapKeys = make(map[string]interface{})
      )

      r	= f.Index(i)
      t	= r.Type()

      for j := 0; j < r.NumField(); j++ {
	if tag, ok := t.FieldByName(t.Field(j).Name); ok {
	  str = tag.Tag.Get("env")

	  if str != "" {
	    mapKeys[str] = fieldToMap(r.Field(j))
	  }
	}
      }

      mapKeysSlice = append(mapKeysSlice, mapKeys)
    }

    return mapKeysSlice
  }

  return nil
}
