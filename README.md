# go-assists

Series of helpers and structs for golang

### Work with Strings
-----------------------------------

##### CamelToSnakeCase(str string) string

 Convert a [camel](https://en.wikipedia.org/wiki/Camel_case) case string 
 to [snake](https://en.wikipedia.org/wiki/Snake_case) case 
 
 _Example:_ java_script -> JavaScript

##### CamelToSnakeCase(str string) string


### Work wit time.Time 

##### IsTimeZero(t *time.Time)
 Check if time is equals to 00:00:00

##### StartOfTheDay(t *time.Time) *time.Time
 Return new time.Time with time equals to start of the day : 00:00:00
 
##### EndOfTheDay(t *time.Time) *time.Time
 Return new time.Time with time equals end of the day : to 23:59:50
