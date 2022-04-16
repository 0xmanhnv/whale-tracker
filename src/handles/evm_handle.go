package handles

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strings"

	"whale-tracker/src/blockchain"
	"whale-tracker/src/blockchain/token"
	"whale-tracker/src/services"
	whaleutils "whale-tracker/src/whale-utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.mongodb.org/mongo-driver/bson"
)

// LogTransfer ..
type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
	TxHash common.Hash
}

type InfoTransfer struct {
	IsPending bool
	Value     *big.Int
}

func GetCurrentBlock(chainId int) uint64 {
	client := blockchain.GetClient(chainId)
	currentBlock, _ := client.BlockNumber(context.Background())

	return currentBlock
}

/*
Load logs from token address
@param: chainId - network chain id
@param: tokenAddress - token address on evm or bsc
@param: fromBlock - from block
@param: toBlock - to block

return []types.Log
*/
func LoadLogs(client *ethclient.Client, addresses []common.Address, fromBlock int64, toBlock int64, topics [][]common.Hash) []types.Log {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
		Addresses: addresses,
		Topics:    topics,
	}

	logs, err := client.FilterLogs(context.Background(), query)

	if err != nil {
		log.Fatal(err)
	}

	return logs
}

func LogHandle(logs []types.Log) {
	logTransferSig := []byte("Transfer(address,address,uint256)")
	LogApprovalSig := []byte("Approval(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	logApprovalSigHash := crypto.Keccak256Hash(LogApprovalSig)
	var transferHash []common.Hash

	contractAbi, err := abi.JSON(strings.NewReader(string(token.TokenABI)))

	if err != nil {
		log.Fatal(err)
	}

	for _, vLog := range logs {

		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			fmt.Printf("Log Name: Transfer\n")

			var transferEvent LogTransfer

			i, err := contractAbi.Unpack("Transfer", vLog.Data)

			if err != nil {
				log.Fatal(err)
			}

			print(i)

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
			transferEvent.TxHash = vLog.TxHash

			fmt.Printf("From: %s\n", transferEvent.From.Hex())
			fmt.Printf("To: %s\n", transferEvent.To.Hex())
			fmt.Printf("TxHash: %s\n", transferEvent.TxHash)

			transferHash = append(transferHash, transferEvent.TxHash)

		case logApprovalSigHash.Hex():
			//TODO: handle fmt.Printf("Log Name: Approval\n")
		}

		fmt.Printf("\n\n")

		// fmt.Printf("\nProcessed events from block %v", log.BlockNumber)

		// if vlog.Address == common.HexToAddress(address) {
		// 	// fmt.Println(common.HexToAddress(address))
		// 	fmt.Println(vlog.Address)
		// }
	}

	fmt.Println("LogHandle Done")
}

