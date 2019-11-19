package command

// Haskell: data Command = Arithmetic String| Pop String Int|  Push String Int

type Command interface {
	isCommand()
}

type Arithmetic struct {
	Name string
}

func (_ Arithmetic) isCommand() {}

type Pop struct {
	Segment string
	Index   int
}

func (_ Pop) isCommand() {}

type Push struct {
	Segment string
	Index   int
}

func (_ Push) isCommand() {}

type Label struct {
	Name string
}

func (_ Label) isCommand() {}

type Goto struct {
	Label string
}

func (_ Goto) isCommand() {}

type IFGoto struct {
	Label string
}

func (_ IFGoto) isCommand() {}

type UndefinedCommand struct {
	Label string
}

func (_ UndefinedCommand) isCommand() {}
