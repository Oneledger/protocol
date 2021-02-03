from sdk import *
import datetime


print Colours.HEADER + '######## Creating Accounts ########' + Colours.END_C
setup_accounts()
addr_list = addresses()


test_case0 = AdminToken('USER01', TOKEN_HOSPITAL, 'ST_MICHAEL', addr_list[0], super_user_addresses['kevin'], kevin.user_id,
                        datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'))
test_case1 = AdminToken('USER02', TOKEN_HOSPITAL, 'TORONTO_GENERAL', addr_list[1], super_user_addresses['kevin'], kevin.user_id,
                        datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'))
test_case2 = AdminToken('USER03', TOKEN_HOSPITAL, 'TORONTO_GENERAL', addr_list[2], super_user_addresses['kevin'], kevin.user_id,
                        datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'))
test_case3 = AdminToken('USER04', TOKEN_HOSPITAL, 'TORONTO_GENERAL', super_user_addresses['charlie'], super_user_addresses['kevin'], kevin.user_id,
                        datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'))
test_case4 = AdminToken('USER02', TOKEN_HOSPITAL, 'TORONTO_GENERAL', addr_list[1], super_user_addresses['kevin'], kevin.user_id,
                        datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'))

# Extra Test cases
test_case5 = AdminToken('USER02', TOKEN_HOSPITAL, 'UMBRELLA', addr_list[1], super_user_addresses['kevin'], kevin.user_id,
                        datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'))
test_case6 = AdminToken('USER02', TOKEN_SCREENER, 'BORDER1', addr_list[1], super_user_addresses['kevin'], kevin.user_id,
                        datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'))

# Test Case list
test_cases = [test_case0, test_case1, test_case2, test_case3, test_case4, test_case5]


if __name__ == '__main__':
    # Create Auth Token for St. Michael's hospital
    test_case0.send(kevin)
    validate_result(True, test_case0.get_resp(), "Test Case 0 - Create Auth Token for St. Michael's hospital")

    # Create Auth Token for Hospital Admin
    test_case1.send(kevin)
    validate_result(True, test_case1.get_resp(), "Test Case 1 - Create Auth Token for Hospital Admin in Toronto General")

    # Build request with wrong signer
    test_case2.send_wrong_signer(kevin)
    validate_result(False, test_case2.get_resp(), "Test Case 2 - Build request with wrong signer")

    # Create auth token for another super user
    test_case3.send(kevin)
    validate_result(True, test_case3.get_resp(), "Test Case 3 - Create auth token for another super user")

    # Create duplicate auth token
    test_case4.send(kevin)
    validate_result(False, test_case4.get_resp(), "Test Case 4 - Create duplicate auth token")

    # Create another hospital admin auth token for a user that has already has a token with different hospital
    test_case5.send(kevin)
    validate_result(True, test_case5.get_resp(), "Test Case 5 - Create another hospital auth token for same person with different hospital")

    # Create another hospital admin auth token for a user that has already has a token with different hospital, should succeed
    test_case6.send(kevin)
    validate_result(True, test_case6.get_resp(), "Test Case 6 - Create a screener auth token for same person")

    print ''
    print Colours.OK_BLUE + 'Query Tonto General Hospital tokens' + Colours.END_C
    tokens = get_auth_tokens('TORONTO_GENERAL', 'USER02')
    print json.dumps(tokens, indent=2)

    print ''
    print Colours.OK_BLUE + 'Query St. Michael Hospital tokens' + Colours.END_C
    tokens = get_auth_tokens('ST_MICHAEL', 'USER01')
    print json.dumps(tokens, indent=2)

    print ''
    print Colours.OK_BLUE + 'Query UMBRELLA tokens' + Colours.END_C
    tokens = get_auth_tokens('UMBRELLA', 'USER02')
    print json.dumps(tokens, indent=2)

    print ''
    print Colours.OK_BLUE + 'Query BORDER1 tokens' + Colours.END_C
    tokens = get_auth_tokens('BORDER1', 'USER02')
    print json.dumps(tokens, indent=2)

    print ''
    print Colours.OK_BLUE + 'Query Super User tokens' + Colours.END_C
    tokens = get_auth_tokens('SuperAdminGroup', kevin.user_id)
    print json.dumps(tokens, indent=2)

    print Colours.OK_GREEN + '$$$$ Hospital Admin Test passed $$$$' + Colours.END_C
