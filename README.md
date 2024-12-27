## Nox programming language

Learning about lexers, parsers and programming languages in general.

### Capabilities
Nox is capable for this much as of right now:-

```c
fn retVal() {
    // return statements and function calls
    return 100 + 10;
}

fn main() {
    // assignments
    let a: int = retVal();
    // printing to stdout
    print(a);

    // mutable variables
    a = 50;

    // expressions and booleans
    print(a / 2);
    print(a >= 5 && a <= 50);

    // strings
    let str: string = "something";
    print(str);

    // If-Else [Control statements]
    let x: int = 8;
    if (x == 10) {
        print("IF WORKED");
    } else if (x == 9) {
        print("else if worked");
    } else {
        print("else worked");
    }

    // For
    for (let i: int = 0; i < 10; i = i + 1;) {
        print("Something");
        print(i);
    }

    return;
}
```

#### Resources
Some resources if you wish to read around these topics:-

- [Recursive expression parser](https://www.stroustrup.com/Programming/calculator00.cpp)
- [Pratt Parsing](https://matklad.github.io/2020/04/13/simple-but-powerful-pratt-parsing.html)
- [QBE backend](https://c9x.me/compile/)
