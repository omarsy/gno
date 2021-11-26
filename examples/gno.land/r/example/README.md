This is a demo of Gno smart contract programming.  This document was
constructed by Gno. To see how it was done, follow the steps below.

The smart contract files that were uploaded to make this
possible can be found here:
https://github.com/gnolang/gno/tree/master/examples/gno.land

## add test account

> make
> ./build/gnokey add test1 --recover

Use this mnemonic:
> source bonus chronic canvas draft south burst lottery vacant surface solve popular case indicate oppose farm nothing bullet exhibit title speed wink action roast

## start gnoland validator.

> ./build/gnoland

(This can be reset with `make reset`).

## sign an addpkg (add avl package) transaction.

```bash
./build/gnokey maketx addpkg test1 --pkgpath "gno.land/p/avl" --pkgdir "examples/gno.land/p/avl" --deposit 100gnot --gas-fee 1gnot --gas-wanted 2000000 > addpkg.avl.unsigned.txt
./build/gnokey query "auth/accounts/g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5"
./build/gnokey sign test1 --txpath addpkg.avl.unsigned.txt --chainid "testchain" --number 0 --sequence 0 > addpkg.avl.signed.txt
./build/gnokey broadcast addpkg.avl.signed.txt
```

## sign an addpkg (add dom package) transaction.

```bash
./build/gnokey maketx addpkg test1 --pkgpath "gno.land/p/dom" --pkgdir "examples/gno.land/p/dom" --deposit 100gnot --gas-fee 1gnot --gas-wanted 2000000 > addpkg.dom.unsigned.txt
./build/gnokey query "auth/accounts/g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5"
./build/gnokey sign test1 --txpath addpkg.dom.unsigned.txt --chainid "testchain" --number 0 --sequence 1 > addpkg.dom.signed.txt
./build/gnokey broadcast addpkg.dom.signed.txt
```

## sign an addpkg (add example realm) transaction.

```bash
./build/gnokey maketx addpkg test1 --pkgpath "gno.land/r/example" --pkgdir "examples/gno.land/r/example" --deposit 100gnot --gas-fee 1gnot --gas-wanted 2000000 > addrealm.unsigned.txt
./build/gnokey query "auth/accounts/g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5"
./build/gnokey sign test1 --txpath addrealm.unsigned.txt --chainid "testchain" --number 0 --sequence 2 > addrealm.signed.txt
./build/gnokey broadcast addrealm.signed.txt
```

## sign a exec (contract call) transaction.

```bash
./build/gnokey maketx exec test1 --pkgpath "gno.land/r/example" --func AddPost --args "Gno Demo" --args#file "./examples/gno.land/r/example/README.md" --gas-fee 1gnot --gas-wanted 2000000 > addpage.unsigned.txt
./build/gnokey query "auth/accounts/g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5"
./build/gnokey sign test1 --txpath addpage.unsigned.txt --chainid "testchain" --number 0 --sequence 3 > addpage.signed.txt
./build/gnokey broadcast addpage.signed.txt
```

## render page with ABCI query (evalquery).

```bash
./build/gnokey query "vm/qeval" --data "gno.land/r/example
Render()"
```