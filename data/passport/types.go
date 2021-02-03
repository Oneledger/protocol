package passport

import (
	"errors"
)

type (
	// test related
	UserID      string
	TestType    int
	TestSubType int
	TestResult  int

	// auth token related
	TokenID         string
	TokenType       int
	TokenSubType    int
	TokenTypeID     string
	TokenRole       int
	TokenPermission int
)

func (id UserID) Err() error {
	switch {
	case len(id) == 0:
		return errors.New("user id is empty")
	}
	return nil
}

func (id UserID) String() string {
	return string(id)
}

func GetTestType(test string) TestType {
	switch test {
	case "COVID-19":
		return TestCOVID19
	default:
		return TestInvalid
	}
}

func (test TestType) Err() error {
	if test != TestCOVID19 {
		return errors.New("Unknown test type")
	}
	return nil
}

func (test TestType) String() string {
	switch test {
	case TestCOVID19:
		return "COVID-19"
	default:
		return "Unknown test type"
	}
}

func (test TestSubType) Err() error {
	if (test != TestSubPCR) && (test != TestSubAntigen) && (test != TestSubAntiBody) {
		return errors.New("Unknown test sub type")
	}
	return nil
}

func (test TestSubType) String() string {
	switch test {
	case TestSubAntigen:
		return "Antigen"
	case TestSubAntiBody:
		return "AntiBody"
	case TestSubPCR:
		return "PCR"
	default:
		return "Unknown test sub type"
	}
}

func GetTestSubType(test string) TestSubType {
	switch test {
	case "AntiBody":
		return TestSubAntiBody
	case "Antigen":
		return TestSubAntigen
	case "PCR":
		return TestSubPCR
	default:
		return TestSubInvalid
	}
}

func (result TestResult) Err() error {
	if result != COVID19Positive && result != COVID19Negative && result != COVID19Pending {
		return errors.New("Unknown test result")
	}
	return nil
}

func (result TestResult) String() string {
	switch result {
	case COVID19Positive:
		return "COVID-19 test positive"
	case COVID19Negative:
		return "COVID-19 test negative"
	case COVID19Pending:
		return "COVID-19 test pending"
	default:
		return "Unknown test result"
	}
}

func (id TokenID) Err() error {
	switch {
	case len(id) == 0:
		return errors.New("auth token id is empty")
	}
	return nil
}

func (id TokenID) String() string {
	return string(id)
}

func (typeId TokenTypeID) Err() error {
	switch {
	case len(typeId) == 0:
		return errors.New("auth token type id is empty")
	}
	return nil
}

func (typeId TokenTypeID) String() string {
	return string(typeId)
}

func (mine TokenPermission) HasPermission(require TokenPermission) bool {
	get := mine & require
	return get == require
}
