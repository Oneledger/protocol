from sdk import *

addr_list = addresses()
_delegator = addr_list[0]
_delegate_amount = (int("2") * 10 ** 9)
_undelegate = Undelegate(_delegator, _delegate_amount)

if __name__ == "__main__":
    # todo delegate

    # undelegate

