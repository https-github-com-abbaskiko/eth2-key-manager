package validator_signer

import (
	"encoding/hex"
	"fmt"

	pb "github.com/wealdtech/eth2-signer-api/pb/v1"

	"github.com/bloxapp/eth2-key-manager/core"
)

func (signer *SimpleSigner) SignBeaconProposal(req *pb.SignBeaconProposalRequest) (*pb.SignResponse, error) {
	// 1. get the account
	if req.GetPublicKey() == nil {
		return nil, fmt.Errorf("account was not supplied")
	}
	account, err := signer.wallet.AccountByPublicKey(hex.EncodeToString(req.GetPublicKey()))
	if err != nil {
		return nil, err
	}

	// 2. lock for current account
	signer.lock(account.ID(), "proposal")
	defer signer.unlock(account.ID(), "proposal")

	// 2. check we can even sign this
	if status := signer.slashingProtector.IsSlashableProposal(account.ValidatorPublicKey(), req); status.Status != core.ValidProposal {
		if status.Error != nil {
			return nil, status.Error
		}
		return nil, fmt.Errorf("err, slashable proposal: %s", status.Status)
	}

	// 3. add to protection storage
	if err := signer.slashingProtector.SaveProposal(account.ValidatorPublicKey(), req); err != nil {
		return nil, err
	}

	// 4. generate ssz root hash and sign
	forSig, err := PrepareProposalReqForSigning(req)
	if err != nil {
		return nil, err
	}
	sig, err := account.ValidationKeySign(forSig)
	if err != nil {
		return nil, err
	}
	res := &pb.SignResponse{
		State:     pb.ResponseState_SUCCEEDED,
		Signature: sig.Marshal(),
	}

	return res, nil
}

// PrepareProposalReqForSigning prepares the given proposal request for signing.
// This is exported to allow use it by custom signing mechanism.
func PrepareProposalReqForSigning(req *pb.SignBeaconProposalRequest) ([]byte, error) {
	data := core.ToCoreBlockData(req)
	forSig, err := prepareForSig(data, req.Domain)
	if err != nil {
		return nil, err
	}
	return forSig[:], nil
}
