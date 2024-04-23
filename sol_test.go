package main

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/pkg/pointer"
	"github.com/blocto/solana-go-sdk/program/associated_token_account"
	"github.com/blocto/solana-go-sdk/program/memo"
	"github.com/blocto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/blocto/solana-go-sdk/program/system"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/mr-tron/base58"
	"log"
	"testing"
)

// this function is working, no action requried
func TestTransfer(t *testing.T) {

	fP, _ := types.AccountFromBase58("5zq5A6JG2FDhe9h7Mr2AWB6fHFzbzVxRZmcqi8t37ugy8aVbYxdPZVagGh5o37yZpAfEvNfBaJXGiwmP3EG9Y3EG")

	c := client.NewClient(rpc.DevnetRPCEndpoint)

	bh, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("get recent block hash error, err: %v\n", err)
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{feePayer},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        fP.PublicKey,
			RecentBlockhash: bh.Blockhash,
			Instructions: []types.Instruction{
				memo.BuildMemo(memo.BuildMemoParam{
					SignerPubkeys: []common.PublicKey{fP.PublicKey},
					Memo:          []byte("12345"),
				}),
				system.Transfer(system.TransferParam{
					From:   fP.PublicKey,
					To:     common.PublicKeyFromString("5yUdx5qUPiqGc2jJAb8xPVYvtVMbSFkAC4KMn2DehUe4"),
					Amount: 1e8,
				}),
			},
		}),
	})

	if err != nil {
		log.Fatalf("failed to send money: %s", err)
	}

	txHash, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("failed: %s", err)
	}

	fmt.Println("tx hash:")
	fmt.Println(txHash)

}

// this function should mint an edition fo the master collection nft to a new wallet, this does not work.
func TestMintTo(t *testing.T) {

	c := client.NewClient(rpc.DevnetRPCEndpoint)

	mint := common.PublicKeyFromString("7v82o655Y3QsQixmpARGCs8xnzDCCE7UU2R2KQz2LLjm")

	receiver := common.PublicKeyFromString("5yUdx5qUPiqGc2jJAb8xPVYvtVMbSFkAC4KMn2DehUe4")

	res, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("get recent block hash error, err: %v\n", err)
	}

	ata, _, err := common.FindAssociatedTokenAddress(receiver, mint)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: res.Blockhash,
			Instructions: []types.Instruction{
				token.MintTo(token.MintToParam{
					Mint:   mint,
					Auth:   feePayer.PublicKey,
					To:     ata,
					Amount: 1,
				}),
			},
		}),
		Signers: []types.Account{feePayer},
	})
	if err != nil {
		log.Fatalf("generate tx error, err: %v\n", err)
	}

	txhash, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("send raw tx error, err: %v\n", err)
	}

	log.Println("txhash:", txhash)

}

// This function is working, no action required.
func TestSendNewMint(t *testing.T) {

	c := client.NewClient(rpc.DevnetRPCEndpoint)

	mint := common.PublicKeyFromString("7v82o655Y3QsQixmpARGCs8xnzDCCE7UU2R2KQz2LLjm")

	receiver := common.PublicKeyFromString("5yUdx5qUPiqGc2jJAb8xPVYvtVMbSFkAC4KMn2DehUe4")

	fromATA, _, err := common.FindAssociatedTokenAddress(feePayer.PublicKey, mint)
	if err != nil {
		log.Fatalf("Error finding associated token account: %v", err)
	}

	toATA, _, err := common.FindAssociatedTokenAddress(receiver, mint)
	if err != nil {
		log.Fatalf("Error finding associated token account: %v", err)
	}

	res, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("get recent block hash error, err: %v\n", err)
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{feePayer},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: res.Blockhash,
			Instructions: []types.Instruction{
				token.TransferChecked(token.TransferCheckedParam{
					From:   fromATA,
					To:     toATA,
					Mint:   mint,
					Auth:   feePayer.PublicKey,
					Amount: 1,
				}),
			},
		}),
	})

	if err != nil {
		log.Fatalf("generate tx error, err: %v\n", err)
	}

	txhash, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("send raw tx error, err: %v\n", err)
	}

	log.Println("txhash:", txhash)
}

