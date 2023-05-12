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

### Simple Concurrent Cache Usage

```go
package main

import (
	"fmt"
	"github.com/GalushkoArt/simpleCache"
	"time"
)

func main() {
	mapCache := simpleCache.NewConcurrentCache(100 * time.Millisecond)
	mapCache.Set("1", 1)
	mapCache.Set("3", 3)

	var value = *mapCache.Get("1")
	var value3 = *mapCache.Get("3")
	fmt.Println("Existed + 1:", value.(int)+1)
	fmt.Println("Value 3 at start:", value3)
	fmt.Println("Not Existed:", mapCache.Get("2"))

	fmt.Println("Deleted:", *mapCache.Delete("1"))
	fmt.Println("Not Deleted:", mapCache.Delete("2"))
	fmt.Println("After Delete:", mapCache.Get("1"))
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Value 3 after purify:", mapCache.Get("3"))
}
```

#### Result

```
Existed + 1: 2
Value 3 at start: 3
Not Existed: <nil>
Deleted: 1
Not Deleted: <nil>
After Delete: <nil>
Value 3 after purify: <nil>
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

	fmt.Println("Deleted + 1:", *mapCache.Delete("1")+1)
	fmt.Println("Not Deleted:", mapCache.Delete("2"))
	fmt.Println("After Delete:", mapCache.Get("1"))
}
```

#### Result

```
Existed + 1: 2
Not Existed: <nil>
Deleted: 2
Not Deleted: <nil>
After Delete: <nil>
```

### Generic Concurrent Cache Usage

```go
package main

import (
	"fmt"
	"github.com/GalushkoArt/simpleCache"
	"time"
)

func main() {
	mapCache := simpleCache.NewGenericConcurrentCache[int](100 * time.Millisecond)
	mapCache.Set("1", 1)
	mapCache.Set("3", 3)

	var value = *mapCache.Get("1")
	var value3 = *mapCache.Get("3")
	fmt.Println("Existed + 1:", value+1)
	fmt.Println("Value 3 at start:", value3)
	fmt.Println("Not Existed:", mapCache.Get("2"))

	fmt.Println("Deleted + 1:", *mapCache.Delete("1")+1)
	fmt.Println("Not Deleted:", mapCache.Delete("2"))
	fmt.Println("After Delete:", mapCache.Get("1"))
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Value 3 after purify:", mapCache.Get("3"))
}
```

#### Result

```
Existed + 1: 2
Value 3 at start: 3
Not Existed: <nil>
Deleted + 1: 2
Not Deleted: <nil>
After Delete: <nil>
Value 3 after purify: <nil>
```