package async

type (
	IOperation interface {
		Execute() (interface{}, error)
	}

	FuncOperation func() (interface{}, error)
)

func NewFuncOperation(f func() (interface{}, error)) FuncOperation {
	return FuncOperation(f)
}

func (f FuncOperation) Execute() (interface{}, error) {
	return f()
}