//Do not worry about the test below
//func TestCreateToken(t *testing.T) {
//
//	c := client.NewClient(rpc.DevnetRPCEndpoint)
//
//	// create an mint account
//	mint := types.NewAccount()
//	fmt.Println("mint:", mint.PublicKey.ToBase58())
//
//	tokenMetadataPubkey, err := token_metadata.GetTokenMetaPubkey(mint.PublicKey)
//	if err != nil {
//		log.Fatalf("failed to find a valid token metadata, err: %v", err)
//	}
//
//	ata, _, err := common.FindAssociatedTokenAddress(feePayer.PublicKey, mint.PublicKey)
//	if err != nil {
//		log.Fatalf("failed to find a valid ata, err: %v", err)
//	}
//
//	// get rent
//	rentExemptionBalance, err := c.GetMinimumBalanceForRentExemption(
//		context.Background(),
//		token.MintAccountSize,
//	)
//	if err != nil {
//		log.Fatalf("get min balacne for rent exemption, err: %v", err)
//	}
//
//	res, err := c.GetLatestBlockhash(context.Background())
//	if err != nil {
//		log.Fatalf("get recent block hash error, err: %v\n", err)
//	}
//
//	tx, err := types.NewTransaction(types.NewTransactionParam{
//		Message: types.NewMessage(types.NewMessageParam{
//			FeePayer:        feePayer.PublicKey,
//			RecentBlockhash: res.Blockhash,
//			Instructions: []types.Instruction{
//				system.CreateAccount(system.CreateAccountParam{
//					From:     feePayer.PublicKey,
//					New:      mint.PublicKey,
//					Owner:    common.TokenProgramID,
//					Lamports: rentExemptionBalance,
//					Space:    token.MintAccountSize,
//				}),
//				token.InitializeMint(token.InitializeMintParam{
//					Decimals:   8,
//					Mint:       mint.PublicKey,
//					MintAuth:   feePayer.PublicKey,
//					FreezeAuth: nil,
//				}),
//				token_metadata.CreateMetadataAccountV3(token_metadata.CreateMetadataAccountV3Param{
//					Metadata:                tokenMetadataPubkey,
//					Mint:                    mint.PublicKey,
//					MintAuthority:           feePayer.PublicKey,
//					Payer:                   feePayer.PublicKey,
//					UpdateAuthority:         feePayer.PublicKey,
//					UpdateAuthorityIsSigner: true,
//					IsMutable:               false,
//					Data: token_metadata.DataV2{
//						Name:                 "a nude croc",
//						Symbol:               "CROC",
//						Uri:                  "https://cmgwtp4c55s7hlrtudc43vncgef4clnp7pl6xcpkg62soisatppa.arweave.net/Ew1pv4LvZfOuM6DFzdWiMQvBLa_71-uJ6je1JyJAm94",
//						SellerFeeBasisPoints: 100,
//						Creators: &[]token_metadata.Creator{
//							{
//								Address:  feePayer.PublicKey,
//								Verified: true,
//								Share:    100,
//							},
//						},
//						Uses: nil,
//					},
//					CollectionDetails: nil,
//				}),
//				associated_token_account.Create(associated_token_account.CreateParam{
//					Funder:                 feePayer.PublicKey,
//					Owner:                  feePayer.PublicKey,
//					Mint:                   mint.PublicKey,
//					AssociatedTokenAccount: ata,
//				}),
//				token.MintTo(token.MintToParam{
//					Mint:   mint.PublicKey,
//					To:     ata,
//					Auth:   feePayer.PublicKey,
//					Amount: 50000000000,
//				}),
//			},
//		}),
//		Signers: []types.Account{feePayer, mint},
//	})
//	if err != nil {
//		log.Fatalf("generate tx error, err: %v\n", err)
//	}
//
//	txhash, err := c.SendTransaction(context.Background(), tx)
//	if err != nil {
//		log.Fatalf("send tx error, err: %v\n", err)
//	}
//
//	log.Println("txhash:", txhash)
//} /  /

