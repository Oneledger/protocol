import sys


def check_query_undelegate(result, pending_count_expected, matured_expected):
    if not result['height']:
        sys.exit(-1)
    if len(result['pendingAmount']) != pending_count_expected:
        sys.exit(-1)
    if result['maturedAmount'] != matured_expected:
        sys.exit(-1)