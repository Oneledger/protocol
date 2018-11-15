package vm

type OLVMService struct {
  Protocol string
  Port int
}

type Container int

type Args struct {
  From, Address, CallString string
  Value int
}

type Reply struct {
  Out, Ret string
}

type OLVMClient struct {
  Protocol string
  ServicePath string
}
