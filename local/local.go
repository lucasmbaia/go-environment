package local

import (
  "reflect"
  "os"
  "errors"
  "strconv"
)

func Set(key string, value interface{}, setTag bool) error {
  var (
    err	error
    v	reflect.Value
    t	reflect.Type
    str	string
  )

  switch reflect.ValueOf(value).Kind() {
  case reflect.String:
    err = os.Setenv(key, value.(string))
  case reflect.Int:
    err = os.Setenv(key, strconv.Itoa(value.(int)))
  case reflect.Int8:
    err = os.Setenv(key, strconv.Itoa(int(value.(int8))))
  case reflect.Int16:
    err = os.Setenv(key, strconv.Itoa(int(value.(int16))))
  case reflect.Int32:
    err = os.Setenv(key, strconv.Itoa(int(value.(int32))))
  case reflect.Int64:
    err = os.Setenv(key, strconv.Itoa(int(value.(int64))))
  case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
    err = os.Setenv(key, strconv.FormatUint(value.(uint64), 64))
  case reflect.Float32, reflect.Float64:
    err = os.Setenv(key, strconv.FormatFloat(value.(float64), 'f', -1, 64))
  case reflect.Bool:
    err = os.Setenv(key, strconv.FormatBool(value.(bool)))
  case reflect.Struct:
    v =	reflect.ValueOf(value)
    t = v.Type()

    if setTag {
      for i := 0; i < v.NumField(); i++ {
	if tag, ok := t.FieldByName(t.Field(i).Name); ok {
	  str = tag.Tag.Get("env")

	  if str != "" {
	    if err = Set(str, v.Field(i).Interface(), setTag); err != nil {
	      return err
	    }
	  }
	}
      }
    } else {
      for i := 0; i < v.NumField(); i++ {
	if err = Set(t.Field(i).Name, v.Field(i).Interface(), setTag); err != nil {
	  return err
	}
      }
    }

  case reflect.Ptr:
    err = Set(key, reflect.ValueOf(value).Elem().Interface(), setTag)
  default:
    return errors.New("Type is not supported")
  }

  return err
}

func Get(key string, value interface{}, tag, defaultTag bool) error {
  var (
    err	error
    env	string
    v	= reflect.ValueOf(value)
  )

  if v.Kind() != reflect.Ptr {
    return errors.New("Expected a pointer to a variable")
  }

  if v.Elem().Kind() != reflect.Struct {
    if env = os.Getenv(key); env == "" {
      return errors.New("Environment does not exist")
    }
  }

  if err = setField(v.Elem(), env, tag, defaultTag); err != nil {
    return err
  }

  return nil
}

func setField(f reflect.Value, value interface{}, getTag, defaultTag bool) error {
  var (
    err		error
    i		interface{}
    v		reflect.Value
    t		reflect.Type
    env		string
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
    v =	reflect.ValueOf(value)
    t = v.Type()

    if getTag {
      for i := 0; i < v.NumField(); i++ {
	if tag, ok := t.FieldByName(t.Field(i).Name); ok {
	  str = tag.Tag.Get("env")

	  if str != "" {
	    if env = os.Getenv(str); env != "" {
	      if err = setField(v.Field(i), env, getTag, defaultTag); err != nil {
		return err
	      }
	    } else {
	      if defaultTag {
		if tag.Tag.Get("envDefault") != "" {
		  if err = setField(v.Field(i), tag.Tag.Get("envDefault"), getTag, defaultTag); err != nil {
		    return err
		  }
		}
	      }
	    }
	  } else {
	    if defaultTag {
	      if tag.Tag.Get("envDefault") != "" {
		if err = setField(v.Field(i), tag.Tag.Get("envDefault"), getTag, defaultTag); err != nil {
		  return err
		}
	      }
	    }
	  }
	}
      }
    } else {
      for i := 0; i < v.NumField(); i++ {
	if env = os.Getenv(t.Field(i).Name); env != "" {
	  if err = setField(v.Field(i), env, getTag, defaultTag); err != nil {
	    return err
	  }
	} else {
	  if defaultTag {
	    if tag, ok := t.FieldByName(t.Field(i).Name); ok {
	      if tag.Tag.Get("envDefault") != "" {
		if err = setField(v.Field(i), tag.Tag.Get("envDefault"), getTag, defaultTag); err != nil {
		  return err
		}
	      }
	    }
	  }
	}
      }
    }
  case reflect.Ptr:
    if err = setField(f.Elem(), value, getTag, defaultTag); err != nil {
      return err
    }
  default:
    return errors.New("Type is not supported")
  }

  return nil
}