// this function partially works, it will create an nft with a master edition. You should run this twice, once to create the parent master edition with a 1 of 1 supply.
// then again to create a child master edition that is tied to the former via collection, with a supply of 5. the master edition also do not show as verified.
func TestMintMasterEdition(t *testing.T) {

	c := client.NewClient(rpc.DevnetRPCEndpoint)

	//receiver := common.PublicKeyFromString("5yUdx5qUPiqGc2jJAb8xPVYvtVMbSFkAC4KMn2DehUe4")

	mint := types.NewAccount()

	mint.PublicKey.ToBase58()

	pk := base58.Encode(mint.PrivateKey)

	fmt.Printf("Mint private key: %s\n", pk)

	fmt.Printf("NFT: %v\n", mint.PublicKey.ToBase58())

	//collection := types.NewAccount()
	//
	//fmt.Printf("collection: %v\n", collection.PublicKey.ToBase58())

	ata, _, err := common.FindAssociatedTokenAddress(feePayer.PublicKey, mint.PublicKey)
	if err != nil {
		log.Fatalf("failed to find a valid ata, err: %v", err)
	}

	fmt.Printf("ata account public key: %s\n", ata.ToBase58())

	tokenMetadataPubkey, err := token_metadata.GetTokenMetaPubkey(mint.PublicKey)
	if err != nil {
		log.Fatalf("failed to find a valid token metadata, err: %v", err)
	}

	masterEditionPubKey, err := token_metadata.GetMasterEdition(mint.PublicKey)
	if err != nil {
		log.Fatalf("failed to find a valid token metadata, err: %v", err)
	}

	fmt.Printf("token metadata pubkey: %s\n", tokenMetadataPubkey.ToBase58())

	fmt.Printf("token master edition pubkey: %s\n", masterEditionPubKey.ToBase58())

	mintAccountRent, err := c.GetMinimumBalanceForRentExemption(context.Background(), token.MintAccountSize)
	if err != nil {
		log.Fatalf("failed to get mint account rent, err: %v", err)
	}

	recentBlockhashResponse, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("failed to get recent blockhash, err: %v", err)
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{mint, feePayer},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: recentBlockhashResponse.Blockhash,
			Instructions: []types.Instruction{
				system.CreateAccount(system.CreateAccountParam{
					From:     feePayer.PublicKey,
					New:      mint.PublicKey,
					Owner:    common.TokenProgramID,
					Lamports: mintAccountRent,
					Space:    token.MintAccountSize,
				}),
				token.InitializeMint(token.InitializeMintParam{
					Decimals:   0,
					Mint:       mint.PublicKey,
					MintAuth:   feePayer.PublicKey,
					FreezeAuth: &feePayer.PublicKey,
				}),
				token_metadata.CreateMetadataAccountV3(token_metadata.CreateMetadataAccountV3Param{
					Metadata:                tokenMetadataPubkey,
					Mint:                    mint.PublicKey,
					MintAuthority:           feePayer.PublicKey,
					Payer:                   feePayer.PublicKey,
					UpdateAuthority:         feePayer.PublicKey,
					UpdateAuthorityIsSigner: true,
					IsMutable:               true,
					Data: token_metadata.DataV2{
						Name:                 "Sunday Nudes",
						Symbol:               "CrocShirt",
						Uri:                  "https://arweave.net/OTWb-97jvMPl5g4gVHKuXaMpwNSJtbPPJAT157zGwcE",
						SellerFeeBasisPoints: 100,
						Creators: &[]token_metadata.Creator{
							{
								Address:  feePayer.PublicKey,
								Verified: true,
								Share:    100,
							},
						},
						//Collection: &token_metadata.Collection{
						//	Verified: true,
						//	Key:      common.PublicKeyFromString("9LAifmytpeteaakKWvyk4rLodY5UUsTZtMGzU6dLakcK"), //this will not validate as true, need to fix this.
						//},
						Uses: nil,
					},
					CollectionDetails: nil,
				}),
				associated_token_account.Create(associated_token_account.CreateParam{
					Funder:                 feePayer.PublicKey,
					Owner:                  feePayer.PublicKey,
					Mint:                   mint.PublicKey,
					AssociatedTokenAccount: ata,
				}),
				token.MintTo(token.MintToParam{
					Mint:   mint.PublicKey,
					To:     ata,
					Auth:   feePayer.PublicKey,
					Amount: 1,
				}),
				token_metadata.CreateMasterEditionV3(token_metadata.CreateMasterEditionParam{
					Edition:         masterEditionPubKey,
					Mint:            mint.PublicKey,
					UpdateAuthority: feePayer.PublicKey,
					MintAuthority:   feePayer.PublicKey,
					Metadata:        tokenMetadataPubkey,
					Payer:           feePayer.PublicKey,
					MaxSupply:       pointer.Get[uint64](1),
				}),
			},
		}),
	})
	if err != nil {
		log.Fatalf("failed to new a tx, err: %v", err)
	}

	sig, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("failed to send tx, err: %v", err)
	}

	fmt.Println("txid:", sig)

}

