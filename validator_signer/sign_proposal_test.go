package validator_signer

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"github.com/wealdtech/eth2-signer-api/pb/v1"
	"testing"
)

//func TestMultipleConcurrentProposalSignatures(t *testing.T) {
//	seed,_ := hex.DecodeString("f51883a4c56467458c3b47d06cd135f862a6266fabdfb9e9e4702ea5511375d7")
//	signer,err := setupWithSlashingProtection(seed)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	// save an initial valid proposal
//	t.Run("initial valid proposal", func(t *testing.T) {
//		_,err = signer.SignBeaconProposal(&v1.SignBeaconProposalRequest{
//			Id:                   &v1.SignBeaconProposalRequest_Account{Account:"1"},
//			Domain:               []byte("domain"),
//			Data:                 &v1.BeaconBlockHeader{
//				Slot:                 99,
//				ProposerIndex:        2,
//				ParentRoot:           []byte("Z"),
//				StateRoot:            []byte("Z"),
//				BodyRoot:             []byte("Z"),
//			},
//		})
//		require.NoError(t, err)
//	})
//
//	for i := 0 ; i < 5 ; i++ {
//		t.Parallel()
//		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
//			t.Parallel()
//			// valid
//			go func() {
//				t.Run("valid proposal", func(t *testing.T) {
//					t.Parallel()
//					_,err = signer.SignBeaconProposal(&v1.SignBeaconProposalRequest{
//						Id:                   &v1.SignBeaconProposalRequest_Account{Account:"1"},
//						Domain:               []byte("domain"),
//						Data:                 &v1.BeaconBlockHeader{
//							Slot:                 99,
//							ProposerIndex:        2,
//							ParentRoot:           []byte("Z"),
//							StateRoot:            []byte("Z"),
//							BodyRoot:             []byte("Z"),
//						},
//					})
//					if i ==3 {
//						require.NotNil(t, err)
//					} else {
//						require.NoError(t, err)
//					}
//				})
//			}()
//
//			// double
//			go func() {
//				t.Run("different state root, should error", func(t *testing.T) {
//					t.Parallel()
//					_,err = signer.SignBeaconProposal(&v1.SignBeaconProposalRequest{
//						Id:                   &v1.SignBeaconProposalRequest_Account{Account:"1"},
//						Domain:               []byte("domain"),
//						Data:                 &v1.BeaconBlockHeader{
//							Slot:                 99,
//							ProposerIndex:        2,
//							ParentRoot:           []byte("Z"),
//							StateRoot:            []byte("A"),
//							BodyRoot:             []byte("Z"),
//						},
//					})
//					require.NotNil(t, err)
//					require.EqualError(t, err, "err, slashable proposal: DoubleProposal")
//				})
//			}()
//		})
//	}
//
//
//}

func TestProposalSlashingSignatures(t *testing.T) {
	seed, _ := hex.DecodeString("f51883a4c56467458c3b47d06cd135f862a6266fabdfb9e9e4702ea5511375d7")
	signer, err := setupWithSlashingProtection(seed)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("valid proposal", func(t *testing.T) {
		_, err = signer.SignBeaconProposal(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: _byteArray("83e04069ed28b637f113d272a235af3e610401f252860ed2063d87d985931229458e3786e9b331cd73d9fc58863d9e4b")},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          99,
				ProposerIndex: 2,
				ParentRoot:    []byte("Z"),
				StateRoot:     []byte("Z"),
				BodyRoot:      []byte("Z"),
			},
		})
		require.NoError(t, err)
	})

	t.Run("valid proposal, sign using account name. Should error", func(t *testing.T) {
		_, err = signer.SignBeaconProposal(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_Account{Account: "1"},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          99,
				ProposerIndex: 2,
				ParentRoot:    []byte("Z"),
				StateRoot:     []byte("Z"),
				BodyRoot:      []byte("Z"),
			},
		})
		require.NotNil(t, err)
		require.Error(t, err, "account was not supplied")
	})

	t.Run("double proposal, different state root. Should error", func(t *testing.T) {
		_, err = signer.SignBeaconProposal(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: _byteArray("83e04069ed28b637f113d272a235af3e610401f252860ed2063d87d985931229458e3786e9b331cd73d9fc58863d9e4b")},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          99,
				ProposerIndex: 2,
				ParentRoot:    []byte("Z"),
				StateRoot:     []byte("A"),
				BodyRoot:      []byte("Z"),
			},
		})
		require.NotNil(t, err)
		require.EqualError(t, err, "err, slashable proposal: DoubleProposal")
	})

	t.Run("double proposal, different body root. Should error", func(t *testing.T) {
		_, err = signer.SignBeaconProposal(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: _byteArray("83e04069ed28b637f113d272a235af3e610401f252860ed2063d87d985931229458e3786e9b331cd73d9fc58863d9e4b")},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          99,
				ProposerIndex: 2,
				ParentRoot:    []byte("Z"),
				StateRoot:     []byte("Z"),
				BodyRoot:      []byte("A"),
			},
		})
		require.NotNil(t, err)
		require.EqualError(t, err, "err, slashable proposal: DoubleProposal")
	})

	t.Run("double proposal, different parent root. Should error", func(t *testing.T) {
		_, err = signer.SignBeaconProposal(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: _byteArray("83e04069ed28b637f113d272a235af3e610401f252860ed2063d87d985931229458e3786e9b331cd73d9fc58863d9e4b")},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          99,
				ProposerIndex: 2,
				ParentRoot:    []byte("A"),
				StateRoot:     []byte("Z"),
				BodyRoot:      []byte("Z"),
			},
		})
		require.NotNil(t, err)
		require.EqualError(t, err, "err, slashable proposal: DoubleProposal")
	})

	t.Run("double proposal, different proposer index. Should error", func(t *testing.T) {
		_, err = signer.SignBeaconProposal(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: _byteArray("83e04069ed28b637f113d272a235af3e610401f252860ed2063d87d985931229458e3786e9b331cd73d9fc58863d9e4b")},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          99,
				ProposerIndex: 3,
				ParentRoot:    []byte("Z"),
				StateRoot:     []byte("Z"),
				BodyRoot:      []byte("Z"),
			},
		})
		require.NotNil(t, err)
		require.EqualError(t, err, "err, slashable proposal: DoubleProposal")
	})
}
