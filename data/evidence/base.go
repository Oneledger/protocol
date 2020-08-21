package evidence

type Options struct {
	// votes count required
	MinVotesRequired int64 `json:"minVotesRequired"`
	// votes boundaries H : H - VotesBlockDiff
	BlockVotesDiff int64 `json:"blockVotesDiff"`

	// penalty 1st iteration
	PenaltyBasePercentage int64 `json:"penaltyBasePercentage"`
	// penalty divisability factor
	PenaltyBaseDecimals int64 `json:"penaltyBaseDecimals"`

	// penalty cut for bounty (as example 13.43 %, stored as 1343)
	PenaltyBountyPercentage int64 `json:"penaltyBountyPercentage"`
	// penalty cut decimals
	PenaltyBountyDecimals int64 `json:"penaltyBountyDecimals"`

	// penalty cut for burn
	PenaltyBurnPercentage int64 `json:"penaltyBurnPercentage"`
	// penalty cut decimals
	PenaltyBurnDecimals int64 `json:"penaltyBurnDecimals"`

	// time to unfreeze validator (number of days) - for 2 scenario
	ValidatorReleaseTime int64 `json:"validatorReleaseTime"`
	// required count to finish voting
	AllegationVotesCount int64 `json:"allegationVotesCount"`

	// allegation persent
	AllegationPercentage int64 `json:"allegationPercentage"`
	// allegation cut decimals
	AllegationDecimals int64 `json:"allegationDecimals"`
	// Active (self staking)
	// Freeze (freeze true, active false)
}
