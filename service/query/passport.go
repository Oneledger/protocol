package query

import (
	"sort"
	"time"

	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/passport"
	pspt "github.com/Oneledger/protocol/data/passport"
)

func (sv *Service) PSPT_QueryTestInfoByID(req client.PSPTFilterTestInfoRequest, reply *client.PSPTFilterTestInfoReply) error {
	if err := req.Admin.Err(); err != nil {
		return pspt.ErrInvalidIdentifier
	}
	if err := req.Person.Err(); err != nil {
		return pspt.ErrInvalidIdentifier
	}

	// check permission
	permitted, _ := sv.authTokens.HasPermission(req.Org, req.Admin, pspt.PermitQueryTest)
	if !permitted {
		return pspt.ErrPermissionRequired
	}

	// get test data
	infoList, err := sv.Tests.GetTestInfoByID(req.Person, req.Test)
	if err != nil {
		return err
	}

	*reply = client.PSPTFilterTestInfoReply{
		InfoList: infoList,
		Height:   sv.Tests.State.Version(),
	}

	return nil
}

// this function will allow user to query their results without permission check
func (sv *Service) PSPT_QueryTestInfoForUser(req client.PSPTQueryTestInfoForUserReq, reply *client.PSPTQueryTestInfoForUserReply) error {
	if err := req.Person.Err(); err != nil {
		return pspt.ErrInvalidIdentifier
	}

	// get test data
	infoList, err := sv.Tests.GetTestInfoByID(req.Person, req.Test)
	if err != nil {
		return err
	}

	*reply = client.PSPTQueryTestInfoForUserReply{
		InfoList: infoList,
		Height:   sv.Tests.State.Version(),
	}

	return nil
}

func parseTime(st string) (tm time.Time, err error) {
	tm = time.Time{}
	if len(st) > 0 {
		tm, err = time.Parse(time.RFC3339, st)
		if err != nil {
			err = pspt.ErrInvalidTimeStamp
		}
	}
	return
}

func (sv *Service) PSPT_FilterTestInfo(req client.PSPTFilterTestInfoRequest, reply *client.PSPTFilterTestInfoReply) (err error) {
	// search by person directly
	if err = req.Person.Err(); err == nil {
		return sv.PSPT_QueryTestInfoByID(req, reply)
	}

	// check time stamps
	testBegin, err := parseTime(req.TestTimeBegin)
	if err != nil {
		return
	}
	testEnd, err := parseTime(req.TestTimeEnd)
	if err != nil {
		return
	}

	// check time stamps
	analyzeBegin, err := parseTime(req.AnalyzeTimeBegin)
	if err != nil {
		return
	}
	analyzeEnd, err := parseTime(req.AnalyzeTimeEnd)
	if err != nil {
		return
	}

	// check permission
	permitted, _ := sv.authTokens.HasPermission(req.Org, req.Admin, pspt.PermitQueryTest)
	if !permitted {
		return pspt.ErrPermissionRequired
	}

	// filter test data
	replyList := []*pspt.TestInfo{}
	replyMap := map[pspt.UserID]bool{}
	sv.Tests.IterateOrgTests(req.Test, req.UploadedOrg, req.UploadedBy, req.Person,
		func(test pspt.TestType, org pspt.TokenTypeID, uploadedBy, person pspt.UserID, num int) bool {
			// continue if already this person checked before
			if _, ok := replyMap[person]; ok {
				return false
			} else {
				replyMap[person] = true
			}
			// fetch person's records
			infoList, err := sv.Tests.GetTestInfoByID(person, test)
			if err != nil {
				return false
			}

			// filter person's records
			for _, info := range infoList {
				// filter by person if specified
				if req.UploadedBy != "" && req.UploadedBy != uploadedBy {
					continue
				}
				// filter by uploadedOrg if specified
				if req.UploadedOrg != "" && req.UploadedOrg != org {
					continue
				}
				// filter by upload time range
				uploadedAt := info.TimeUpload()
				if !testBegin.IsZero() && uploadedAt.Before(testBegin) {
					continue
				}
				if !testEnd.IsZero() && uploadedAt.After(testEnd) {
					continue
				}
				// filter by analysisOrg if specified
				if req.AnalysisOrg != "" && req.AnalysisOrg != info.AnalysisOrg {
					continue
				}
				// filter by analyzedBy if specified
				if req.AnalyzedBy != "" && req.AnalyzedBy != info.AnalyzedBy {
					continue
				}
				// filter by analysis time range
				analyzedAt := info.TimeAnalyze()
				if !analyzeBegin.IsZero() && analyzedAt.Before(analyzeBegin) {
					continue
				}
				if !analyzeEnd.IsZero() && analyzedAt.After(analyzeEnd) {
					continue
				}
				replyList = append([]*pspt.TestInfo{info}, replyList...)
			}
			return false
		})

	// Sort by upload time
	sort.Slice(replyList, func(i, j int) bool {
		return replyList[i].TimeUpload().After(replyList[j].TimeUpload())
	})

	*reply = client.PSPTFilterTestInfoReply{
		InfoList: replyList,
		Height:   sv.Tests.State.Version(),
	}

	return nil
}

func (sv *Service) PSPT_FilterReadLogs(req client.PSPTFilterReadLogsRequest, reply *client.PSPTFilterReadLogsReply) (err error) {
	// check time stamps
	tBegin, err := parseTime(req.TimeBegin)
	if err != nil {
		return
	}
	tEnd, err := parseTime(req.TimeEnd)
	if err != nil {
		return
	}

	// check permission
	permitted, _ := sv.authTokens.HasPermission(req.Org, req.Admin, pspt.PermitQueryRead)
	if !permitted {
		return pspt.ErrPermissionRequired
	}

	// filter test data
	replyList := []*pspt.ReadLog{}
	sv.Tests.IterateReadLogs(req.Test, req.ReadOrg, req.ReadBy, func(id pspt.UserID, logs []*pspt.ReadLog) bool {
		for _, log := range logs {

			// filter by person if specified
			if req.Person != "" && req.Person != log.Person {
				return false
			}
			// filter by time range
			readAt := log.TimeRead()
			if !tBegin.IsZero() && readAt.Before(tBegin) {
				continue
			}
			if !tEnd.IsZero() && readAt.After(tEnd) {
				continue
			}
			replyList = append(replyList, log)
		}
		return false
	})

	// Sort by read time
	sort.Slice(replyList, func(i, j int) bool {
		return replyList[i].TimeRead().After(replyList[j].TimeRead())
	})

	*reply = client.PSPTFilterReadLogsReply{
		Logs:   replyList,
		Height: sv.Tests.State.Version(),
	}

	return nil
}

func (sv *Service) GetAuthToken(req client.AuthTokenRequest, reply *client.AuthTokenReply) error {
	token, err := sv.authTokens.GetAuthToken(req.Org, req.ID)
	if err != nil {
		return err
	}
	reply.Token = token

	return nil
}

func (sv *Service) GetTokensByOrganization(req client.AuthTokenRequest, reply *client.OrgTokenReply) error {
	tokenList := []*passport.AuthToken{}

	// check permission
	permitted, _ := sv.authTokens.HasPermission(req.Org, req.ID, pspt.PermitQueryTokens)
	if !permitted {
		return pspt.ErrPermissionRequired
	}

	sv.authTokens.IterateOrgTokens(req.Org, func(token *passport.AuthToken) bool {
		tokenList = append(tokenList, token)
		return false
	})

	reply.Tokens = tokenList

	return nil
}
