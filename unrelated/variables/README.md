
<b> G0lang Variables </b>

Small example, named var1.go:

 ```
package main

import "fmt"

func main(){
	var i string = "henlo" //declare and initialize variable
	j := "uarld"           // use the short variable declaration
	fmt.Println(i, j)
}
 ```
 
Build it, and run it:

```
root@kr03nen:/home/gog0#  go build var1.go 
root@kr03nen:/home/gog0#  ./var1 
henlo uarld
```
 

Now, let's compile the program so that we can use it with gdb:

```
go build -gcflags "-N -l" -o debug_var var1.go
```
Start gdb, and execute instructions (comment added for better understanding):

```
root@kr03nen:/home/gog0# gdb -q ./debug_var 
(gdb) #list program lines
(gdb) list
1	package main
2	
3	import "fmt"
4	
5	func main(){
6		var i string = "henlo"
7		j := "uarld"
8		fmt.Println(i, j)
9	}
(gdb) #breakpoint at line 5, where function main starts
(gdb) break 5
Breakpoint 1 at 0x498ea0: file /home/gog0/var1.go, line 5.
(gdb) #start the program
(gdb) run
Starting program: /home/gog0/debug_var 
[New LWP 7100]
[New LWP 7101]
[New LWP 7102]
[New LWP 7103]
[New LWP 7104]

Thread 1 "debug_var" hit Breakpoint 1, main.main () at /home/gog0/var1.go:5
5	func main(){
(gdb) #next line - 6
(gdb) step
6		var i string = "henlo"
(gdb) #find the type of i
(gdb) whatis i
type = struct string
(gdb) #next line - 7
(gdb) s
7		j := "uarld"
(gdb) #find the type of j
(gdb) whatis j
type = struct string
(gdb) s
8		fmt.Println(i, j)
(gdb) s
runtime.convTstring (val=..., x=<optimized out>) at /usr/local/go/src/runtime/iface.go:370
370	func convTstring(val string) (x unsafe.Pointer) {
(gdb) continue
Continuing.
henlo uarld
[LWP 7104 exited]
[LWP 7103 exited]
[LWP 7102 exited]
[LWP 7101 exited]
[LWP 7100 exited]
[Inferior 1 (process 7096) exited normally]
(gdb) #exit gdb
(gdb) q
root@kr03nen:/home/gog0# 

```
<br></br>

<i> What did we just notice in the gdb? </i>
<br>
Both variables i, respectively j have been found as string type; however, we have a variable i, declared as type "string". 
And then, we don't specify any type for j variable - here, Golang will deduce the type of j variable (hence, using the ":=")


<b> Variables scope - Local and global </b>

Just like in other programming languages (Java, c/c++, python, etc), the scope of variables is existent even in Golang.
As in other programming language, we encounter global and local variables. <br>
<i>Short recap: <br>
Global variables - variables that can be called anywhere in a program <br>
Local variables - variables with strict access / called only in specific parts of a program </i>

Small example, program g_l.go

```
package main

import(
	"fmt"
)

//global
var j string = "uarld"

func local_var(){

        //local
	var i string = "henlo"
	
	fmt.Println(" ")
	fmt.Println("I am global j called from local_var():", j)
        fmt.Println(" ")	
	fmt.Println("I am local i called in main() from local_var():", i)
}

func main(){
	local_var()
	fmt.Println(" ")
	fmt.Println("I am global j called from main():", j)
}
```

Let's simply run it:

```
root@kr03nen:/home/gog0# go run g_l.go 

I am global j called from local_var(): uarld
 
I am local i called in main() from local_var(): henlo
 
I am global j called from main(): uarld
```

Now, let's change a bit the program and add two more lines at the ending, in main() function:

```
package main

import(
	"fmt"
)

var j string = "uarld"

func local_var(){

	var i string = "henlo"
	
	fmt.Println("I am global j called from local_var():", j)
        fmt.Println(" ")	
	fmt.Println("I am local i called in main() from local_var():", i)
}

func main(){
	local_var()
	fmt.Println(" ")
	fmt.Println("I am global j called from main()", j)
	fmt.Println(" ")
	fmt.Println("now call local variable i: ", i)
}

```

Let's run it:

```
root@kr03nen:/home/gog0# go run g_l.go 
# command-line-arguments
./g_l.go:23:45: undefined: i
```

Obviously, we have obtained the "undefined: i" error, since we are calling variable i outside its own function, local_var()


