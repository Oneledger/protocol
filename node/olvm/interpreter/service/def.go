package service

type OLVMService int

type Args struct {
  From, Address, CallString string
  Value int
}

type Reply struct {
  Out, Ret string
}
