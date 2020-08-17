package main

import(
	"fmt"
)

//global var
var j string = "uarld"

func local_var(){

        //local var
	var i string = "henlo"
	
	fmt.Println("I am global j called from local_var():", j)
        fmt.Println(" ")	
	fmt.Println("I am local i called in main() from local_var():", i)
}

func main(){
	local_var()
	fmt.Println(" ")
	fmt.Println("I am global j called from main()", j)
	//fmt.Println(" ")
	//fmt.Println("now call local variable i: ", i)
}
