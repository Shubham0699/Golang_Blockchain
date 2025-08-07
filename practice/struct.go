package main

import  "fmt"
//import  "bytes"

type rectangle struct {
	length int
	breadth int
}

func (r rectangle) area() int{

	return r.length * r.breadth

	
}


func main(){
	r:= rectangle{
		length:10,
		breadth:5,
	}
	fmt.Println("the area is:" , r.area())

	b:=[]byte("hello")
	fmt.Println(b)

}