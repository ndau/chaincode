package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

var section = `
	{
	fm := c.globalStore["functions"].(map[string]int)
		name := n.(string)
			ctr := c.globalStore["functionCounter"].(int)
  	fm[name] = ctr
	ctr   ++
	c.globalStore [ "functionCounter"] =ctr;
	return newFunctionDef(name, s)
	}
`

func formatSection(b []byte) ([]byte, error) {
	cmd := exec.Command("gofmt")
	cmd.Stdin = bytes.NewReader(b)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func doNothing(b []byte) ([]byte, error) {
	c := bytes.ToUpper(b)
	return c, nil
}

func main() {
	x, err := ParseReader("stdin", os.Stdin)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(x.([]byte)))
}
