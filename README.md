# DI

[![Go Reference](https://pkg.go.dev/badge/github.com/andresmeireles/di.svg)](https://pkg.go.dev/github.com/andresmeireles/di)
![WorkflowBadge](https://github.com/andresmeireles/di/actions/workflows/go.yml/badge.svg)

A reflection based dependency injection and service locator.

## How Install

`go get -u github.com/andresmeireles/di`

## How use

Ha duas peças fundamentais para ao usar essa lib, o conceito de dependencia e o de container, a dependencia é uma struct que implementa a interface `Dependency` e o container é a implementacao de varias dependencias. Entao, para usar esta lib é necessario ter um container e ele receber dependencias.

### Criando uma dependencia

Temos o seguinte struct 

```
package mypkg

type Dep struct {
    Name string
}
```

Essa struct precisa virar um depedencia, o que pode ser feito de dois modos, com uma `NamedDependency` ou uma `TypedDependency`, uma vai receber uma string como identificador e outra recebera um tipo como identificador. Em ambas um parametro com uma funcao que criar essa struct deve ser recebida.

```
\\ named
d := di.NewNamedDependency("mypkg.dep", func () {return Dep{"Andre"}})

\\ on named dependency name can be any string
d := di.NewNamedDependency("mydep", func () {return Dep{"Andre"}})

\\ or typed
d := di.NewTypedDependency[mypkg.Dep](func () (return Dep{"Andre}))
```

Essas dois modos de criar uma implementacao da interface `Dependency`. Com uma dependencia criada ela deve ser adicionado no container.

````
container := new(di.Container)

err := container.Add(d)

if err != nil {
    panic("cannot add")
}

// it can be invoked in two ways:
// with string, you must know what name you give to dependency
dep, err := container.Get("mypkg.Dep") // or container.Get("mydep")

// With helper
// when use this helper, the it will be called always by type
dep, err := di.Get[mypkg.Dep](container)
````


