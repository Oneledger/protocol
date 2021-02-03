from sdk import *
import datetime


print Colours.HEADER + '######## Creating Accounts ########' + Colours.END_C
addr_list = addresses()


test_case0 = AdminToken('USER05', TOKEN_SCREENER, 'US_BOARDER', addr_list[3], super_user_addresses['kevin'], kevin.user_id,
                        datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'))
test_case1 = AdminToken('USER06', TOKEN_SCREENER, 'CAD_BOARDER', addr_list[4], super_user_addresses['kevin'], kevin.user_id,
                        datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'))
test_case2 = AdminToken('USER07', TOKEN_SCREENER, 'CAD_BOARDER', addr_list[5], super_user_addresses['kevin'], kevin.user_id,
                        datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'))
test_case3 = AdminToken('USER08', TOKEN_SCREENER, 'CAD_BOARDER', super_user_addresses['charlie'], super_user_addresses['kevin'], kevin.user_id,
                        datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'))
test_case4 = AdminToken('USER06', TOKEN_SCREENER, 'CAD_BOARDER', addr_list[4], super_user_addresses['kevin'], kevin.user_id,
                        datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'))

# Extra Test cases
test_case5 = AdminToken('USER06', TOKEN_SCREENER, 'US_BOARDER', addr_list[4], super_user_addresses['kevin'], kevin.user_id,
                        datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'))
test_case6 = AdminToken('USER06', TOKEN_HOSPITAL, 'HOSPITAL1', addr_list[4], super_user_addresses['kevin'], kevin.user_id,
                        datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'))

# Test Case list
test_cases = [test_case0, test_case1, test_case2, test_case3, test_case4, test_case5]


if __name__ == '__main__':
    # Create Auth Token for the US Boarder
    test_case0.send(kevin)
    validate_result(True, test_case0.get_resp(), "Test Case 0 - Create Auth Token for the US Boarder")

    # Create Auth Token for Hospital Admin
    test_case1.send(kevin)
    validate_result(True, test_case1.get_resp(), "Test Case 1 - Create Auth Token for Screener at Canadian boarder")

    # Build request with wrong signer
    test_case2.send_wrong_signer(kevin)
    validate_result(False, test_case2.get_resp(), "Test Case 2 - Build request with wrong signer")

    # Create auth token for another super user
    test_case3.send(kevin)
    validate_result(True, test_case3.get_resp(), "Test Case 3 - Create auth token for another super user")

    # Create duplicate auth token
    test_case4.send(kevin)
    validate_result(False, test_case4.get_resp(), "Test Case 4 - Create duplicate auth token")

    # Create another screener auth token for same person, different org
    test_case5.send(kevin)
    validate_result(True, test_case5.get_resp(), "Test Case 5 - Create another screener auth token for same person, different org")

    # Create a hospital auth token for same person
    test_case6.send(kevin)
    validate_result(True, test_case6.get_resp(), "Test Case 6 - Create a hospital auth token for same person")

    print ''
    print Colours.OK_BLUE + 'Query Canadian Boarder tokens' + Colours.END_C
    tokens = get_auth_tokens('CAD_BOARDER', 'USER06')
    print json.dumps(tokens, indent=2)

    print ''
    print Colours.OK_BLUE + 'Query US Boarder tokens' + Colours.END_C
    tokens = get_auth_tokens('US_BOARDER', 'USER05')
    print json.dumps(tokens, indent=2)

    print ''
    print Colours.OK_BLUE + 'Query HOSPITAL1 tokens' + Colours.END_C
    tokens = get_auth_tokens('HOSPITAL1', 'USER06')
    print json.dumps(tokens, indent=2)

    print ''
    print Colours.OK_BLUE + 'Query Super User tokens' + Colours.END_C
    tokens = get_auth_tokens('SuperAdminGroup', hao.user_id)
    print json.dumps(tokens, indent=2)

    print Colours.OK_GREEN + '$$$$ Screener Admin Test passed $$$$' + Colours.END_C
