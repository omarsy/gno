package demo

func Render(_ string) string {
	return `# Gno Demo Packages

Welcome to the "demo" namespace. Here you will find realms that demonstrate 
how the Gno language can be used. Check out some examples below:

- [Art](/r/demo/art)
- [AtomicSwap](/r/demo/atomicswap)
- [Userbook](/r/demo/userbook)
- [Boards](/r/demo/boards)
- [Counter](/r/demo/counter)
- GRC:
	- [Foo20](/r/demo/foo20)
	- [Bar20](/r/demo/bar20)
	- [GRC20 Registry](/r/demo/grc20reg)
	- [GRC20 Factory](/r/demo/grc20factory)
	- [Foo721](/r/demo/foo721)
- ...


> This is a non-exhaustive list. Check out the [Gno monorepo](https://github.com/gnolang/gno/tree/master/examples/gno.land/r/demo)
for more examples.

`
}
