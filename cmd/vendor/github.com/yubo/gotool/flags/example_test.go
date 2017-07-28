/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package flags_test

import (
	"flag"
	"fmt"

	"github.com/yubo/gotool/flags"
)

var (
	b        bool
	s        string
	sb1, sb2 bool
	ss1, ss2 string
)

func reset() {
	b = false
	sb1 = false
	sb2 = false
	s = "hello"
	ss1 = "hello"
	ss2 = "hello"
}

func init() {

	flag.BoolVar(&b, "b", false, "a bool")
	flag.StringVar(&s, "s", "hello", "a string")

	cmd1 := flags.NewCommand("cmd1", "cmd1,usage", flag.ExitOnError)
	cmd1.BoolVar(&sb1, "b", false, "a bool")
	cmd1.StringVar(&ss1, "s", "hello", "a string")

	cmd2 := flags.NewCommand("cmd2", "cmd2,usage", flag.ExitOnError)
	cmd2.BoolVar(&sb2, "b", false, "a bool")
	cmd2.StringVar(&ss2, "s", "hello", "a string")

}

func ExampleFlag() {

	reset()
	flags.Parse()

	fmt.Println("b:", b)
	fmt.Println("s:", s)
	// Output:
	// b: false
	// s: hello
}

func ExampleFlags1() {
	reset()
	flags.CommandLine.Parse([]string{
		"-b", "-s", "world",
		"cmd1", "-b", "-s", "007",
	})
	fmt.Println("b:", b)
	fmt.Println("s:", s)
	fmt.Println("cmd1.b:", sb1)
	fmt.Println("cmd1.s:", ss1)
	fmt.Println("cmd2.b:", sb2)
	fmt.Println("cmd2.s:", ss2)
	// Output:
	// b: true
	// s: world
	// cmd1.b: true
	// cmd1.s: 007
	// cmd2.b: false
	// cmd2.s: hello
}

func ExampleFlags2() {
	reset()
	flags.CommandLine.Parse([]string{
		"-s", "hello,world",
		"cmd2", "-b", "-s", "007",
	})
	fmt.Println("b:", b)
	fmt.Println("s:", s)
	fmt.Println("cmd1.b:", sb1)
	fmt.Println("cmd1.s:", ss1)
	fmt.Println("cmd2.b:", sb2)
	fmt.Println("cmd2.s:", ss2)
	// Output:
	// b: false
	// s: hello,world
	// cmd1.b: false
	// cmd1.s: hello
	// cmd2.b: true
	// cmd2.s: 007
}
