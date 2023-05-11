Go Simple Cache Example
==========================

Simple example how to implement in-memory cache with Go

You can use simple cache and parametrized with generics cache

## Usage

### Simple Cache Usage

```go
package main

import (
	"fmt"
	"github.com/GalushkoArt/simpleCache"
)

func main() {
	mapCache := simpleCache.NewMapCache()
	mapCache.Set("1", 1)

	var value = *mapCache.Get("1")
	fmt.Println("Existed + 1:", value.(int)+1)
	fmt.Println("Not Existed:", mapCache.Get("2"))

	fmt.Println("Deleted:", *mapCache.Delete("1"))
	fmt.Println("Not Deleted:", mapCache.Delete("2"))
	fmt.Println("After Delete:", mapCache.Get("1"))
}
```

#### Result

```
Existed + 1: 2
Not Existed: <nil>
Deleted: 1
Not Deleted: <nil>
After Delete: <nil>
```

### Generic Cache Usage

```go
package main

import (
	"fmt"
	"github.com/GalushkoArt/simpleCache"
)

func main() {
	mapCache := simpleCache.NewGenericMapCache[int]()
	mapCache.Set("1", 1)

	fmt.Println("Existed + 1:", *mapCache.Get("1")+1)
	fmt.Println("Not Existed:", mapCache.Get("2"))

	fmt.Println("Deleted:", *mapCache.Delete("1"))
	fmt.Println("Not Deleted:", mapCache.Delete("2"))
	fmt.Println("After Delete:", mapCache.Get("1"))
}
```

#### Result

```
Existed + 1: 2
Not Existed: <nil>
Deleted: 1
Not Deleted: <nil>
After Delete: <nil>
```