import datetime
import pytz
from testinfo import *
from sdk import *
from datetime import datetime, timedelta

def create_hospital_admins():
    # Create Auth Token for St. Michael's hospital Admin
    hsptAdmin1.send(charlie)
    validate_result(True, hsptAdmin1.get_resp(), "Create Auth Token for Hospital Admin in St. Michael's hospital")

    print Colours.OK_BLUE + 'Query St. Michael Hospital tokens' + Colours.END_C
    tokens = get_auth_tokens('ST-MICHAEL', 'AdminUser01')
    print json.dumps(tokens, indent=2)

    # Create Auth Token for Sunnybrook Hospital Admin
    hsptAdmin2.send(charlie)
    validate_result(True, hsptAdmin2.get_resp(), "Create Auth Token for Hospital Admin in Sunnybrook hospital")

    print Colours.OK_BLUE + 'Query Sunnybrook Hospital tokens' + Colours.END_C
    tokens = get_auth_tokens('SUNNY-BROOK', 'AdminUser02')
    print json.dumps(tokens, indent=2)

def test_upload_update():
    # upload and update test info

    upload_info0.send_add_test(admin1, admin_addresses['admin1'])
    print Colours.OK_GREEN + '$$$$ Upload test info0 succeed $$$$' + Colours.END_C
    update_info0.send_update_test(admin2, admin_addresses['admin2'])
    print Colours.OK_GREEN + '$$$$ Update test info0 succeed $$$$' + Colours.END_C

    upload_info1.send_add_test(admin1, admin_addresses['admin1'])
    print Colours.OK_GREEN + '$$$$ Upload test info1 succeed $$$$' + Colours.END_C
    
    upload_info2.send_add_test(admin2, admin_addresses['admin2'])
    print Colours.OK_GREEN + '$$$$ Upload test info2 succeed $$$$' + Colours.END_C
    update_info2.send_update_test(admin2, admin_addresses['admin2'])
    print Colours.OK_GREEN + '$$$$ Update test info2 succeed $$$$' + Colours.END_C

    upload_info3.send_add_test(admin2, admin_addresses['admin2'])
    print Colours.OK_GREEN + '$$$$ Upload test info3 succeed $$$$' + Colours.END_C
    print Colours.OK_GREEN + '$$$$ Upload test info passed $$$$' + Colours.END_C

def test_filter_test_info():
    # filter test info by person    
    infoList = filter_test_info(COVID19, 'ST-MICHAEL', "AdminUser01", "person1", "", "", 'ST-MICHAEL', "AdminUser01")
    check_record_list(infoList, [upload_info1, final_info0])

    infoList = filter_test_info(COVID19, 'SUNNY-BROOK', "AdminUser02", "person2", "", "", 'SUNNY-BROOK', "AdminUser02")
    check_record_list(infoList, [final_info2])

    infoList = filter_test_info(COVID19, 'SUNNY-BROOK', "AdminUser02", "person3", "", "", 'SUNNY-BROOK', "AdminUser02")
    check_record_list(infoList, [upload_info3])
    print Colours.OK_GREEN + '$$$$ Filter test info by person passed $$$$' + Colours.END_C

    # filter test info by hospital admin
    infoList = filter_test_info(COVID19, 'ST-MICHAEL', "AdminUser01", "", "", "", 'ST-MICHAEL', "AdminUser01")
    check_record_list(infoList, [upload_info1, final_info0])

    infoList = filter_test_info(COVID19, 'SUNNY-BROOK', "AdminUser02", "", "", "", 'SUNNY-BROOK', "AdminUser02")
    check_record_list(infoList, [upload_info3, final_info2])
    print Colours.OK_GREEN + '$$$$ Filter test info by hospital admin passed $$$$' + Colours.END_C

    # filter test info by hospital id
    infoList = filter_test_info(COVID19, 'ST-MICHAEL', "AdminUser01", "", "", "",'ST-MICHAEL', "")
    check_record_list(infoList, [upload_info1, final_info0])

    infoList = filter_test_info(COVID19, 'SUNNY-BROOK', "AdminUser02", "", "", "", 'SUNNY-BROOK', "")
    check_record_list(infoList, [upload_info3, final_info2])
    print Colours.OK_GREEN + '$$$$ Filter test info by hospital id passed $$$$' + Colours.END_C

    # filter test info by test type
    infoList = filter_test_info(COVID19, 'ST-MICHAEL', "AdminUser01", "", "", "", "", "")
    check_record_list(infoList, [upload_info1, upload_info3, final_info2, final_info0])
    print Colours.OK_GREEN + '$$$$ Filter test info by test type passed $$$$' + Colours.END_C

    # filter test info by [test & upload time]
    infoList = filter_test_info(COVID19, 'ST-MICHAEL', "AdminUser01", "", uploaded_at1, "", "", "")
    check_record_list(infoList, [upload_info1, upload_info3, final_info2])
    print Colours.OK_GREEN + '$$$$ Filter test info by upload time passed $$$$' + Colours.END_C

    # filter test info by analysis org
    infoList = filter_test_info(COVID19, 'ST-MICHAEL', "AdminUser01", "", "", "", "", "", "SUNNY-BROOK")
    check_record_list(infoList, [final_info2, final_info0])
    print Colours.OK_GREEN + '$$$$ Filter test info by analysis org passed $$$$' + Colours.END_C

    # cross-filter test info by [test & analyze time]
    infoList = filter_test_info(COVID19, 'ST-MICHAEL', "AdminUser01", "", "", "", "", "", "", "", analyzed_at2)
    check_record_list(infoList, [final_info2])
    print Colours.OK_GREEN + '$$$$ Filter test info by analyze time passed $$$$' + Colours.END_C

