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
	PenaltyBountyPercentage int64 `json:"penaltyBountyPerncentage"`
	// penalty cut decimals
	PenaltyBountyDecimals int64 `json:"penaltyBountyDecimals"`

	// penalty cut for burn
	PenaltyBurnPercentage int64 `json:"penaltyBurnPerncentage"`
	// penalty cut decimals
	PenaltyBurnDecimals int64 `json:"penaltyBurnDecimals"`

	// time to unfreeze validator (number of days) - for 2 scenario
	ValidatorReleaseTime int64 `json:"validatorReleaseTime"`
	// required validator votes
	ValidatorVotePercentage int64 `json:"validatorVotePercentage"`
	// required validator decimals
	ValidatorVoteDecimals int64 `json:"validatorVoteDecimals"`

	// allegation persent
	AllegationPercentage int64 `json:"allegationPercentage"`
	// allegation cut decimals
	AllegationDecimals int64 `json:"allegationDecimals"`
}
