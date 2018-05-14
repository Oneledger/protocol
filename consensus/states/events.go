package states

type BasicEvent interface {
	New()
}

type EventActionPrevote struct {
	Prevote *Prevote `json:"prevote"`
}

type EventActionPropose struct {
	Propose *Propose `json:"propose"`
}

type EventActionPrecommit struct {
	Precommit *Precommit `json:"precommit"`
}

type EventActionNewHeight struct {
	NewHeight *NewHeight `json:"new height"`
}

type EventActionCommit struct {
	Commit *Commit `json:"commit"`
}

type EventActionDriveChain struct {
	DriveChain *DriveChain `json:"drive chain"`
}