// this function does not work, it should mint an edition from the second master edition created above. however it fails signature and other issues. I believe we are using the wrong keys in the wrong places.
func TestMintFromMasterEdition(t *testing.T) {

	c := client.NewClient(rpc.DevnetRPCEndpoint)

	//receiver := common.PublicKeyFromString("5yUdx5qUPiqGc2jJAb8xPVYvtVMbSFkAC4KMn2DehUe4")

	mint := types.NewAccount()

	mint.PublicKey.ToBase58()

	pk := base58.Encode(mint.PrivateKey)

	fmt.Printf("Mint private key: %s\n", pk)

	fmt.Printf("NFT: %v\n", mint.PublicKey.ToBase58())

	//collection := types.NewAccount()
	//
	//fmt.Printf("collection: %v\n", collection.PublicKey.ToBase58())

	ata, _, err := common.FindAssociatedTokenAddress(feePayer.PublicKey, mint.PublicKey)
	if err != nil {
		log.Fatalf("failed to find a valid ata, err: %v", err)
	}

	fmt.Printf("ata account public key: %s\n", ata.ToBase58())

	tokenMetadataPubkey, err := token_metadata.GetTokenMetaPubkey(mint.PublicKey)
	if err != nil {
		log.Fatalf("failed to find a valid token metadata, err: %v", err)
	}

	masterEditionPubKey, err := token_metadata.GetMasterEdition(mint.PublicKey)
	if err != nil {
		log.Fatalf("failed to find a valid token metadata, err: %v", err)
	}

	fmt.Printf("token metadata pubkey: %s\n", tokenMetadataPubkey.ToBase58())

	fmt.Printf("token master edition pubkey: %s\n", masterEditionPubKey.ToBase58())

	//mintAccountRent, err := c.GetMinimumBalanceForRentExemption(context.Background(), token.MintAccountSize)
	//if err != nil {
	//	log.Fatalf("failed to get mint account rent, err: %v", err)
	//}

	recentBlockhashResponse, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("failed to get recent blockhash, err: %v", err)
	}

	editionMark, err := token_metadata.GetEditionMark(mint.PublicKey, 1)
	if err != nil {
		log.Fatalf("failed to get recent blockhash, err: %v", err)
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: recentBlockhashResponse.Blockhash,
			Instructions: []types.Instruction{
				token_metadata.MintNewEditionFromMasterEditionViaToken(token_metadata.MintNewEditionFromMasterEditionViaTokeParam{
					NewMetaData:                tokenMetadataPubkey,
					NewEdition:                 mint.PublicKey,
					MasterEdition:              common.PublicKeyFromString("6ceBeT5wbwbbPZEwGvC5vh4s68NvbCpAeRefdvLfuZyh"),
					NewMint:                    mint.PublicKey,
					EditionMark:                editionMark,
					NewMintAuthority:           feePayer.PublicKey,
					Payer:                      feePayer.PublicKey,
					TokenAccountOwner:          feePayer.PublicKey,
					TokenAccount:               common.PublicKeyFromString("BcQfpWd6gDnZCumNbDey1F2BfT2XSusfBK44UEQY7iqa"),
					NewMetadataUpdateAuthority: feePayer.PublicKey,
					MasterMetadata:             common.PublicKeyFromString("BJsQckGDuieimFbi6dJUMSKTKniCbtryzCE2PBRooozQ"),
					Edition:                    1,
				}),
			},
		}),
		Signers: []types.Account{feePayer},
	})
	if err != nil {
		log.Fatalf("failed to new a tx, err: %v", err)
	}

	sig, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("failed to send tx, err: %v", err)
	}

	fmt.Println("txid:", sig)
}
