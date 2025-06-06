package options

import (
	"encoding/json"
	"github.com/spf13/pflag"
	"github.com/yuanbaopig/app/fname"
	"reflect"
)

type Options struct {
	MongoSharding *MongoShardingOptions
}

func New() *Options {
	return &Options{
		MongoSharding: NewMongoShardingOptions(),
	}
}

// ApplyTo applies the run options to the method receiver and returns self.
func (o *Options) ApplyTo() error {
	return nil
}

func (o *Options) Flags() (fss fname.NamedFlagSets) {
	v := reflect.ValueOf(o).Elem()

	for i := 0; i < v.NumField(); i++ {
		if field, ok := v.Field(i).Interface().(opts); ok {
			field.AddFlags(fss.FlagSet(field.Name()))
		}
	}
	return
}

func (o *Options) String() string {
	data, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}

	return string(data)
}

func (o *Options) Complete() error {
	v := reflect.ValueOf(o).Elem()
	for i := 0; i < v.NumField(); i++ {
		if field, ok := v.Field(i).Interface().(opts); ok {
			if err := field.Complete(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *Options) Validate() []error {
	var errs []error
	v := reflect.ValueOf(o).Elem()
	for i := 0; i < v.NumField(); i++ {
		if field, ok := v.Field(i).Interface().(opts); ok {
			errs = append(errs, field.Validate()...)
		}
	}
	return errs
}

type OptsInterfaces[T any] interface {
	ApplyTo() T
	opts
}

type opts interface {
	AddFlags(*pflag.FlagSet)
	Validate() []error
	Complete() error
	Name() string
}