def test_read():
    # read person1
    infoList = TestInfo.read('ST-MICHAEL', "AdminUser01", admin_addresses["admin1"], "person1", addrs["person1"], COVID19, read_at0, admin1.password)
    check_record_list(infoList, [upload_info1, final_info0])

    # read person2, scan twice
    infoList = TestInfo.read('ST-MICHAEL', "AdminUser01", admin_addresses["admin1"], "person2", addrs["person2"], COVID19, read_at1, admin1.password)
    check_record_list(infoList, [final_info2])
    infoList = TestInfo.read('ST-MICHAEL', "AdminUser01", admin_addresses["admin1"], "person2", addrs["person2"], COVID19, read_at2, admin1.password)
    check_record_list(infoList, [final_info2])

    # read person3
    infoList = TestInfo.read('SUNNY-BROOK', "AdminUser02", admin_addresses["admin2"], "person3", addrs["person3"], COVID19, read_at2, admin2.password)
    check_record_list(infoList, [upload_info3])

    # query test info from person, this does not require permission
    infoList = read_test_info_from_user(COVID19, "person3")
    check_record_list(infoList, [upload_info3])

    print Colours.OK_GREEN + '$$$$ Read test info passed $$$$' + Colours.END_C

def test_read_logs():
    # filter test info by org, cross query
    logs = filter_read_logs(COVID19, 'SUNNY-BROOK', "AdminUser02", "ST-MICHAEL", "", "", "", "")
    check_record_list(logs, [log2, log1, log0], ReadLog)

    # filter test info by [org+admin]
    logs = filter_read_logs(COVID19, 'ST-MICHAEL', "AdminUser01", "ST-MICHAEL", "AdminUser01", "", "", "")
    check_record_list(logs, [log2, log1, log0], ReadLog)

    # filter test info by [org+admin+person]
    logs = filter_read_logs(COVID19, 'ST-MICHAEL', "AdminUser01", "ST-MICHAEL", "AdminUser01", "person2", "", "")
    check_record_list(logs, [log2, log1], ReadLog)

    # filter test info by [org+time]
    logs = filter_read_logs(COVID19, 'ST-MICHAEL', "AdminUser01", "ST-MICHAEL", "", "", read_at1, "")
    check_record_list(logs, [log2, log1], ReadLog)

    # filter test info by [test+person]
    logs = filter_read_logs(COVID19, 'SUNNY-BROOK', "AdminUser02", "", "", "person2", "", "")
    check_record_list(logs, [log2, log1], ReadLog)

    print Colours.OK_GREEN + '$$$$ Read logs passed $$$$' + Colours.END_C

if __name__ == '__main__':
    create_hospital_admins()

    test_upload_update()

    test_filter_test_info()

    test_read()
    
    test_read_logs()
    