func LogHandleToWhale(client *ethclient.Client, logs []types.Log, tokenAddresses []common.Address) {

	for _, vLog := range logs {
		tokenAddress := vLog.Address
		priceToken := GetPriceFromPanCakeSwap(tokenAddress.String())

		switch {
		case AddressesContains(tokenAddress, tokenAddresses):

			from := vLog.Topics[1].Hex()
			to := vLog.Topics[2].Hex()
			// TxHash := vLog.TxHash.Hex()

			/* Handle calc price token */
			const spec = 200
			price := priceToken.Data.Price.SetPrec(spec)

			isWhaleValue := new(big.Float).SetInt(big.NewInt(50000)).SetPrec(spec) //usd

			// infoTransfer := LoadInfoTransfer(client, vLog.TxHash.Hex())

			fmt.Println("----" + tokenAddress.String() + "----")

			// handle if TO is not smart contract
			if !CheckIsSmartContract(client, common.HexToAddress(to)) {
				isWhale := services.FindOnenWhale(tokenAddress.String(), common.HexToAddress(to).String())

				if isWhale == nil {
					// handle get balance
					toBalanceInt := GetBalanceToken(client, tokenAddress, common.HexToAddress(to))
					toBalanceEther := whaleutils.WeiToEther(toBalanceInt)

					toValue := new(big.Float)
					toValue.Mul(price, toBalanceEther)

					if toValue.Cmp(isWhaleValue) >= 0 {

						_, err := services.CreateWhale(bson.D{
							{"address", common.HexToAddress(to).String()},
							{"token_address", tokenAddress.String()},
						})

						if err != nil {
							fmt.Println("Whale add success")
						}
					}
				}
			}

			// handle if FROM is not smart contract
			if !CheckIsSmartContract(client, common.HexToAddress(from)) {
				isWhale := services.FindOnenWhale(tokenAddress.String(), common.HexToAddress(from).String())

				if isWhale == nil {
					// handle get balance
					fromBalanceInt := GetBalanceToken(client, tokenAddress, common.HexToAddress(from))
					fromBalanceEther := whaleutils.WeiToEther(fromBalanceInt)

					fromValue := new(big.Float)
					fromValue.Mul(price, fromBalanceEther)

					if fromValue.Cmp(isWhaleValue) >= 0 {
						_, err := services.CreateWhale(bson.D{
							{"address", common.HexToAddress(from).String()},
							{"token_address", tokenAddress.String()},
						})

						if err != nil {
							fmt.Println("Whale add success")
						}
					}
				}
			}

			// GetTokens(vLog)

			// fmt.Println("----------test")
			// fmt.Println(infoTransfer)
			// fmt.Println(price)

			// // end calc price token

			// fmt.Println("FROM: Token balance of " + common.HexToAddress(from).String() + " : " + GetBalanceToken(client, tokenAddress, common.HexToAddress(from)).String())

			// fmt.Println("TO:Token balance of " + common.HexToAddress(to).String() + " : " + GetBalanceToken(client, tokenAddress, common.HexToAddress(to)).String())

			// fmt.Println(vLog.BlockNumber)
			// fmt.Println(tokenAddress)
			// fmt.Println(TxHash)
			// fmt.Println(infoTransfer)
			// fmt.Println()

		}
	}
}

func CheckIsSmartContract(client *ethclient.Client, address common.Address) bool {
	bytecode, err := client.CodeAt(context.Background(), address, nil) // nil is latest block
	if err != nil {
		log.Fatal(err)
		return false
	}

	return len(bytecode) > 0
}

func GetTokens(vLog types.Log) {

	// read value token
	contractAbi, err := abi.JSON(strings.NewReader(string(token.TokenABI)))
	if err != nil {
		log.Fatal(err)
	}

	value, err := contractAbi.Unpack("Transfer", vLog.Data)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(value[0])
}

func LoadInfoTransfer(client *ethclient.Client, TxHash string) InfoTransfer {
	tx, isPending, err := client.TransactionByHash(context.Background(), common.HexToHash(TxHash))

	if err != nil {
		log.Fatal(err)
	}

	var infoTransfer InfoTransfer

	infoTransfer.IsPending = isPending
	infoTransfer.Value = tx.Value()

	return infoTransfer
}

func CheckIsWhale(client *ethclient.Client, address common.Address) {
	//TODO: check is whale
}

func AddressesContains(address common.Address, addresses []common.Address) bool {
	for _, i := range addresses {
		if i == address {
			return true
		}
	}
	return false
}

func CheckAdressIsSwapAddress(address common.Address) bool {
	// 0x10ed43c718714eb63d5aa57b78b54704e256024e
	swappAddresses := []common.Address{
		common.HexToAddress("0x10ed43c718714eb63d5aa57b78b54704e256024e"),
	}

	return AddressesContains(address, swappAddresses)
}

func GetBalanceToken(client *ethclient.Client, tokenAddress common.Address, address common.Address) *big.Int {

	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}

	return bal
}

type DataPrice struct {
	Name      string    `json:"name"`
	Symbol    string    `json:"symbol"`
	Price     big.Float `json:"price"`
	Price_BNB big.Float `json:"price_BNB"`
}

type ResponsePrice struct {
	Data DataPrice `json:"data"`
}

func GetPriceFromPanCakeSwap(address string) ResponsePrice {

	UrlApi := "https://api.pancakeswap.info/api/v2/tokens/" + address

	resp, err := http.Get(UrlApi)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var result ResponsePrice
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	fmt.Println(result)

	return result
}
