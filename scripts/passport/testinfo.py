from sdk import *

print Colours.HEADER + '######## Creating Accounts ########' + Colours.END_C
setup_accounts()
addr_list = addresses()

hsptAdmin1 = AdminToken('AdminUser01', TOKEN_HOSPITAL, 'ST-MICHAEL', addr_list[0], super_user_addresses['charlie'],
                        charlie.user_id,
                        datetime.now().strftime('%Y-%m-%d %H:%M:%S'))
hsptAdmin2 = AdminToken('AdminUser02', TOKEN_HOSPITAL, 'SUNNY-BROOK', addr_list[1], super_user_addresses['charlie'],
                        charlie.user_id,
                        datetime.now().strftime('%Y-%m-%d %H:%M:%S'))

# RFC3339 time stamps
now = pytz.UTC.localize(datetime.utcnow())

tested_at0 = (now - timedelta(hours=32)).isoformat('T')
tested_at1 = (now - timedelta(hours=30)).isoformat('T')
tested_at2 = (now - timedelta(hours=24)).isoformat('T')
tested_at3 = (now - timedelta(hours=24)).isoformat('T')

uploaded_at0 = (now - timedelta(hours=18)).isoformat('T')
uploaded_at1 = (now - timedelta(hours=12)).isoformat('T')
uploaded_at2 = (now - timedelta(hours=12)).isoformat('T')
uploaded_at3 = now.isoformat('T')


analyzed_at0 = (now - timedelta(hours=10)).isoformat('T')
analyzed_at1 = (now - timedelta(hours=4)).isoformat('T')
analyzed_at2 = (now - timedelta(hours=5)).isoformat('T')
analyzed_at3 = (now - timedelta(hours=3)).isoformat('T')

read_at0 = (now + timedelta(hours=24)).isoformat('T')
read_at1 = (now + timedelta(hours=36)).isoformat('T')
read_at2 = (now + timedelta(hours=48)).isoformat('T')

# test info person1
# test related fields in update_info is not passed to protocol, just for result validation when reading test info
upload_info0 = TestInfo('xh45yh71kbo6nii9v2k37zxrf05d4i3m1qysd3v16ulfzeb09xslhygenl1c9li1', 'ST-MICHAEL', 'AdminUser01',
                        "person1", COVID19, TestSubAntiBody, "AudaciaBioscience", COVID19Pending, tested_at0,
                        uploaded_at0, "Edward", "showed similar symptoms")
update_info0 = TestInfo('xh45yh71kbo6nii9v2k37zxrf05d4i3m1qysd3v16ulfzeb09xslhygenl1c9li1', 'ST-MICHAEL', 'AdminUser01',
                        "person1", COVID19, TestSubAntiBody, "AudaciaBioscience", COVID19Negative, tested_at0,
                        uploaded_at0, "Edward", "some more notes", 'SUNNY-BROOK', analyzed_at0, "AdminUser02")
final_info0 = TestInfo('xh45yh71kbo6nii9v2k37zxrf05d4i3m1qysd3v16ulfzeb09xslhygenl1c9li1', 'ST-MICHAEL', 'AdminUser01',
                       "person1", COVID19, TestSubAntiBody, "AudaciaBioscience", COVID19Negative, tested_at0,
                       uploaded_at0, "Edward", "showed similar symptoms\nsome more notes (notes updated at: " + analyzed_at0 + ")\n",
                       'SUNNY-BROOK', analyzed_at0, "AdminUser02")

upload_info1 = TestInfo('79x1w29dbjg6v9gzmrjwdbdunbc1teii54eaq761mqwmbnbm6r8pqo61bdog8tco', 'ST-MICHAEL', 'AdminUser01',
                        "person1", COVID19, TestSubAntigen, "MAG ELISA", COVID19Pending, tested_at3,
                        uploaded_at3, "Edward", "want to test again")

# test info person2 person3
upload_info2 = TestInfo('28s149ez17tms7xo9pz13f8t3tbez78eh0fsedlqjw8ds3udprg60ky1bb5jcxw6', 'SUNNY-BROOK', 'AdminUser02',
                        "person2", COVID19, TestSubAntigen, "Gnomegen", COVID19Pending,
                        tested_at1, uploaded_at1, "Edward", "had a headache")
update_info2 = TestInfo('28s149ez17tms7xo9pz13f8t3tbez78eh0fsedlqjw8ds3udprg60ky1bb5jcxw6', 'SUNNY-BROOK', 'AdminUser02',
                        "person2", COVID19, TestSubAntigen, "Gnomegen", COVID19Positive,
                        tested_at1, uploaded_at1, "Edward", "some more notes", 'SUNNY-BROOK', analyzed_at2, "AdminUser02")
final_info2 = TestInfo('28s149ez17tms7xo9pz13f8t3tbez78eh0fsedlqjw8ds3udprg60ky1bb5jcxw6', 'SUNNY-BROOK', 'AdminUser02',
                        "person2", COVID19, TestSubAntigen, "Gnomegen", COVID19Positive,
                        tested_at1, uploaded_at1, "Edward", "had a headache\nsome more notes (notes updated at: " + analyzed_at2 + ")\n",
                       'SUNNY-BROOK', analyzed_at2, "AdminUser02")


upload_info3 = TestInfo('db1u6k4jb02l73ls1zyi4xs4lb9w6fbcaipbzonv08eb8y0g7oter1jrzlxceix8', 'SUNNY-BROOK', 'AdminUser02',
                        "person3", COVID19, TestSubPCR, "AudaciaBioscience", COVID19Pending, tested_at2,
                        uploaded_at2, "Edward", "had a fever")

# read logs
log0 = ReadLog('ST-MICHAEL', "AdminUser01", admin_addresses["admin1"], "person1", addrs["person1"], COVID19, read_at0)
log1 = ReadLog('ST-MICHAEL', "AdminUser01", admin_addresses["admin1"], "person2", addrs["person2"], COVID19, read_at1)
log2 = ReadLog('ST-MICHAEL', "AdminUser01", admin_addresses["admin1"], "person2", addrs["person2"], COVID19, read_at2)
log3 = ReadLog('SUNNY-BROOK', "AdminUser02", admin_addresses["admin2"], "person3", addrs["person3"], COVID19, read_at2)


def check_record_list(actual, expected, Record=TestInfo):
    if len(actual) != len(expected):
        sys.exit(-1)
    # check records
    for i, values in enumerate(actual):
        info = Record()
        info.from_dict(**values)
        if info != expected[i]:
            sys.exit(-1)
