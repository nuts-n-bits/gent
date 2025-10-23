package main

import "fmt"

func main() {
	var aruploadmix Upload3;
	a, b := aruploadmix.Parse(`
	
arinit -ns zhwp
arsync -i mapjas3981 "-n 20251992"
addrev -r 299821 -t 2015-25-19T24:53:45 -u 424353 -un "" -s a -s b -c fjreopgjeipgrt
	`);

	fmt.Printf("ERR: %#v\n MIS: %#v\n PARSED: %#v\n\n", b, a, aruploadmix);


	fmt.Print(aruploadmix.Write());

}

