import time

from sdk import *

addr_list = addresses()
_delegator = addr_list[0]
_delegate_amount = int("10") * 10 ** 9
_undelegate_amount_1 = int("1") * 10 ** 9
_undelegate_amount_2 = int("2") * 10 ** 9
_undelegate_amount_3 = int("3") * 10 ** 9
_undelegate = Undelegate(_delegator)

if __name__ == "__main__":
    # todo delegate

    # undelegate
    _undelegate.send_tx(_undelegate_amount_1)
    time.sleep(10)
    _undelegate.send_tx(_undelegate_amount_2)
    _undelegate.send_tx(_undelegate_amount_3)
    # query result
    result = _undelegate.query_undelegate()

    check_query_undelegate(result, 2, )
    print bcolors.OKGREEN + "#### Undelegation Test Succeded" + bcolors.ENDC
