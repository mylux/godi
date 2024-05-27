# godi

> [!Caution]
> This is still experimental, therefore limited and will eventually break. Think twice before using in critical systems

## Go Dependency Injector
Do you love Inversion of Control through dependency injection?  
Would you like that SpringBoot flavour in youg golang code? You won't find it here, but almost. Give this a try.

## Using godi
Install it with go get:
`go get github.com/mylux/godi/`

Use in your code:

```go
import "github.com/mylux/godi/container"

type Car interface(){
  ShiftGear()
}

type Civic struct {

}

func (c Civic) ShiftGear(){
    fmt.Println("CVT Module simulating gear...")
}

func WireStuff(){
  container.Wire(new(Car), new(Civic))
}

func main(){
  WireStuff()
  var c Car = container.Construct[Car]()
  c.ShiftGear()
}
```

This should output:
`"CVT Module simulating gear..."`

## AutoWire support
Using the `autowired` tag in a struct field allows the injector to automatically construct the object by using the `AutoWire` function
```go
import "github.com/mylux/godi/container"

type Car interface(){
  ShiftGear()
}

type Civic struct {

}

func (c Civic) ShiftGear(){
    fmt.Println("CVT Module simulating gear...")
}

type Person struct {
 DailyVehicle Car `autowired:""`
 RoadTripVehicle Car
}

func WireStuff(){
  container.Wire(new(Car), new(Civic))
}

func main(){
  p := Person{}
  container.AutoWire(&p)
  p.DailyVehicle.ShiftGear() // RoadTripVehicle will be unset, because it does not have the AutoWired tag
}

```

## Factory Function Wire
Another easier way of defining the implementation of an interface with godi is using a function that constructs the object. Example below:

```go
import "github.com/mylux/godi/container"

type Car interface(){
  ShiftGear()
}

type Civic struct {

}

func (c Civic) ShiftGear(){
    fmt.Println("CVT Module simulating gear...")
}

func createCar() Car{
  return Civic{}
}

func WireStuff(){
  container.Wire(createCar)
}

func main(){
  WireStuff()
  var c Car = container.Construct[Car]()
  c.ShiftGear()
}
```

This should output:
`"CVT Module simulating gear..."`

In the example above, the function return type identifies the interface that will be implemented. If godi has a function that can construct the implementation, it does not have to know what is the underlying type. It will simply run the function.

### Parametrized Factory
If your factory function has parameters, they can be automatically constructed too, as long as you have already wired them before. See the example below:

```go
import "github.com/mylux/godi/container"

type Car interface(){
  ShiftGear()
}

type Civic struct {

}

type MercedesMaybachS680 struct {
}

type OutdoorEvent interface {
  IsTrip() bool
}

type RoadTripEvent {
  
}

func (o RoadTripEvent) IsTrip() bool{
  return true
}

func (m MercedesMaybachS680) ShiftGear(){
    fmt.Println("9GTronic transmission selecting gear...")
}

func (c Civic) ShiftGear(){
    fmt.Println("CVT Module simulating gear...")
}

func createCar(e OutdoorEvent) Car{
  if e.IsTrip() {
    return MercedesMaybachS680{}
  }
  return Civic{}
}

func WireStuff(){
  container.Wire(new(OutdoorEvent), new(RoadTripEvent))
  container.Wire(createCar)
}

func main(){
  WireStuff()
  var c Car = container.Construct[Car]()
  c.ShiftGear()
}
```

This should output:
`"9GTronic transmission selecting gear..."`

In the example before, the OutdoorEvent is set to be implemented by RoadTripEvent, so the IsTrip() function will be always true. When constructing the Car object, the OutdoorEvent will be constructed first, using the type that was defined earlier and passed to the factory function that makes a car. As IsTrip() will be true, the car wil be Mercedes Maybach S680 (What a car, BTW!) and the gear will be shifted using its 9Gtronic transmission

## Constructing and wiring non-interfaced types
> **_NOTE:_**  New in version 0.2.0

You can provide constructors for non-interfaced wires too if your facory function requires a parameter.

If all you need is a simple `MyType{}` instanciation without anything special, godi can construct this with no wiring or configuration.

Why would you use it?

Imagine you Have the following struct:
```go
type NationalIdDocument struct {

}

type Car interface {
  ShiftGear()
}

type Civic struct {

}

func (c Civic) ShiftGear(){
  return "Shifting gear with CVT transmission..."
}

type Person struct {
  DailyDriver Car                 `autowired:""`
  RoadTrip    Car                 `autowired:""`
  IdDoc       NationalIdDocument  `autowired:""`
}

```

In the example above, you have the interface Car that is implemented by the Civic type, but NationalIdDocument is not an interface, but a struct.

Thanks to the ability of constructing the struct with a simple instantiation (as you would do with a `container.Construct[NationalIdDocument]()`), you can safely AutoWire the Person struct and it will construct the Person with two Civic cars (No Mercedes this time :disappointed:) and an empty NationalIdDocument in the IdDoc field, by doing this:
```go
func main(){
  container.Wire(new(Car), new(Civic))
  p := Person{}
  container.AutoWire(&p)
}
```

## Contributing
Interested in the project? Feel free to open issues, send Pull Requests and I will gladly look at those!
