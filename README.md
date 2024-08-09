# Jelly


	Version: 0.0.1
	Date: 09/08/2024 UK

	Jelly is just another way of storing variables in a readable file
	basically another version of a .env file but in a single file and different format

	Jelly files are basic, currently only support variable and addition

	I opted to have any vars or strings on the same line to just append
	No need to have a + operator as its implied

	This is very basic right now and only allows strings or variables
	If a number is needed, it must be done as a string

## Example Usage


### -- data.jelly
```
greeting = "Hello"
name = "Jelly"

full_greeting = @greeting ", " @name "!"
```
### -- main.go
```go
func main() {
	store := NewStore("data.jelly")

	fg := store.Get("full_greeting")
	fmt.Println("Full greeting:", fg)
}
```

### Output
    Hello, Jelly!